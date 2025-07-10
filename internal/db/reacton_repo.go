package db

import (
	"context"
	"time"

	"github.com/gocql/gocql"
	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
	"github.com/malyshEvhen/meow_mingle/pkg/logger"
)

type reactionRepository struct {
	session *gocql.Session
	logger  *logger.Logger
}

// ReactionRepository defines the interface for reaction data operations
type ReactionRepository interface {
	Save(ctx context.Context, targetID, authorID, content string) error
	SaveReaction(ctx context.Context, reaction *app.Reaction) error
	Delete(ctx context.Context, targetID, authorID string) error
	GetByTarget(ctx context.Context, targetID, targetType string) ([]app.Reaction, error)
	GetByAuthor(ctx context.Context, authorID string, limit int) ([]app.Reaction, error)
	Exists(ctx context.Context, targetID, authorID string) (bool, error)
	CountByTarget(ctx context.Context, targetID, targetType string) (map[string]int, error)
	GetReactionTypes(ctx context.Context, targetID, targetType string) ([]string, error)
}

// Save creates a new reaction with the given parameters (legacy method)
func (rr *reactionRepository) Save(ctx context.Context, targetID, authorID, content string) error {
	if targetID == "" {
		return errors.NewValidationError("target ID is required")
	}

	if authorID == "" {
		return errors.NewValidationError("author ID is required")
	}

	if content == "" {
		return errors.NewValidationError("reaction content is required")
	}

	reaction := &app.Reaction{
		TargetID:  targetID,
		AuthorID:  authorID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return rr.SaveReaction(ctx, reaction)
}

// SaveReaction saves a complete reaction object
func (rr *reactionRepository) SaveReaction(ctx context.Context, reaction *app.Reaction) error {
	if reaction == nil {
		return errors.NewValidationError("reaction cannot be nil")
	}

	if reaction.TargetID == "" {
		return errors.NewValidationError("target ID is required")
	}

	if reaction.AuthorID == "" {
		return errors.NewValidationError("author ID is required")
	}

	if reaction.Content == "" {
		return errors.NewValidationError("reaction content is required")
	}

	now := time.Now()
	if reaction.CreatedAt.IsZero() {
		reaction.CreatedAt = now
	}
	reaction.UpdatedAt = now

	// Determine target type based on content or default to "post"
	targetType := "post" // Default assumption
	if len(reaction.Content) > 0 {
		// You could implement logic to determine target type
		// For now, we'll assume posts, but this could be enhanced
	}

	// Insert into main reactions table
	query := `INSERT INTO mingle.reactions (target_id, target_type, author_id, reaction_type, created_at)
			  VALUES (?, ?, ?, ?, ?)`

	err := rr.session.Query(query,
		reaction.TargetID,
		targetType,
		reaction.AuthorID,
		reaction.Content,
		reaction.CreatedAt,
	).WithContext(ctx).Exec()
	if err != nil {
		rr.logger.WithComponent("reaction-repository").Error("Failed to save reaction to main table",
			"target_id", reaction.TargetID,
			"author_id", reaction.AuthorID,
			"reaction_type", reaction.Content,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	// Insert into reactions_by_target table for efficient target queries
	targetQuery := `INSERT INTO mingle.reactions_by_target (target_id, target_type, reaction_type, author_id, created_at)
					VALUES (?, ?, ?, ?, ?)`

	err = rr.session.Query(targetQuery,
		reaction.TargetID,
		targetType,
		reaction.Content,
		reaction.AuthorID,
		reaction.CreatedAt,
	).WithContext(ctx).Exec()
	if err != nil {
		rr.logger.WithComponent("reaction-repository").Error("Failed to save reaction to target table",
			"target_id", reaction.TargetID,
			"author_id", reaction.AuthorID,
			"reaction_type", reaction.Content,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	rr.logger.WithComponent("reaction-repository").Info("Reaction saved successfully",
		"target_id", reaction.TargetID,
		"author_id", reaction.AuthorID,
		"reaction_type", reaction.Content,
	)

	return nil
}

// Delete removes a reaction
func (rr *reactionRepository) Delete(ctx context.Context, targetID, authorID string) error {
	if targetID == "" {
		return errors.NewValidationError("target ID is required")
	}

	if authorID == "" {
		return errors.NewValidationError("author ID is required")
	}

	// Check if reaction exists
	exists, err := rr.Exists(ctx, targetID, authorID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.NewNotFoundError("reaction not found")
	}

	// Default target type - in production this should be determined properly
	targetType := "post"

	// Delete from main reactions table
	query := `DELETE FROM mingle.reactions WHERE target_id = ? AND target_type = ? AND author_id = ?`
	err = rr.session.Query(query, targetID, targetType, authorID).WithContext(ctx).Exec()
	if err != nil {
		rr.logger.WithComponent("reaction-repository").Error("Failed to delete reaction from main table",
			"target_id", targetID,
			"author_id", authorID,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	// Delete from reactions_by_target table
	// Note: We need to get the reaction type first for proper deletion
	targetQuery := `DELETE FROM mingle.reactions_by_target WHERE target_id = ? AND target_type = ? AND author_id = ?`
	err = rr.session.Query(targetQuery, targetID, targetType, authorID).WithContext(ctx).Exec()
	if err != nil {
		rr.logger.WithComponent("reaction-repository").Error("Failed to delete reaction from target table",
			"target_id", targetID,
			"author_id", authorID,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	rr.logger.WithComponent("reaction-repository").Info("Reaction deleted successfully",
		"target_id", targetID,
		"author_id", authorID,
	)

	return nil
}

// GetByTarget retrieves reactions for a specific target
func (rr *reactionRepository) GetByTarget(ctx context.Context, targetID, targetType string) ([]app.Reaction, error) {
	if targetID == "" {
		return nil, errors.NewValidationError("target ID is required")
	}

	if targetType == "" {
		targetType = "post" // Default
	}

	var reactions []app.Reaction

	query := `SELECT reaction_type, author_id, created_at
			  FROM mingle.reactions_by_target WHERE target_id = ? AND target_type = ?`

	iter := rr.session.Query(query, targetID, targetType).WithContext(ctx).Iter()
	defer iter.Close()

	var reactionType, authorID string
	var createdAt time.Time

	for iter.Scan(&reactionType, &authorID, &createdAt) {
		reactions = append(reactions, app.Reaction{
			TargetID:  targetID,
			AuthorID:  authorID,
			Content:   reactionType,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		})
	}

	if err := iter.Close(); err != nil {
		rr.logger.WithComponent("reaction-repository").Error("Failed to get reactions by target",
			"target_id", targetID,
			"target_type", targetType,
			"error", err.Error(),
		)
		return nil, errors.NewDatabaseError(err)
	}

	rr.logger.WithComponent("reaction-repository").Debug("Reactions by target retrieved successfully",
		"target_id", targetID,
		"target_type", targetType,
		"reactions_count", len(reactions),
	)

	return reactions, nil
}

// GetByAuthor retrieves reactions by a specific author
func (rr *reactionRepository) GetByAuthor(ctx context.Context, authorID string, limit int) ([]app.Reaction, error) {
	if authorID == "" {
		return nil, errors.NewValidationError("author ID is required")
	}

	if limit <= 0 {
		limit = 50 // Default limit
	}

	var reactions []app.Reaction

	query := `SELECT target_id, target_type, reaction_type, created_at
			  FROM mingle.reactions WHERE author_id = ? LIMIT ?`

	iter := rr.session.Query(query, authorID, limit).WithContext(ctx).Iter()
	defer iter.Close()

	var targetID, targetType, reactionType string
	var createdAt time.Time

	for iter.Scan(&targetID, &targetType, &reactionType, &createdAt) {
		reactions = append(reactions, app.Reaction{
			TargetID:  targetID,
			AuthorID:  authorID,
			Content:   reactionType,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		})
	}

	if err := iter.Close(); err != nil {
		rr.logger.WithComponent("reaction-repository").Error("Failed to get reactions by author",
			"author_id", authorID,
			"error", err.Error(),
		)
		return nil, errors.NewDatabaseError(err)
	}

	rr.logger.WithComponent("reaction-repository").Debug("Reactions by author retrieved successfully",
		"author_id", authorID,
		"reactions_count", len(reactions),
	)

	return reactions, nil
}

// Exists checks if a reaction exists for a target by a specific author
func (rr *reactionRepository) Exists(ctx context.Context, targetID, authorID string) (bool, error) {
	if targetID == "" {
		return false, errors.NewValidationError("target ID is required")
	}

	if authorID == "" {
		return false, errors.NewValidationError("author ID is required")
	}

	var count int
	targetType := "post" // Default

	query := `SELECT COUNT(*) FROM mingle.reactions WHERE target_id = ? AND target_type = ? AND author_id = ?`

	err := rr.session.Query(query, targetID, targetType, authorID).WithContext(ctx).Scan(&count)
	if err != nil {
		rr.logger.WithComponent("reaction-repository").Error("Failed to check reaction existence",
			"target_id", targetID,
			"author_id", authorID,
			"error", err.Error(),
		)
		return false, errors.NewDatabaseError(err)
	}

	return count > 0, nil
}

// CountByTarget counts reactions by type for a specific target
func (rr *reactionRepository) CountByTarget(ctx context.Context, targetID, targetType string) (map[string]int, error) {
	if targetID == "" {
		return nil, errors.NewValidationError("target ID is required")
	}

	if targetType == "" {
		targetType = "post" // Default
	}

	reactionCounts := make(map[string]int)

	query := `SELECT reaction_type, COUNT(*) FROM mingle.reactions_by_target
			  WHERE target_id = ? AND target_type = ? GROUP BY reaction_type`

	iter := rr.session.Query(query, targetID, targetType).WithContext(ctx).Iter()
	defer iter.Close()

	var reactionType string
	var count int

	for iter.Scan(&reactionType, &count) {
		reactionCounts[reactionType] = count
	}

	if err := iter.Close(); err != nil {
		rr.logger.WithComponent("reaction-repository").Error("Failed to count reactions by target",
			"target_id", targetID,
			"target_type", targetType,
			"error", err.Error(),
		)
		return nil, errors.NewDatabaseError(err)
	}

	return reactionCounts, nil
}

// GetReactionTypes retrieves unique reaction types for a target
func (rr *reactionRepository) GetReactionTypes(ctx context.Context, targetID, targetType string) ([]string, error) {
	if targetID == "" {
		return nil, errors.NewValidationError("target ID is required")
	}

	if targetType == "" {
		targetType = "post" // Default
	}

	var reactionTypes []string

	query := `SELECT DISTINCT reaction_type FROM mingle.reactions_by_target
			  WHERE target_id = ? AND target_type = ?`

	iter := rr.session.Query(query, targetID, targetType).WithContext(ctx).Iter()
	defer iter.Close()

	var reactionType string

	for iter.Scan(&reactionType) {
		reactionTypes = append(reactionTypes, reactionType)
	}

	if err := iter.Close(); err != nil {
		rr.logger.WithComponent("reaction-repository").Error("Failed to get reaction types",
			"target_id", targetID,
			"target_type", targetType,
			"error", err.Error(),
		)
		return nil, errors.NewDatabaseError(err)
	}

	return reactionTypes, nil
}

func NewReactionRepository(session *gocql.Session) ReactionRepository {
	return &reactionRepository{
		session: session,
		logger:  logger.GetLogger(),
	}
}
