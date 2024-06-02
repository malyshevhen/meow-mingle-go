package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/errors"
	"github.com/malyshEvhen/meow_mingle/internal/middleware"
	"github.com/malyshEvhen/meow_mingle/internal/mock"
	"github.com/malyshEvhen/meow_mingle/internal/utils"
	"github.com/stretchr/testify/assert"
)

type Comments []db.CommentInfo

func TestHandleCreateComment(t *testing.T) {
	type input struct {
		comment db.CreateCommentParams
		error   errors.Error
	}

	type want struct {
		status  int
		comment db.Comment
	}

	testCases := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "happy path",
			input: input{
				comment: db.CreateCommentParams{
					PostID:  1,
					Content: "Test Comment",
				},
			},
			want: want{
				status: http.StatusCreated,
				comment: db.Comment{
					ID:       1,
					AuthorID: 1,
					Content:  "Test Comment",
				},
			},
		},
		{
			name: "returns 400 if invalid post ID",
			input: input{
				comment: db.CreateCommentParams{
					Content:  "",
					AuthorID: 0,
					PostID:   0,
				},
			},
			want: want{
				status:  http.StatusBadRequest,
				comment: db.Comment{},
			},
		},
		{
			name: "returns 500 if store returns an error",
			input: input{
				comment: db.CreateCommentParams{
					PostID:  1,
					Content: "Test Comment",
				},
				error: errors.NewDatabaseError(fmt.Errorf("%s", "db error")),
			},
			want: want{
				status:  http.StatusInternalServerError,
				comment: db.Comment{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := &mock.MockStore{}
			store.SetError(tc.input.error)
			store.SetComment(tc.want.comment)

			mux := mux.NewRouter()
			mux.HandleFunc("/posts/{id}/comments",
				middleware.MiddlewareChain(
					HandleCreateComment(store),
					middleware.LoggerMW,
					middleware.ErrorHandler,
					fakeAuth(1),
				),
			).Methods("POST")
			server := httptest.NewServer(mux)
			defer server.Close()

			req, err := http.NewRequest(
				"POST",
				fmt.Sprintf("%s/posts/1/comments", server.URL),
				reqBodyOf(tc.input.comment),
			)
			assert.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err, "perform the request")
			defer resp.Body.Close()

			t.Run("check that request body is the same as expected", func(t *testing.T) {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err, "read response body")
				createdComment, err := utils.Unmarshal[db.Comment](body)
				assert.NoErrorf(t, err, "unmarshal response body")
				assert.Truef(
					t,
					reflect.DeepEqual(createdComment, tc.want.comment),
					"handler returned wrong body: got %v want %v",
					createdComment,
					tc.want.comment,
				)
			})

			t.Run("check that status is correct", func(t *testing.T) {
				assert.Equal(t, tc.want.status, resp.StatusCode)
			})
		})
	}
}

func TestHandleCreateCommentUnauthenticated(t *testing.T) {
	store := &mock.MockStore{}

	mux := mux.NewRouter()
	mux.HandleFunc("/posts/{id}/comments",
		middleware.MiddlewareChain(
			HandleCreateComment(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
			middleware.WithJWTAuth(store, testCfg),
		),
	).Methods("POST")

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
	store := &mock.MockStore{}

	mux := mux.NewRouter()
	mux.HandleFunc("/posts/{id}/comments",
		middleware.MiddlewareChain(
			HandleGetComments(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
			fakeAuth(1),
		),
	).Methods("GET")

	commentRow := db.CommentInfo{
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

		respComments, err := utils.Unmarshal[Comments](body)
		assert.NoError(t, err)
		assert.NotEmpty(t, respComments)
		assert.Equalf(t, commentRow, respComments[0],
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
	store := &mock.MockStore{}

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

	mux := mux.NewRouter()
	mux.HandleFunc("/comments/{id}",
		middleware.MiddlewareChain(
			HandleUpdateComments(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
			fakeAuth(1),
		),
	).Methods("PUT")

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

		comment, err := utils.Unmarshal[db.Comment](body)
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
	store := &mock.MockStore{}

	mux := mux.NewRouter()
	mux.HandleFunc("/comments/{id}",
		middleware.MiddlewareChain(
			HandleUpdateComments(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
			middleware.WithJWTAuth(store, testCfg),
		),
	).Methods("PUT")

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
	store := &mock.MockStore{}

	comment := db.Comment{
		ID: 1,
	}
	store.SetComment(comment)

	mux := mux.NewRouter()
	mux.HandleFunc("/comments/{id}",
		middleware.MiddlewareChain(
			HandleDeleteComments(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
			fakeAuth(1),
		),
	).Methods("DELETE")

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
	store := &mock.MockStore{}

	mux := mux.NewRouter()
	mux.HandleFunc("/comments/{id}",
		middleware.MiddlewareChain(
			HandleDeleteComments(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
			middleware.WithJWTAuth(store, testCfg),
		),
	).Methods("DELETE")

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
	store := &mock.MockStore{}

	validParams := db.CreateCommentLikeParams{
		UserID:    1,
		CommentID: 1,
	}

	mux := mux.NewRouter()
	mux.HandleFunc("/comments/{id}/likes",
		middleware.MiddlewareChain(
			HandleLikeComment(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
			fakeAuth(1),
		),
	).Methods("POST")

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
	store := &mock.MockStore{}

	mux := mux.NewRouter()
	mux.HandleFunc("/comments/{id}/likes",
		middleware.MiddlewareChain(
			HandleRemoveLikeFromComment(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
			fakeAuth(1),
		),
	).Methods("DELETE")

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
	store := &mock.MockStore{}

	mux := mux.NewRouter()
	mux.HandleFunc("/comments/{id}/likes",
		middleware.MiddlewareChain(
			HandleRemoveLikeFromComment(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
			middleware.WithJWTAuth(store, testCfg),
		),
	).Methods("DELETE")

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
