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
	. "github.com/malyshEvhen/meow_mingle/internal/utils"
	"github.com/stretchr/testify/assert"
)

type Comments []db.CommentInfo

func TestHandleCreateComment(t *testing.T) {
	var (
		authUserID int64 = 1
		store            = &mock.MockStore{}
		router           = mux.NewRouter()
		path             = "/posts/{id}/comments"
		method           = "POST"
	)

	router.HandleFunc(path,
		middleware.MiddlewareChain(
			HandleCreateComment(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
			middleware.WithJWTAuth(store, testCfg),
		),
	).Methods(method)

	server := httptest.NewServer(router)
	defer server.Close()

	type input struct {
		userID  int64
		postID  int
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
				userID: authUserID,
				postID: 1,
				comment: db.CreateCommentParams{
					PostID:  1,
					Content: "Test Comment",
				},
			},
			want: want{
				status: http.StatusCreated,
				comment: db.Comment{
					ID:       1,
					AuthorID: authUserID,
					Content:  "Test Comment",
				},
			},
		},
		{
			name: "returns 400 if invalid post content",
			input: input{
				userID: authUserID,
				postID: 1,
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
	}

	for _, tc := range testCases {
		store.SetError(tc.input.error)
		store.SetComment(tc.want.comment)

		t.Run(tc.name, func(t *testing.T) {
			req, err := newAuthRequest(
				method,
				fmt.Sprintf("%s/posts/%d/comments", server.URL, tc.input.postID),
				reqBodyOf(tc.input.comment),
				tc.input.userID,
			)
			assert.NoError(t, err, "creating request")

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err, "perform the request")
			defer resp.Body.Close()

			t.Run("check that request body is the same as expected", func(t *testing.T) {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err, "read response body")

				createdComment, err := Unmarshal[db.Comment](body)
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

		t.Run("returns 401 if unauthenticated", func(t *testing.T) {
			req, err := http.NewRequest("POST", fmt.Sprintf("%s/posts/1/comments", server.URL), nil)
			assert.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err, "perform the request")
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	}
}

func TestHandleGetComments(t *testing.T) {
	var (
		authUserID int64 = 1
		router           = mux.NewRouter()
		store            = &mock.MockStore{}
		path             = "/posts/{id}/comments"
		method           = "GET"
	)

	router.HandleFunc(path,
		middleware.MiddlewareChain(
			HandleGetComments(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
		),
	).Methods(method)

	server := httptest.NewServer(router)
	defer server.Close()

	type input struct {
		userID int64
		postID int
		error  errors.Error
	}

	type want struct {
		status  int
		comment db.CommentInfo
	}

	testCases := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "returns 200 if comment exists",
			input: input{
				userID: authUserID,
				postID: 1,
			},
			want: want{
				status: http.StatusOK,
				comment: db.CommentInfo{
					ID:       1,
					AuthorID: 1,
					PostID:   1,
					Content:  "Test Comment",
				},
			},
		},
		{
			name: "returns 404 if post not found",
			input: input{
				userID: authUserID,
				postID: 1,
				error:  errors.NewNotFoundError("post not found"),
			},
			want: want{
				status: http.StatusNotFound,
			},
		},
	}

	for _, tc := range testCases {
		store.SetError(tc.input.error)
		store.AddComments(tc.want.comment)

		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(method, fmt.Sprintf("%s/posts/%d/comments", server.URL, tc.input.postID), nil)
			assert.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err, "perform the request")
			defer resp.Body.Close()

			t.Run("check if status code is correct", func(t *testing.T) {
				assert.Equalf(t, tc.want.status, resp.StatusCode,
					"handler returned wrong status code: got %v want %v",
					resp.StatusCode, tc.want.status)
			})

			if tc.want.status == 200 {
				t.Run("check that body is not empty and contains comment that we expect", func(t *testing.T) {
					body, err := io.ReadAll(resp.Body)
					assert.NoError(t, err)

					respComments, err := Unmarshal[Comments](body)
					assert.NoErrorf(t, err, "unmarshal response body")
					assert.NotEmpty(t, respComments, "empty comments")
					assert.Equalf(t, tc.want.comment, respComments[0],
						"handler returned unexpected body: got %v want %v",
						respComments[0], tc.want.comment)
				})
			}
		})
	}
}

func TestHandleUpdateCommentAuthorized(t *testing.T) {
	var (
		router = mux.NewRouter()
		store  = &mock.MockStore{}
		path   = "/comments/{id}"
		method = "PUT"
	)

	type input struct {
		commentID int
		comment   db.UpdateCommentParams
		error     errors.Error
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
			name: "should return 200 if params is valid",
			input: input{
				commentID: 1,
				comment: db.UpdateCommentParams{
					ID:      1,
					Content: "Updated Comment",
				},
			},
			want: want{
				status: http.StatusOK,
				comment: db.Comment{
					ID:       1,
					AuthorID: 1,
					Content:  "Updated Comment",
				},
			},
		},
		// TODO: write TC for each invalid parameter.
		{
			name: "should return 400 if params is invalid",
			input: input{
				commentID: 1,
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "should return 500 if DB error was occured",
			input: input{
				error:     errors.NewInternalServerError(fmt.Errorf("unexpected error")),
				commentID: 1,
				comment: db.UpdateCommentParams{
					ID:      1,
					Content: "Updated Comment",
				},
			},
			want: want{
				status: http.StatusInternalServerError,
			},
		},
	}

	for _, tc := range testCases {
		store.SetError(tc.input.error)
		store.SetComment(tc.want.comment)

		router.HandleFunc(path,
			middleware.MiddlewareChain(
				HandleUpdateComments(store),
				middleware.LoggerMW,
				middleware.ErrorHandler,
				fakeAuth(1),
			),
		).Methods(method)

		server := httptest.NewServer(router)
		defer server.Close()

		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(
				method,
				fmt.Sprintf("%s/comments/%d", server.URL, tc.input.commentID),
				reqBodyOf(tc.input.comment),
			)
			assert.NoError(t, err, "creating request")

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err, "perform the request")
			defer resp.Body.Close()

			t.Run("check if status code is correct", func(t *testing.T) {
				assert.Equal(t, tc.want.status, resp.StatusCode)
			})

			if tc.want.status == 200 {
				t.Run("check that body is not empty and contains comment that we expect", func(t *testing.T) {
					body, err := io.ReadAll(resp.Body)
					assert.NoError(t, err, "read response body")

					comment, err := Unmarshal[db.Comment](body)
					assert.NoErrorf(t, err, "unmarshal response body")
					assert.Truef(
						t,
						reflect.DeepEqual(comment, tc.want.comment),
						"handler returned wrong body: got %v want %v",
						comment,
						tc.want.comment,
					)
				})
			}
		})
	}
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
	var (
		authUserID int64 = 1
		router           = mux.NewRouter()
		store            = &mock.MockStore{}
		path             = "/comments/{id}"
		method           = "DELETE"
	)

	router.HandleFunc(path,
		middleware.MiddlewareChain(
			HandleDeleteComments(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
			middleware.WithJWTAuth(store, testCfg),
		),
	).Methods(method)

	server := httptest.NewServer(router)
	defer server.Close()

	type input struct {
		commentID int
		userID    int64
		error     errors.Error
	}

	type want struct {
		status      int
		storeCalled bool
		comment     db.Comment
	}

	testCases := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "returns 204 if delete parametrs are valid",
			input: input{
				commentID: 1,
				userID:    authUserID,
			},
			want: want{
				status: http.StatusNoContent,
				comment: db.Comment{
					ID:       1,
					Content:  "Updated Content",
					AuthorID: authUserID,
					PostID:   1,
				},
				storeCalled: true,
			},
		},
		// {
		// 	name: "returns 403 if authenticated user isn`t the autor of the comment",
		// 	input: input{
		// 		commentID: 1,
		// 		userID:    authUserID,
		// 	},
		// 	want: want{
		// 		status: http.StatusForbidden,
		// 		comment: db.Comment{
		// 			ID:       1,
		// 			Content:  "Updated Content",
		// 			AuthorID: 100,
		// 			PostID:   1,
		// 		},
		// 		storeCalled: true,
		// 	},
		// },
		// {
		// 	name: "returns 404 if comment not found",
		// 	input: input{
		// 		commentID: 1,
		// 		userID:    1,
		// 		error:     errors.NewNotFoundError("comment not found"),
		// 	},
		// 	want: want{
		// 		status: http.StatusForbidden,
		// 		comment: db.Comment{
		// 			ID:       1,
		// 			Content:  "Test Comment",
		// 			AuthorID: 2,
		// 			PostID:   1,
		// 		},
		// 		storeCalled: true,
		// 	},
		// },
		// {
		// 	name: "returns 500 on unexpected error",
		// 	input: input{
		// 		commentID: 1,
		// 		userID:    1,
		// 		error:     errors.NewInternalServerError(fmt.Errorf("unexpected error")),
		// 	},
		// 	want: want{
		// 		status:      http.StatusInternalServerError,
		// 		storeCalled: true,
		// 	},
	}
	for _, tc := range testCases {
		store.SetError(tc.input.error)
		store.SetComment(tc.want.comment)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req, err := newAuthRequest(
					method,
					fmt.Sprintf("%s/comments/%d", server.URL, tc.input.commentID),
					nil,
					int64(tc.input.userID),
				)
				assert.NoError(t, err, "create request")

				resp, err := http.DefaultClient.Do(req)
				assert.NoError(t, err, "perform request")
				defer resp.Body.Close()

				assert.Equal(t, tc.want.storeCalled, store.DeleteCommentCalled())
				assert.Equal(t, tc.want.status, resp.StatusCode)
			})
		}
	}
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
