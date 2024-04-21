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

type Comments []db.GetCommentRow

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
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s/posts/1/comments", server.URL),
			reqBodyOf(validParams),
		)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err, "read response body")

		if status := resp.StatusCode; status != http.StatusCreated {
			t.Errorf(
				"handler returned wrong status code: got %v want %v",
				status,
				http.StatusCreated,
			)
		}

		user, err := Unmarshal[db.Comment](body)

		assert.NoErrorf(t, err, "unmarshal response body")
		assert.Truef(
			t,
			reflect.DeepEqual(user, comment),
			"handler returned wrong body: got %v want %v",
			user,
			comment,
		)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("returns 400 if invalid post ID", func(t *testing.T) {
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s/posts/invalid/comments", server.URL),
			reqBodyOf(validParams),
		)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("returns 400 if invalid content", func(t *testing.T) {
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s/posts/1/comments", server.URL),
			reqBodyOf(invalidParams),
		)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("returns 500 if error creating comment", func(t *testing.T) {
		store.SetError(errors.NewInternalServerError(fmt.Errorf("error creating comment")))

		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s/posts/1/comments", server.URL),
			reqBodyOf(validParams),
		)
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

func TestHandleGetComments(t *testing.T) {
	store := &db.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /posts/{id}/comments",
		MiddlewareChain(
			handleGetComments(store),
			LoggerMW,
			ErrorHandler,
			fakeAuth(1),
		),
	)

	commentRow := db.ListPostCommentsRow{
		ID:       1,
		AuthorID: 1,
		Content:  "Test Comment",
	}

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returns 200 and comments if valid", func(t *testing.T) {
		store.AddComments(commentRow)

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/posts/1/comments", server.URL), nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")
		defer resp.Body.Close()

		assert.Equalf(t, http.StatusOK, resp.StatusCode,
			"handler returned wrong status code: got %v want %v",
			resp.StatusCode, http.StatusOK)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		respComments, err := Unmarshal[Comments](body)
		assert.NoError(t, err)
		assert.NotEmpty(t, respComments)
		assert.Equalf(t, db.GetCommentRow(commentRow), respComments[0],
			"handler returned unexpected body: got %v want %v",
			respComments[0], commentRow)
	})

	t.Run("returns 404 if post not found", func(t *testing.T) {
		store.SetError(errors.NewNotFoundError("post not found"))

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/posts/1/comments", server.URL), nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")
		defer resp.Body.Close()

		assert.Equalf(t, http.StatusNotFound, resp.StatusCode,
			"handler returned wrong status code: got %v want %v",
			resp.StatusCode, http.StatusOK)
	})

	t.Run("returns 500 on unexpected error", func(t *testing.T) {
		store.SetError(errors.NewInternalServerError(fmt.Errorf("unexpected error")))

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/posts/1/comments", server.URL), nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")
		defer resp.Body.Close()

		assert.Equalf(t, http.StatusInternalServerError, resp.StatusCode,
			"handler returned wrong status code: got %v want %v",
			resp.StatusCode, http.StatusOK)
	})
}

func TestHandleUpdateComment(t *testing.T) {
	store := &db.MockStore{}

	validParams := db.UpdateCommentParams{
		ID:      1,
		Content: "Updated Comment",
	}

	invalidParams := db.UpdateCommentParams{
		ID:      0,
		Content: "",
	}

	updatedComment := db.Comment{
		ID:       1,
		AuthorID: 1,
		Content:  "Updated Comment",
	}
	store.SetComment(updatedComment)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /comments/{id}",
		MiddlewareChain(
			handleUpdateComments(store),
			LoggerMW,
			ErrorHandler,
			fakeAuth(1),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returns 200 and updated comment if valid", func(t *testing.T) {
		req, err := http.NewRequest(
			"PUT",
			fmt.Sprintf("%s/comments/1", server.URL),
			reqBodyOf(validParams),
		)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err, "read response body")

		comment, err := Unmarshal[db.Comment](body)
		assert.NoErrorf(t, err, "unmarshal response body")
		assert.Truef(
			t,
			reflect.DeepEqual(comment, updatedComment),
			"handler returned wrong body: got %v want %v",
			comment,
			updatedComment,
		)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("returns 400 if invalid comment ID", func(t *testing.T) {
		req, err := http.NewRequest(
			"PUT",
			fmt.Sprintf("%s/comments/1", server.URL),
			reqBodyOf(invalidParams),
		)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("returns 500 on unexpected error", func(t *testing.T) {
		store.SetError(errors.NewInternalServerError(fmt.Errorf("unexpected error")))

		req, err := http.NewRequest(
			"PUT",
			fmt.Sprintf("%s/comments/1", server.URL),
			reqBodyOf(validParams),
		)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestHandleUpdateCommentUnauthorized(t *testing.T) {
	store := &db.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /comments/{id}",
		MiddlewareChain(
			handleUpdateComments(store),
			LoggerMW,
			ErrorHandler,
			WithJWTAuth(store),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returns 401 if not authenticated", func(t *testing.T) {
		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/comments/1", server.URL), nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestHandleDeleteComment(t *testing.T) {
	store := &db.MockStore{}

	comment := db.Comment{
		ID: 1,
	}
	store.SetComment(comment)

	mux := http.NewServeMux()
	mux.HandleFunc("/comments/{id}",
		MiddlewareChain(
			handleDeleteComments(store),
			LoggerMW,
			ErrorHandler,
			fakeAuth(1),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returns 204 if valid", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/comments/1", server.URL), nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("returns 404 if comment not found", func(t *testing.T) {
		store.SetError(errors.NewNotFoundError("comment not found"))

		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/comments/1", server.URL), nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("returns 500 on unexpected error", func(t *testing.T) {
		store.SetError(errors.NewInternalServerError(fmt.Errorf("unexpected error")))

		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/comments/1", server.URL), nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestHandleDeleteCommentUnauthorized(t *testing.T) {
	store := &db.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("/comments/{id}",
		MiddlewareChain(
			handleDeleteComments(store),
			LoggerMW,
			ErrorHandler,
			WithJWTAuth(store),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returns 401 if unauthorized", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/comments/1", server.URL), nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestHandleLikeComment(t *testing.T) {
	store := &db.MockStore{}

	validParams := db.CreateCommentLikeParams{
		UserID:    1,
		CommentID: 1,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/comments/{id}/likes",
		MiddlewareChain(
			handleLikeComment(store),
			LoggerMW,
			ErrorHandler,
			fakeAuth(1),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returns 204 if valid", func(t *testing.T) {
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s/comments/1/likes", server.URL),
			reqBodyOf(validParams),
		)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.True(t, store.LikeCommentCalled())
	})

	t.Run("returns 500 on error", func(t *testing.T) {
		store.SetError(errors.NewInternalServerError(fmt.Errorf("error liking comment")))

		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s/comments/1/likes", server.URL),
			reqBodyOf(validParams),
		)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestHandleRemoveLikeFromComment(t *testing.T) {
	store := &db.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("/comments/{id}/likes",
		MiddlewareChain(
			handleRemoveLikeFromComment(store),
			LoggerMW,
			ErrorHandler,
			fakeAuth(1),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returns 204 if valid", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/comments/1/likes", server.URL), nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.True(t, store.UnlikeCommentCalled())
	})

	t.Run("returns 404 if comment not found", func(t *testing.T) {
		store.SetError(errors.NewNotFoundError("comment not found"))

		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/comments/1/likes", server.URL), nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
	t.Run("returns 500 on unexpected error", func(t *testing.T) {
		store.SetError(errors.NewInternalServerError(fmt.Errorf("unexpected error")))

		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/comments/1/likes", server.URL), nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestHandleRemoveLikeFromCommentUnauthorized(t *testing.T) {
	store := &db.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("/comments/{id}/likes",
		MiddlewareChain(
			handleRemoveLikeFromComment(store),
			LoggerMW,
			ErrorHandler,
			WithJWTAuth(store),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returns 401 if unauthorized", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/comments/1/likes", server.URL), nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.False(t, store.UnlikeCommentCalled())
	})
}
