package graph

import (
	"context"
	_ "embed"

	"github.com/google/uuid"
	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	//go:embed cypher/create_comment.cypher
	createCommentCypher string

	//go:embed cypher/list_post_comments.cypher
	listPostComments string

	//go:embed cypher/update_comment.cypher
	updateCommentCypher string

	//go:embed cypher/delete_comment.cypher
	deleteCommentCypher string
)

type commentNeo4jRepository struct {
	query Neo4jQuerier[app.Comment]
}

func NewCommentRepository(driver neo4j.DriverWithContext) *commentNeo4jRepository {
	return &commentNeo4jRepository{
		query: Neo4jQuerier[app.Comment]{
			driver: driver,
		},
	}
}

// Save implements app.CommentRepository.
func (cr *commentNeo4jRepository) Save(ctx context.Context, authorId, postId, content string) (comment app.Comment, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return
	}

	params := struct {
		ID       string `json:"id"`
		Content  string `json:"content" validate:"required"`
		AuthorID string `json:"author_id" validate:"required"`
		PostID   string `json:"post_id" validate:"required"`
	}{
		ID:       id.String(),
		Content:  content,
		AuthorID: authorId,
		PostID:   postId,
	}

	return cr.query.Create(ctx, params, createCommentCypher)
}

// GetAll implements app.CommentRepository.
func (cr *commentNeo4jRepository) GetAll(ctx context.Context, id string) (posts []app.Comment, err error) {
	return cr.query.List(ctx, listPostComments, id)
}

// Update implements app.CommentRepository.
func (cr *commentNeo4jRepository) Update(ctx context.Context, commentId, content string) (comment app.Comment, err error) {
	params := struct {
		ID      string `json:"id"`
		Content string `json:"content" validate:"required"`
	}{
		ID:      commentId,
		Content: content,
	}

	return cr.query.Update(ctx, updateCommentCypher, params)
}

// Delete implements app.CommentRepository.
func (cr *commentNeo4jRepository) Delete(ctx context.Context, userId, commentId string) (err error) {
	return cr.query.Delete(ctx, deleteCommentCypher, map[string]any{
		"id":        commentId,
		"author_id": userId,
	})
}
