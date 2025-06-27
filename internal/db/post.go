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

type postRepository struct {
	session *gocql.Session
	logger  *logger.Logger
}

// PostRepository defines the interface for post data operations
type PostRepository interface {
	Save(ctx context.Context, authorId, content string) (app.Post, error)
	SavePost(ctx context.Context, post *app.Post) error
	Get(ctx context.Context, postId string) (app.Post, error)
	Feed(ctx context.Context, userId string) ([]app.Post, error)
	List(ctx context.Context, profileId string) ([]app.Post, error)
	Update(ctx context.Context, postId, content string) (app.Post, error)
	Delete(ctx context.Context, postId string) error
	Exists(ctx context.Context, postId string) (bool, error)
	GetByAuthor(ctx context.Context, authorId string, limit int) ([]app.Post, error)
}

// Save creates a new post with the given parameters
func (pr *postRepository) Save(ctx context.Context, authorId, content string) (app.Post, error) {
	if authorId == "" {
		return app.Post{}, errors.NewValidationError("author ID is required")
	}

	if content == "" {
		return app.Post{}, errors.NewValidationError("content is required")
	}

	post := app.Post{
		ID:        uuid.New().String(),
		AuthorID:  authorId,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := pr.SavePost(ctx, &post); err != nil {
		return app.Post{}, err
	}

	return post, nil
}

// SavePost saves a complete post object
func (pr *postRepository) SavePost(ctx context.Context, post *app.Post) error {
	if post == nil {
		return errors.NewValidationError("post cannot be nil")
	}

	if post.ID == "" {
		post.ID = uuid.New().String()
	}

	if post.AuthorID == "" {
		return errors.NewValidationError("author ID is required")
	}

	if post.Content == "" {
		return errors.NewValidationError("content is required")
	}

	now := time.Now()
	if post.CreatedAt.IsZero() {
		post.CreatedAt = now
	}
	post.UpdatedAt = now

	// Insert into main posts table
	query := `INSERT INTO mingle.posts (id, author_id, content, image_urls, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?)`

	var imageUrls []string // Empty list for now
	err := pr.session.Query(query,
		post.ID,
		post.AuthorID,
		post.Content,
		imageUrls,
		post.CreatedAt,
		post.UpdatedAt,
	).WithContext(ctx).Exec()
	if err != nil {
		pr.logger.WithComponent("post-repository").Error("Failed to save post to main table",
			"post_id", post.ID,
			"author_id", post.AuthorID,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	// Insert into posts_by_author table for efficient author queries
	authorQuery := `INSERT INTO mingle.posts_by_author (author_id, created_at, post_id, content, image_urls, updated_at)
					VALUES (?, ?, ?, ?, ?, ?)`

	err = pr.session.Query(authorQuery,
		post.AuthorID,
		post.CreatedAt,
		post.ID,
		post.Content,
		imageUrls,
		post.UpdatedAt,
	).WithContext(ctx).Exec()
	if err != nil {
		pr.logger.WithComponent("post-repository").Error("Failed to save post to author table",
			"post_id", post.ID,
			"author_id", post.AuthorID,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	pr.logger.WithComponent("post-repository").Info("Post saved successfully",
		"post_id", post.ID,
		"author_id", post.AuthorID,
	)

	return nil
}

// Get retrieves a post by ID
func (pr *postRepository) Get(ctx context.Context, postId string) (app.Post, error) {
	if postId == "" {
		return app.Post{}, errors.NewValidationError("post ID is required")
	}

	var post app.Post
	var imageUrls []string

	query := `SELECT id, author_id, content, image_urls, created_at, updated_at
			  FROM mingle.posts WHERE id = ?`

	err := pr.session.Query(query, postId).WithContext(ctx).Scan(
		&post.ID,
		&post.AuthorID,
		&post.Content,
		&imageUrls,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		if err == gocql.ErrNotFound {
			pr.logger.WithComponent("post-repository").Info("Post not found",
				"post_id", postId,
			)
			return app.Post{}, errors.NewNotFoundError("post not found")
		}
		pr.logger.WithComponent("post-repository").Error("Failed to get post",
			"post_id", postId,
			"error", err.Error(),
		)
		return app.Post{}, errors.NewDatabaseError(err)
	}

	pr.logger.WithComponent("post-repository").Debug("Post retrieved successfully",
		"post_id", postId,
	)

	return post, nil
}

// Feed retrieves posts for a user's feed (could be enhanced with following logic)
func (pr *postRepository) Feed(ctx context.Context, userId string) ([]app.Post, error) {
	if userId == "" {
		return nil, errors.NewValidationError("user ID is required")
	}

	// For now, return recent posts from user_feed table
	// In a full implementation, this would be populated by a feed generation service
	var posts []app.Post

	query := `SELECT post_id, author_id, content, image_urls, created_at
			  FROM mingle.user_feed WHERE user_id = ?
			  ORDER BY created_at DESC LIMIT 20`

	iter := pr.session.Query(query, userId).WithContext(ctx).Iter()
	defer iter.Close()

	var postId, authorId, content string
	var imageUrls []string
	var createdAt time.Time

	for iter.Scan(&postId, &authorId, &content, &imageUrls, &createdAt) {
		posts = append(posts, app.Post{
			ID:        postId,
			AuthorID:  authorId,
			Content:   content,
			CreatedAt: createdAt,
			UpdatedAt: createdAt, // Use created_at as fallback
		})
	}

	if err := iter.Close(); err != nil {
		pr.logger.WithComponent("post-repository").Error("Failed to get user feed",
			"user_id", userId,
			"error", err.Error(),
		)
		return nil, errors.NewDatabaseError(err)
	}

	pr.logger.WithComponent("post-repository").Debug("Feed retrieved successfully",
		"user_id", userId,
		"posts_count", len(posts),
	)

	return posts, nil
}

// List retrieves posts by author ID
func (pr *postRepository) List(ctx context.Context, profileId string) ([]app.Post, error) {
	return pr.GetByAuthor(ctx, profileId, 50) // Default limit of 50
}

// GetByAuthor retrieves posts by author with limit
func (pr *postRepository) GetByAuthor(ctx context.Context, authorId string, limit int) ([]app.Post, error) {
	if authorId == "" {
		return nil, errors.NewValidationError("author ID is required")
	}

	if limit <= 0 {
		limit = 20 // Default limit
	}

	var posts []app.Post

	query := `SELECT post_id, content, image_urls, created_at, updated_at
			  FROM mingle.posts_by_author WHERE author_id = ?
			  ORDER BY created_at DESC LIMIT ?`

	iter := pr.session.Query(query, authorId, limit).WithContext(ctx).Iter()
	defer iter.Close()

	var postId, content string
	var imageUrls []string
	var createdAt, updatedAt time.Time

	for iter.Scan(&postId, &content, &imageUrls, &createdAt, &updatedAt) {
		posts = append(posts, app.Post{
			ID:        postId,
			AuthorID:  authorId,
			Content:   content,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	if err := iter.Close(); err != nil {
		pr.logger.WithComponent("post-repository").Error("Failed to get posts by author",
			"author_id", authorId,
			"error", err.Error(),
		)
		return nil, errors.NewDatabaseError(err)
	}

	pr.logger.WithComponent("post-repository").Debug("Posts by author retrieved successfully",
		"author_id", authorId,
		"posts_count", len(posts),
	)

	return posts, nil
}

// Update updates a post's content
func (pr *postRepository) Update(ctx context.Context, postId, content string) (app.Post, error) {
	if postId == "" {
		return app.Post{}, errors.NewValidationError("post ID is required")
	}

	if content == "" {
		return app.Post{}, errors.NewValidationError("content is required")
	}

	// Check if post exists and get current data
	currentPost, err := pr.Get(ctx, postId)
	if err != nil {
		return app.Post{}, err
	}

	now := time.Now()

	// Update main posts table
	query := `UPDATE mingle.posts SET content = ?, updated_at = ? WHERE id = ?`
	err = pr.session.Query(query, content, now, postId).WithContext(ctx).Exec()
	if err != nil {
		pr.logger.WithComponent("post-repository").Error("Failed to update post in main table",
			"post_id", postId,
			"error", err.Error(),
		)
		return app.Post{}, errors.NewDatabaseError(err)
	}

	// Update posts_by_author table
	authorQuery := `UPDATE mingle.posts_by_author SET content = ?, updated_at = ?
					WHERE author_id = ? AND created_at = ? AND post_id = ?`
	err = pr.session.Query(authorQuery, content, now, currentPost.AuthorID, currentPost.CreatedAt, postId).WithContext(ctx).Exec()
	if err != nil {
		pr.logger.WithComponent("post-repository").Error("Failed to update post in author table",
			"post_id", postId,
			"author_id", currentPost.AuthorID,
			"error", err.Error(),
		)
		return app.Post{}, errors.NewDatabaseError(err)
	}

	// Return updated post
	updatedPost := currentPost
	updatedPost.Content = content
	updatedPost.UpdatedAt = now

	pr.logger.WithComponent("post-repository").Info("Post updated successfully",
		"post_id", postId,
	)

	return updatedPost, nil
}

// Delete removes a post
func (pr *postRepository) Delete(ctx context.Context, postId string) error {
	if postId == "" {
		return errors.NewValidationError("post ID is required")
	}

	// Get post details before deletion
	post, err := pr.Get(ctx, postId)
	if err != nil {
		return err
	}

	// Delete from main posts table
	query := `DELETE FROM mingle.posts WHERE id = ?`
	err = pr.session.Query(query, postId).WithContext(ctx).Exec()
	if err != nil {
		pr.logger.WithComponent("post-repository").Error("Failed to delete post from main table",
			"post_id", postId,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	// Delete from posts_by_author table
	authorQuery := `DELETE FROM mingle.posts_by_author
					WHERE author_id = ? AND created_at = ? AND post_id = ?`
	err = pr.session.Query(authorQuery, post.AuthorID, post.CreatedAt, postId).WithContext(ctx).Exec()
	if err != nil {
		pr.logger.WithComponent("post-repository").Error("Failed to delete post from author table",
			"post_id", postId,
			"author_id", post.AuthorID,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	pr.logger.WithComponent("post-repository").Info("Post deleted successfully",
		"post_id", postId,
	)

	return nil
}

// Exists checks if a post exists
func (pr *postRepository) Exists(ctx context.Context, postId string) (bool, error) {
	if postId == "" {
		return false, errors.NewValidationError("post ID is required")
	}

	var count int
	query := `SELECT COUNT(*) FROM mingle.posts WHERE id = ?`

	err := pr.session.Query(query, postId).WithContext(ctx).Scan(&count)
	if err != nil {
		pr.logger.WithComponent("post-repository").Error("Failed to check post existence",
			"post_id", postId,
			"error", err.Error(),
		)
		return false, errors.NewDatabaseError(err)
	}

	return count > 0, nil
}

func NewPostRepository(session *gocql.Session) PostRepository {
	return &postRepository{
		session: session,
		logger:  logger.GetLogger(),
	}
}
