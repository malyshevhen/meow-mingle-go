package db

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

	//go:embed cypher/create_like_on_comment.cypher
	createLikeOnCommentCypher string

	//go:embed cypher/update_comment.cypher
	updateCommentCypher string

	//go:embed cypher/delete_comment.cypher
	deleteCommentCypher string

	//go:embed cypher/delete_comment_like.cypher
	deleteCommentLikeCypher string

	//go:embed cypher/list_post_comments.cypher
	listPostComments string
)

type commentNeo4jRepository struct {
	Neo4jRepository[app.Comment]
}

func NewCommentRepository(driver neo4j.DriverWithContext) *commentNeo4jRepository {
	return &commentNeo4jRepository{
		Neo4jRepository: Neo4jRepository[app.Comment]{
			driver: driver,
		},
	}
}

func (cr *commentNeo4jRepository) CreateComment(ctx context.Context, params app.CreateCommentParams) (comment app.Comment, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return
	}
	params.ID = id.String()

	return cr.Create(ctx, params, createCommentCypher)
}

func (cr *commentNeo4jRepository) ListPostComments(ctx context.Context, id string) (posts []app.Comment, err error) {
	return cr.List(ctx, listPostComments, id)
}

func (cr *commentNeo4jRepository) UpdateComment(ctx context.Context, params app.UpdateCommentParams) (comment app.Comment, err error) {
	return cr.Update(ctx, updateCommentCypher, params)
}

func (cr *commentNeo4jRepository) DeleteComment(ctx context.Context, userId, commentId string) (err error) {
	return cr.Delete(ctx, deleteCommentCypher, map[string]any{
		"id":        commentId,
		"author_id": userId,
	})
}
