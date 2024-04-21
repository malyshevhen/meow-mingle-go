package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/errors"
	"github.com/malyshEvhen/meow_mingle/internal/mock"
	"github.com/malyshEvhen/meow_mingle/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleCreatePost(t *testing.T) {
	store := &mock.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /posts",
		MiddlewareChain(
			handleCreatePost(store),
			LoggerMW,
			ErrorHandler,
			fakeAuth(1),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	validParams := db.CreatePostParams{
		Content: "Hello world",
	}

	t.Run("returning 201 on successful post creation", func(t *testing.T) {
		post := db.Post{
			ID:       1,
			Content:  validParams.Content,
			AuthorID: 1,
		}
		store.SetPost(post)
		store.SetError(nil)

		paramsBytes, _ := json.Marshal(validParams)
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s/posts", server.URL),
			strings.NewReader(string(paramsBytes)),
		)
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, 201, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		postResp, err := utils.Unmarshal[db.Post](body)
		assert.NoError(t, err)

		assert.Equal(t, post, postResp)
	})

	t.Run("returning 400 on invalid post params", func(t *testing.T) {
		invalidParams := db.CreatePostParams{
			Content: "",
		}

		paramsBytes, _ := json.Marshal(invalidParams)
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s/posts", server.URL),
			strings.NewReader(string(paramsBytes)),
		)
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, 400, res.StatusCode)
	})

	t.Run("returning 500 on db error", func(t *testing.T) {
		store.SetError(errors.NewDatabaseError(fmt.Errorf("db error")))

		paramsBytes, _ := json.Marshal(validParams)
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s/posts", server.URL),
			strings.NewReader(string(paramsBytes)),
		)
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, 500, res.StatusCode)
	})
}

func TestHandleCreatePostUnauthorized(t *testing.T) {
	store := &mock.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /posts",
		MiddlewareChain(handleCreatePost(store),
			LoggerMW,
			ErrorHandler,
			WithJWTAuth(store),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	validParams := db.CreatePostParams{
		Content: "Hello world",
	}

	t.Run("returning 401 on invalid auth", func(t *testing.T) {
		paramsBytes, _ := json.Marshal(validParams)
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s/posts", server.URL),
			strings.NewReader(string(paramsBytes)),
		)
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, 401, res.StatusCode)
	})
}

func TestHandleGetUserPosts(t *testing.T) {
	store := &mock.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}/posts",
		MiddlewareChain(
			handleGetUserPosts(store),
			LoggerMW,
			ErrorHandler,
			fakeAuth(1),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returning 200 on successful get posts", func(t *testing.T) {
		post := db.ListUserPostsRow{
			ID:       1,
			AuthorID: 1,
			Content:  "Test post",
		}
		store.AddListUserPostsRows(post)
		store.SetError(nil)

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/1/posts", server.URL), nil)
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, 200, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		postsResp, err := utils.Unmarshal[[]db.ListUserPostsRow](body)
		assert.NoError(t, err)
		assert.Equal(t, post, postsResp[0])

	})

	t.Run("returning 404 if no posts found for user", func(t *testing.T) {
		store.SetError(errors.NewNotFoundError("no posts found"))

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/1/posts", server.URL), nil)
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, 404, res.StatusCode)
	})

	t.Run("returning 500 on db error", func(t *testing.T) {
		store.SetError(errors.NewDatabaseError(fmt.Errorf("db error")))

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/1/posts", server.URL), nil)
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, 500, res.StatusCode)
	})
}

func TestHandleGetUserPostsUnauthorized(t *testing.T) {
	store := &mock.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}/posts",
		MiddlewareChain(
			handleGetUserPosts(store),
			LoggerMW,
			ErrorHandler,
			WithJWTAuth(store),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returning 401 if unauthorized", func(t *testing.T) {
		post := db.ListUserPostsRow{
			ID:       1,
			AuthorID: 1,
			Content:  "Test post",
		}
		store.AddListUserPostsRows(post)

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/1/posts", server.URL), nil)
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, 401, res.StatusCode)
	})
}

func TestHandleGetPostsById(t *testing.T) {
	store := &mock.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /posts/{id}",
		MiddlewareChain(
			handleGetPostsById(store),
			LoggerMW,
			ErrorHandler,
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returns 200 and post if post found", func(t *testing.T) {
		post := samplePostRow()
		store.SetGetPostRow(post)

		url := fmt.Sprintf("%s/posts/%d", server.URL, post.ID)
		res, err := http.Get(url)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		gotPost, err := utils.Unmarshal[db.GetPostRow](body)
		require.NoError(t, err)

		require.Equal(t, post, gotPost)
	})

	t.Run("returns 404 if post not found", func(t *testing.T) {
		store.SetError(errors.NewNotFoundError("post not found"))

		url := fmt.Sprintf("%s/posts/%d", server.URL, 123)
		res, err := http.Get(url)
		require.NoError(t, err)

		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("returns 500 on unexpected error", func(t *testing.T) {
		store.SetError(errors.NewDatabaseError(fmt.Errorf("db error")))

		url := fmt.Sprintf("%s/posts/%d", server.URL, 123)
		res, err := http.Get(url)
		require.NoError(t, err)

		require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}

func TestHandleUpdatePostsById(t *testing.T) {
	store := &mock.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /posts/{id}",
		MiddlewareChain(
			handleUpdatePostsById(store),
			LoggerMW,
			ErrorHandler,
			fakeAuth(1),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returns 200 and updated post if successful", func(t *testing.T) {
		post := samplePost()
		store.SetPost(post)

		postBytes, err := json.Marshal(db.UpdatePostParams{Content: post.Content})
		require.NoError(t, err)

		url := fmt.Sprintf("%s/posts/%d", server.URL, post.ID)
		req, err := http.NewRequest("PUT", url, strings.NewReader(string(postBytes)))
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		gotPost, err := utils.Unmarshal[db.Post](responseBody)
		require.NoError(t, err)

		require.Equal(t, post, gotPost)
	})

	t.Run("returns 404 if post not found", func(t *testing.T) {
		store.SetError(errors.NewNotFoundError("post not found"))

		url := fmt.Sprintf("%s/posts/123", server.URL)
		body := `{"content": "Updated content"}`

		req, err := http.NewRequest("PUT", url, strings.NewReader(body))
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func TestHandleDeletePostsById(t *testing.T) {
	store := &mock.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /posts/{id}",
		MiddlewareChain(
			handleDeletePostsById(store),
			LoggerMW,
			ErrorHandler,
			fakeAuth(1),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returns 204 if post deleted successfully", func(t *testing.T) {
		store.SetError(nil)

		url := fmt.Sprintf("%s/posts/1", server.URL)
		req, err := http.NewRequest("DELETE", url, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusNoContent, res.StatusCode)
	})

	t.Run("returns 404 if post not found", func(t *testing.T) {
		store.SetError(errors.NewNotFoundError("post not found"))

		url := fmt.Sprintf("%s/posts/1", server.URL)
		req, err := http.NewRequest("DELETE", url, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("returns 500 on unexpected error", func(t *testing.T) {
		store.SetError(errors.NewDatabaseError(fmt.Errorf("db error")))

		url := fmt.Sprintf("%s/posts/1", server.URL)
		req, err := http.NewRequest("DELETE", url, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}

func TestHandleLikePost(t *testing.T) {
	store := &mock.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /posts/{id}/like",
		MiddlewareChain(
			handleLikePost(store),
			LoggerMW,
			ErrorHandler,
			fakeAuth(1),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returns 204 if post liked successfully", func(t *testing.T) {
		store.SetError(nil)

		url := fmt.Sprintf("%s/posts/1/like", server.URL)
		req, err := http.NewRequest("POST", url, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusNoContent, res.StatusCode)
	})

	t.Run("returns 404 if post not found", func(t *testing.T) {
		store.SetError(errors.NewNotFoundError("post not found"))

		url := fmt.Sprintf("%s/posts/1/like", server.URL)
		req, err := http.NewRequest("POST", url, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("returns 401 if unauthorized", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("POST /posts/{id}/like",
			MiddlewareChain(
				handleLikePost(store),
				LoggerMW,
				ErrorHandler,
				WithJWTAuth(store),
			),
		)

		server := httptest.NewServer(mux)
		defer server.Close()

		url := fmt.Sprintf("%s/posts/1/like", server.URL)
		req, err := http.NewRequest("POST", url, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("returns 400 if invalid post ID", func(t *testing.T) {
		url := fmt.Sprintf("%s/posts/invalid/like", server.URL)
		req, err := http.NewRequest("POST", url, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func TestHandleRemoveLikeFromPost(t *testing.T) {
	store := &mock.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /posts/{id}/like",
		MiddlewareChain(
			handleRemoveLikeFromPost(store),
			LoggerMW,
			ErrorHandler,
			fakeAuth(1),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returns 204 if like removed successfully", func(t *testing.T) {
		store.SetError(nil)

		url := fmt.Sprintf("%s/posts/1/like", server.URL)
		req, err := http.NewRequest("DELETE", url, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusNoContent, res.StatusCode)
	})

	t.Run("returns 404 if post not found", func(t *testing.T) {
		store.SetError(errors.NewNotFoundError("post not found"))

		url := fmt.Sprintf("%s/posts/1/like", server.URL)
		req, err := http.NewRequest("DELETE", url, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("returns 401 if unauthorized", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("DELETE /posts/{id}/like",
			MiddlewareChain(
				handleRemoveLikeFromPost(store),
				LoggerMW,
				ErrorHandler,
				WithJWTAuth(store),
			),
		)

		server := httptest.NewServer(mux)
		defer server.Close()

		url := fmt.Sprintf("%s/posts/1/like", server.URL)
		req, err := http.NewRequest("DELETE", url, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("returns 400 if invalid post ID", func(t *testing.T) {
		url := fmt.Sprintf("%s/posts/invalid/like", server.URL)
		req, err := http.NewRequest("DELETE", url, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func samplePost() db.Post {
	return db.Post{
		ID:       1,
		AuthorID: 1,
		Content:  "Hello world",
	}
}

func samplePostRow() db.GetPostRow {
	return db.GetPostRow{
		ID:       1,
		AuthorID: 1,
		Content:  "Hello world",
	}
}
