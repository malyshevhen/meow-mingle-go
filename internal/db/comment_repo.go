package db

import (
	"context"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
	"github.com/malyshEvhen/meow_mingle/pkg/logger"
)

type commentRepository struct {
	session *gocql.Session
	logger  *logger.Logger
}

// CommentRepository defines the interface for comment data operations
type CommentRepository interface {
	Save(ctx context.Context, authorId, postId, content string) (app.Comment, error)
	SaveComment(ctx context.Context, comment *app.Comment) error
	GetAll(ctx context.Context, id string) ([]app.Comment, error)
	GetByPost(ctx context.Context, postId string, limit int) ([]app.Comment, error)
	GetById(ctx context.Context, commentId string) (app.Comment, error)
	Update(ctx context.Context, commentId, content string) (app.Comment, error)
	Delete(ctx context.Context, userId, commentId string) error
	Exists(ctx context.Context, commentId string) (bool, error)
	CountByPost(ctx context.Context, postId string) (int, error)
}

// Save creates a new comment with the given parameters
func (cr *commentRepository) Save(ctx context.Context, authorId, postId, content string) (app.Comment, error) {
	if authorId == "" {
		return app.Comment{}, errors.NewValidationError("author ID is required")
	}

	if postId == "" {
		return app.Comment{}, errors.NewValidationError("post ID is required")
	}

	if content == "" {
		return app.Comment{}, errors.NewValidationError("content is required")
	}

	comment := app.Comment{
		ID:        uuid.New().String(),
		AuthorID:  authorId,
		PostID:    postId,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := cr.SaveComment(ctx, &comment); err != nil {
		return app.Comment{}, err
	}

	return comment, nil
}

// SaveComment saves a complete comment object
func (cr *commentRepository) SaveComment(ctx context.Context, comment *app.Comment) error {
	if comment == nil {
		return errors.NewValidationError("comment cannot be nil")
	}

	if comment.ID == "" {
		comment.ID = uuid.New().String()
	}

	if comment.AuthorID == "" {
		return errors.NewValidationError("author ID is required")
	}

	if comment.PostID == "" {
		return errors.NewValidationError("post ID is required")
	}

	if comment.Content == "" {
		return errors.NewValidationError("content is required")
	}

	now := time.Now()
	if comment.CreatedAt.IsZero() {
		comment.CreatedAt = now
	}
	comment.UpdatedAt = now

	// Insert into main comments table
	query := `
INSERT INTO mingle.comments
(
	id,
	post_id,
	author_id,
	content,
	created_at,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?)`

	err := cr.session.Query(query,
		comment.ID,
		comment.PostID,
		comment.AuthorID,
		comment.Content,
		comment.CreatedAt,
		comment.UpdatedAt,
	).WithContext(ctx).Exec()
	if err != nil {
		cr.logger.WithComponent("comment-repository").Error("Failed to save comment to main table",
			"comment_id", comment.ID,
			"post_id", comment.PostID,
			"author_id", comment.AuthorID,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	// Insert into comments_by_post table for efficient post comment queries
	postQuery := `
INSERT INTO mingle.comments_by_post
(
	post_id,
	created_at,
	comment_id,
	author_id,
	content,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?)`

	err = cr.session.Query(postQuery,
		comment.PostID,
		comment.CreatedAt,
		comment.ID,
		comment.AuthorID,
		comment.Content,
		comment.UpdatedAt,
	).WithContext(ctx).Exec()
	if err != nil {
		cr.logger.WithComponent("comment-repository").Error("Failed to save comment to post table",
			"comment_id", comment.ID,
			"post_id", comment.PostID,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	cr.logger.WithComponent("comment-repository").Info("Comment saved successfully",
		"comment_id", comment.ID,
		"post_id", comment.PostID,
		"author_id", comment.AuthorID,
	)

	return nil
}

// GetAll retrieves all comments for a post (legacy method for compatibility)
func (cr *commentRepository) GetAll(ctx context.Context, id string) ([]app.Comment, error) {
	return cr.GetByPost(ctx, id, 100) // Default limit of 100
}

// GetByPost retrieves comments for a specific post with limit
func (cr *commentRepository) GetByPost(ctx context.Context, postId string, limit int) ([]app.Comment, error) {
	if postId == "" {
		return nil, errors.NewValidationError("post ID is required")
	}

	if limit <= 0 {
		limit = 50 // Default limit
	}

	var comments []app.Comment

	query := `
SELECT
	comment_id,
	author_id,
	content,
	created_at,
	updated_at
FROM mingle.comments_by_post
WHERE post_id = ?
ORDER BY created_at DESC
LIMIT ?`

	iter := cr.session.Query(query, postId, limit).WithContext(ctx).Iter()
	defer iter.Close()

	var commentId, authorId, content string
	var createdAt, updatedAt time.Time

	for iter.Scan(&commentId, &authorId, &content, &createdAt, &updatedAt) {
		comments = append(comments, app.Comment{
			ID:        commentId,
			PostID:    postId,
			AuthorID:  authorId,
			Content:   content,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	if err := iter.Close(); err != nil {
		cr.logger.WithComponent("comment-repository").Error("Failed to get comments by post",
			"post_id", postId,
			"error", err.Error(),
		)
		return nil, errors.NewDatabaseError(err)
	}

	cr.logger.WithComponent("comment-repository").Debug("Comments by post retrieved successfully",
		"post_id", postId,
		"comments_count", len(comments),
	)

	return comments, nil
}

// GetById retrieves a comment by ID
func (cr *commentRepository) GetById(ctx context.Context, commentId string) (app.Comment, error) {
	if commentId == "" {
		return app.Comment{}, errors.NewValidationError("comment ID is required")
	}

	var comment app.Comment

	query := `
SELECT
	id,
	post_id,
	author_id,
	content,
	created_at,
	updated_at
FROM mingle.comments
WHERE id = ?`

	err := cr.session.Query(query, commentId).WithContext(ctx).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.AuthorID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		if err == gocql.ErrNotFound {
			cr.logger.WithComponent("comment-repository").Info("Comment not found",
				"comment_id", commentId,
			)
			return app.Comment{}, errors.NewNotFoundError("comment not found")
		}
		cr.logger.WithComponent("comment-repository").Error("Failed to get comment",
			"comment_id", commentId,
			"error", err.Error(),
		)
		return app.Comment{}, errors.NewDatabaseError(err)
	}

	cr.logger.WithComponent("comment-repository").Debug("Comment retrieved successfully",
		"comment_id", commentId,
	)

	return comment, nil
}

// Update updates a comment's content
func (cr *commentRepository) Update(ctx context.Context, commentId, content string) (app.Comment, error) {
	if commentId == "" {
		return app.Comment{}, errors.NewValidationError("comment ID is required")
	}

	if content == "" {
		return app.Comment{}, errors.NewValidationError("content is required")
	}

	// Check if comment exists and get current data
	currentComment, err := cr.GetById(ctx, commentId)
	if err != nil {
		return app.Comment{}, err
	}

	now := time.Now()

	// Update main comments table
	query := `
UPDATE mingle.comments
SET content = ?, updated_at = ?
WHERE id = ?`

	err = cr.session.Query(query, content, now, commentId).WithContext(ctx).Exec()
	if err != nil {
		cr.logger.WithComponent("comment-repository").Error("Failed to update comment in main table",
			"comment_id", commentId,
			"error", err.Error(),
		)
		return app.Comment{}, errors.NewDatabaseError(err)
	}

	// Update comments_by_post table
	postQuery := `
UPDATE mingle.comments_by_post
SET content = ?, updated_at = ?
WHERE post_id = ?
AND created_at = ?
AND comment_id = ?`

	err = cr.session.Query(postQuery, content, now, currentComment.PostID, currentComment.CreatedAt, commentId).WithContext(ctx).Exec()
	if err != nil {
		cr.logger.WithComponent("comment-repository").Error("Failed to update comment in post table",
			"comment_id", commentId,
			"post_id", currentComment.PostID,
			"error", err.Error(),
		)
		return app.Comment{}, errors.NewDatabaseError(err)
	}

	// Return updated comment
	updatedComment := currentComment
	updatedComment.Content = content
	updatedComment.UpdatedAt = now

	cr.logger.WithComponent("comment-repository").Info("Comment updated successfully",
		"comment_id", commentId,
	)

	return updatedComment, nil
}

// Delete removes a comment
func (cr *commentRepository) Delete(ctx context.Context, userId, commentId string) error {
	if userId == "" {
		return errors.NewValidationError("user ID is required")
	}

	if commentId == "" {
		return errors.NewValidationError("comment ID is required")
	}

	// Get comment details before deletion for authorization and cleanup
	comment, err := cr.GetById(ctx, commentId)
	if err != nil {
		return err
	}

	// Check if the user is authorized to delete this comment
	if comment.AuthorID != userId {
		return errors.NewForbiddenError()
	}

	// Delete from main comments table
	query := `DELETE FROM mingle.comments WHERE id = ?`
	err = cr.session.Query(query, commentId).WithContext(ctx).Exec()
	if err != nil {
		cr.logger.WithComponent("comment-repository").Error("Failed to delete comment from main table",
			"comment_id", commentId,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	// Delete from comments_by_post table
	postQuery := `
DELETE FROM mingle.comments_by_post
WHERE post_id = ?
AND created_at = ?
AND comment_id = ?`

	err = cr.session.Query(postQuery, comment.PostID, comment.CreatedAt, commentId).WithContext(ctx).Exec()
	if err != nil {
		cr.logger.WithComponent("comment-repository").Error("Failed to delete comment from post table",
			"comment_id", commentId,
			"post_id", comment.PostID,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	cr.logger.WithComponent("comment-repository").Info("Comment deleted successfully",
		"comment_id", commentId,
		"user_id", userId,
	)

	return nil
}

// Exists checks if a comment exists
func (cr *commentRepository) Exists(ctx context.Context, commentId string) (bool, error) {
	if commentId == "" {
		return false, errors.NewValidationError("comment ID is required")
	}

	var count int
	query := `SELECT COUNT(*) FROM mingle.comments WHERE id = ?`

	err := cr.session.Query(query, commentId).WithContext(ctx).Scan(&count)
	if err != nil {
		cr.logger.WithComponent("comment-repository").Error("Failed to check comment existence",
			"comment_id", commentId,
			"error", err.Error(),
		)
		return false, errors.NewDatabaseError(err)
	}

	return count > 0, nil
}

// CountByPost counts comments for a specific post
func (cr *commentRepository) CountByPost(ctx context.Context, postId string) (int, error) {
	if postId == "" {
		return 0, errors.NewValidationError("post ID is required")
	}

	var count int
	query := `SELECT COUNT(*) FROM mingle.comments_by_post WHERE post_id = ?`

	err := cr.session.Query(query, postId).WithContext(ctx).Scan(&count)
	if err != nil {
		cr.logger.WithComponent("comment-repository").Error("Failed to count comments by post",
			"post_id", postId,
			"error", err.Error(),
		)
		return 0, errors.NewDatabaseError(err)
	}

	return count, nil
}

func NewCommentRepository(session *gocql.Session) CommentRepository {
	return &commentRepository{
		session: session,
		logger:  logger.GetLogger(),
	}
}
