package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
	"github.com/malyshEvhen/meow_mingle/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreateComment(t *testing.T) {
	store := &db.MockStore{}

	validParams := db.CreateCommentParams{
		PostID:  1,
		Content: "Test Comment",
	}

	invalidParams := db.CreateCommentParams{
		PostID:  0,
		Content: "",
	}

	comment := db.Comment{
		ID:       1,
		AuthorID: 1,
		Content:  "Test Comment",
	}
	store.SetComment(comment)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /posts/{id}/comments",
		MiddlewareChain(
			handleCreateComment(store),
			LoggerMW,
			ErrorHandler,
			fakeAuth(1),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returns 201 and created comment if valid", func(t *testing.T) {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/posts/1/comments", server.URL), reqBodyOf(validParams))
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err, "read response body")

		if status := resp.StatusCode; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		user, err := Unmarshal[db.Comment](body)

		assert.NoErrorf(t, err, "unmarshal response body")
		assert.Truef(t, reflect.DeepEqual(user, comment), "handler returned wrong body: got %v want %v", user, comment)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("returns 400 if invalid post ID", func(t *testing.T) {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/posts/invalid/comments", server.URL), reqBodyOf(validParams))
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("returns 400 if invalid content", func(t *testing.T) {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/posts/1/comments", server.URL), reqBodyOf(invalidParams))
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("returns 500 if error creating comment", func(t *testing.T) {
		store.SetError(errors.NewInternalServerError(fmt.Errorf("error creating comment")))

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/posts/1/comments", server.URL), reqBodyOf(validParams))
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestHandleCreateCommentUnauthenticated(t *testing.T) {
	store := &db.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /posts/{id}/comments",
		MiddlewareChain(
			handleCreateComment(store),
			LoggerMW,
			ErrorHandler,
			WithJWTAuth(store),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returns 401 if unauthenticated", func(t *testing.T) {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/posts/1/comments", server.URL), nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
