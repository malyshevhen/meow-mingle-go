package db

import (
	"context"
	_ "embed"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	//go:embed cypher/create_post.cypher
	createPostCypher string

	//go:embed cypher/match_post_by_id.cypher
	getPostCypher string

	//go:embed cypher/create_like_on_post.cypher
	createLikeOnPostCypher string

	//go:embed cypher/update_post.cypher
	updatePostCypher string

	//go:embed cypher/delete_post.cypher
	deletePostCypher string

	//go:embed cypher/delete_post_like.cypher
	deletePostLikeCypher string

	//go:embed cypher/list_user_posts.cypher
	listUserPostsCypher string

	//go:embed cypher/list_feed.cypher
	listFeed string
)

type IPostRepository interface {
	CreatePost(ctx context.Context, params CreatePostParams) (post Post, err error)
	CreatePostLike(ctx context.Context, params CreatePostLikeParams) error
	GetPost(ctx context.Context, id string) (post Post, err error)
	GetFeed(ctx context.Context, userId string) (feed []Post, err error)
	ListUserPosts(ctx context.Context, userId string) (posts []Post, err error)
	UpdatePost(ctx context.Context, params UpdatePostParams) (post Post, err error)
	DeletePost(ctx context.Context, userId, postId string) error
	DeletePostLike(ctx context.Context, params DeletePostLikeParams) error
}

type PostNeo4jRepository struct {
	Neo4jRepository[Post]
}

func NewPostRepository(driver neo4j.DriverWithContext) *PostNeo4jRepository {
	return &PostNeo4jRepository{
		Neo4jRepository: Neo4jRepository[Post]{
			driver: driver,
		},
	}
}

func (pr *PostNeo4jRepository) CreatePost(ctx context.Context, params CreatePostParams) (post Post, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return
	}
	params.ID = id.String()

	return pr.Create(ctx, params, createPostCypher)
}

func (pr *PostNeo4jRepository) CreatePostLike(ctx context.Context, params CreatePostLikeParams) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	params.ID = id.String()

	return pr.Write(ctx, createLikeOnPostCypher, params)
}

func (pr *PostNeo4jRepository) GetPost(ctx context.Context, id string) (post Post, err error) {
	return pr.Retrieve(ctx, getPostCypher, map[string]any{
		"id": id,
	})
}

func (pr *PostNeo4jRepository) GetFeed(ctx context.Context, userId string) (feed []Post, err error) {
	return pr.List(ctx, listFeed, userId)
}

func (pr *PostNeo4jRepository) ListUserPosts(ctx context.Context, userId string) (posts []Post, err error) {
	return pr.List(ctx, listUserPostsCypher, userId)
}

func (pr *PostNeo4jRepository) UpdatePost(ctx context.Context, params UpdatePostParams) (post Post, err error) {
	return pr.Update(ctx, updatePostCypher, params)
}

func (pr *PostNeo4jRepository) DeletePost(ctx context.Context, userId, postId string) error {
	return pr.Delete(ctx, deletePostCypher, map[string]any{
		"id":        postId,
		"author_id": userId,
	})
}

func (pr *PostNeo4jRepository) DeletePostLike(ctx context.Context, params DeletePostLikeParams) error {
	return pr.Delete(ctx, deletePostLikeCypher, params)
}
