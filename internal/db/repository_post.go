package db

import (
	"context"
	_ "embed"

	"github.com/google/uuid"
	"github.com/malyshEvhen/meow_mingle/internal/app"
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

type postNeo4jRepository struct {
	Neo4jRepository[app.Post]
}

func NewPostRepository(driver neo4j.DriverWithContext) *postNeo4jRepository {
	return &postNeo4jRepository{
		Neo4jRepository: Neo4jRepository[app.Post]{
			driver: driver,
		},
	}
}

func (pr *postNeo4jRepository) CreatePost(ctx context.Context, params app.CreatePostParams) (post app.Post, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return
	}
	params.ID = id.String()

	return pr.Create(ctx, params, createPostCypher)
}

func (pr *postNeo4jRepository) GetPost(ctx context.Context, id string) (post app.Post, err error) {
	return pr.Retrieve(ctx, getPostCypher, map[string]any{
		"id": id,
	})
}

func (pr *postNeo4jRepository) GetFeed(ctx context.Context, userId string) (feed []app.Post, err error) {
	return pr.List(ctx, listFeed, userId)
}

func (pr *postNeo4jRepository) ListUserPosts(ctx context.Context, userId string) (posts []app.Post, err error) {
	return pr.List(ctx, listUserPostsCypher, userId)
}

func (pr *postNeo4jRepository) UpdatePost(ctx context.Context, params app.UpdatePostParams) (post app.Post, err error) {
	return pr.Update(ctx, updatePostCypher, params)
}

func (pr *postNeo4jRepository) DeletePost(ctx context.Context, userId, postId string) error {
	return pr.Delete(ctx, deletePostCypher, map[string]any{
		"id":        postId,
		"author_id": userId,
	})
}
