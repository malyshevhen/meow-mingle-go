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

type PostRepository struct {
	Reposytory[Post]
}

func NewPostRepository(driver neo4j.DriverWithContext) *PostRepository {
	return &PostRepository{
		Reposytory: Reposytory[Post]{
			driver: driver,
		},
	}
}

func (pr *PostRepository) CreatePost(ctx context.Context, params CreatePostParams) (post Post, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return
	}
	params.ID = id.String()

	return pr.Create(ctx, params, createPostCypher)
}

func (pr *PostRepository) CreatePostLike(ctx context.Context, params CreatePostLikeParams) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	params.ID = id.String()

	return pr.Write(ctx, createLikeOnPostCypher, params)
}

func (pr *PostRepository) GetPost(ctx context.Context, id string) (post Post, err error) {
	return pr.GetById(ctx, getPostCypher, id)
}

func (pr *PostRepository) GetFeed(ctx context.Context, userId string) (feed []Post, err error) {
	return pr.List(ctx, listFeed, userId)
}

func (pr *PostRepository) ListUserPosts(ctx context.Context, userId string) (posts []Post, err error) {
	return pr.List(ctx, listUserPostsCypher, userId)
}

func (pr *PostRepository) UpdatePost(ctx context.Context, params UpdatePostParams) (post Post, err error) {
	return pr.Update(ctx, updatePostCypher, params)
}

func (pr *PostRepository) DeletePost(ctx context.Context, userId, postId string) error {
	return pr.Delete(ctx, deletePostCypher, map[string]interface{}{
		"id":        postId,
		"author_id": userId,
	})
}

func (pr *PostRepository) DeletePostLike(ctx context.Context, params DeletePostLikeParams) error {
	return pr.Delete(ctx, deletePostLikeCypher, params)
}
