package db

import (
	"context"
	_ "embed"

	"github.com/google/uuid"
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

type ICommentRepository interface {
	CreateComment(ctx context.Context, params CreateCommentParams) (comment Comment, err error)
	CreateCommentLike(ctx context.Context, params CreateCommentLikeParams) (err error)
	ListPostComments(ctx context.Context, id string) (posts []Comment, err error)
	UpdateComment(ctx context.Context, params UpdateCommentParams) (comment Comment, err error)
	DeleteComment(ctx context.Context, userId, commentId string) (err error)
	DeleteCommentLike(ctx context.Context, params DeleteCommentLikeParams) error
}

type CommentNeo4jRepository struct {
	Neo4jRepository[Comment]
}

func NewCommentRepository(driver neo4j.DriverWithContext) *CommentNeo4jRepository {
	return &CommentNeo4jRepository{
		Neo4jRepository: Neo4jRepository[Comment]{
			driver: driver,
		},
	}
}

func (cr *CommentNeo4jRepository) CreateComment(ctx context.Context, params CreateCommentParams) (comment Comment, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return
	}
	params.ID = id.String()

	return cr.Create(ctx, params, createCommentCypher)
}

func (cr *CommentNeo4jRepository) CreateCommentLike(ctx context.Context, params CreateCommentLikeParams) (err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return
	}
	params.ID = id.String()

	return cr.Write(ctx, createLikeOnCommentCypher, params)
}

func (cr *CommentNeo4jRepository) ListPostComments(ctx context.Context, id string) (posts []Comment, err error) {
	return cr.List(ctx, listPostComments, id)
}

func (cr *CommentNeo4jRepository) UpdateComment(ctx context.Context, params UpdateCommentParams) (comment Comment, err error) {
	return cr.Update(ctx, updateCommentCypher, params)
}

func (cr *CommentNeo4jRepository) DeleteComment(ctx context.Context, userId, commentId string) (err error) {
	return cr.Delete(ctx, deleteCommentCypher, map[string]any{
		"id":        commentId,
		"author_id": userId,
	})
}

func (cr *CommentNeo4jRepository) DeleteCommentLike(ctx context.Context, params DeleteCommentLikeParams) error {
	return cr.Delete(ctx, deleteCommentLikeCypher, params)
}
