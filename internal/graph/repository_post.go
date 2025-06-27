package graph

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

	//go:embed cypher/update_post.cypher
	updatePostCypher string

	//go:embed cypher/delete_post.cypher
	deletePostCypher string

	//go:embed cypher/list_user_posts.cypher
	listUserPostsCypher string

	//go:embed cypher/list_feed.cypher
	listFeed string
)

type postNeo4jRepository struct {
	query Neo4jQuerier[app.Post]
}

func NewPostRepository(driver neo4j.DriverWithContext) *postNeo4jRepository {
	return &postNeo4jRepository{
		query: Neo4jQuerier[app.Post]{
			driver: driver,
		},
	}
}

// Save implements app.PostRepository.
func (pr *postNeo4jRepository) Save(ctx context.Context, authorId, content string) (post app.Post, err error) {
	var id uuid.UUID

	id, err = uuid.NewRandom()
	if err != nil {
		return
	}

	params := struct {
		ID       string `json:"id"`
		Content  string `json:"content"`
		AuthorID string `json:"author_id"`
	}{
		ID:       id.String(),
		Content:  content,
		AuthorID: authorId,
	}

	return pr.query.Create(ctx, params, createPostCypher)
}

// Get implements app.PostRepository.
func (pr *postNeo4jRepository) Get(ctx context.Context, postId string) (post app.Post, err error) {
	return pr.query.Retrieve(ctx, getPostCypher, map[string]any{"id": postId})
}

// Feed implements app.PostRepository.
func (pr *postNeo4jRepository) Feed(ctx context.Context, userId string) (feed []app.Post, err error) {
	return pr.query.List(ctx, listFeed, userId)
}

// List implements app.PostRepository.
func (pr *postNeo4jRepository) List(ctx context.Context, profileId string) (posts []app.Post, err error) {
	return pr.query.List(ctx, listUserPostsCypher, profileId)
}

// Update implements app.PostRepository.
func (pr *postNeo4jRepository) Update(ctx context.Context, postId, content string) (post app.Post, err error) {
	params := struct {
		ID      string `json:"id"`
		Content string `json:"content"`
	}{
		ID:      postId,
		Content: content,
	}

	return pr.query.Update(ctx, updatePostCypher, params)
}

// Delete implements app.PostRepository.
func (pr *postNeo4jRepository) Delete(ctx context.Context, postId string) error {
	return pr.query.Delete(ctx, deletePostCypher, map[string]any{"id": postId})
}
