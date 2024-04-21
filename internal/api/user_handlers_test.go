package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/errors"
	"github.com/malyshEvhen/meow_mingle/internal/middleware"
	"github.com/malyshEvhen/meow_mingle/internal/mock"
	"github.com/malyshEvhen/meow_mingle/internal/types"
	"github.com/malyshEvhen/meow_mingle/internal/utils"
	"github.com/stretchr/testify/assert"
)

type Feed []db.ListUserPostsRow

func TestHandleCreateUser(t *testing.T) {

	user := db.User{
		ID:        1,
		Email:     "john@doe.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password",
		CreatedAt: time.Time{},
	}

	validParams := db.CreateUserParams{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
	}

	invalidParams := db.CreateUserParams{
		Email:     "",
		FirstName: "",
		LastName:  "",
		Password:  "",
	}

	t.Run("should create a new user", func(t *testing.T) {
		store := &mock.MockStore{}
		store.SetUser(user)

		req, err := http.NewRequest("POST", "/users", reqBodyOf(validParams))
		assert.NoError(t, err, "create request")
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(
			middleware.MiddlewareChain(
				HandleCreateUser(store),
				middleware.LoggerMW,
				middleware.ErrorHandler,
			),
		)
		handler.ServeHTTP(rr, req)

		assert.Equalf(t, http.StatusCreated, rr.Code,
			"handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusCreated)
	})

	t.Run("should return 400 if user already exists", func(t *testing.T) {
		store := &mock.MockStore{}
		store.SetUser(user)
		store.SetError(errors.NewValidationError("User already exists"))

		req, err := http.NewRequest("POST", "/users", reqBodyOf(validParams))
		assert.NoError(t, err, "create request")
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(
			middleware.MiddlewareChain(
				HandleCreateUser(store),
				middleware.LoggerMW,
				middleware.ErrorHandler,
			),
		)
		handler.ServeHTTP(rr, req)

		assert.Equalf(t, http.StatusBadRequest, rr.Code,
			"handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusBadRequest)
	})

	t.Run("should return 400 if params are invalid", func(t *testing.T) {
		store := &mock.MockStore{}

		req, err := http.NewRequest("POST", "/users", reqBodyOf(invalidParams))
		assert.NoError(t, err, "create request")
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(
			middleware.MiddlewareChain(
				HandleCreateUser(store),
				middleware.LoggerMW,
				middleware.ErrorHandler,
			),
		)
		handler.ServeHTTP(rr, req)

		assert.Equalf(t, http.StatusBadRequest, rr.Code,
			"handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusBadRequest)

		resp, err := utils.Unmarshal[middleware.ErrorResponse](rr.Body.Bytes())

		assert.NoError(t, err, "unmarshal error response")
		assert.True(t, strings.Contains(resp.Error, "Email"),
			"Message hasn`t contains error about 'Email' field",
		)
	})

	t.Run("should return 400 if body is empty", func(t *testing.T) {
		store := &mock.MockStore{}

		req, err := http.NewRequest("POST", "/users", reqBodyOf(db.CreateUserParams{}))
		assert.NoError(t, err, "create request")
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(
			middleware.MiddlewareChain(
				HandleCreateUser(store),
				middleware.LoggerMW,
				middleware.ErrorHandler,
				fakeAuth(1),
			),
		)
		handler.ServeHTTP(rr, req)

		assert.Equalf(t, http.StatusBadRequest, rr.Code,
			"handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusBadRequest)
	})
}

func TestHandleGetUser(t *testing.T) {

	userRow := db.GetUserRow{
		ID:        1,
		Email:     "john@doe.com",
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Time{},
	}

	store := &mock.MockStore{}
	store.SetGetUserRow(userRow)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}",
		middleware.MiddlewareChain(
			HandleGetUser(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
			fakeAuth(userRow.ID),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("should return 200 and user if user exists", func(t *testing.T) {

		emptyRequest := bytes.NewBuffer([]byte{})
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/1", server.URL), emptyRequest)
		assert.NoError(t, err, "create request")
		defer req.Body.Close()

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err, "read response body")

		if status := resp.StatusCode; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		user, err := utils.Unmarshal[db.GetUserRow](body)

		assert.NoErrorf(t, err, "unmarshal response body")
		assert.Truef(
			t,
			reflect.DeepEqual(user, userRow),
			"handler returned wrong body: got %v want %v",
			user,
			userRow,
		)
	})

	t.Run("should return 403 if user ID does not match auth user ID", func(t *testing.T) {

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/2", server.URL), nil)
		assert.NoError(t, err, "create request")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")
		defer resp.Body.Close()

		assert.Equalf(t, http.StatusForbidden, resp.StatusCode,
			"Expected status code: %d, but was: %d",
			http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 404 if user not found", func(t *testing.T) {

		store.SetError(errors.NewNotFoundError("user not found"))

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/1", server.URL), nil)
		assert.NoError(t, err, "create request")

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		assert.Equalf(t, http.StatusNotFound, rr.Code,
			"Expected status code: %d, but was: %d",
			http.StatusNotFound, rr.Code)
	})
}

func TestHandleAuthenticatedSubscribe(t *testing.T) {

	store := &mock.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users/{id}/subscriptions",
		middleware.MiddlewareChain(
			HandleSubscribe(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
			fakeAuth(int64(1)),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("should subscribe user to cat", func(t *testing.T) {
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s/users/2/subscriptions", server.URL),
			nil,
		)
		assert.NoError(t, err, "create request")

		resp, err := http.DefaultClient.Do(req)

		assert.NoError(t, err, "perform the request")
		assert.Equalf(t, http.StatusNoContent, resp.StatusCode,
			"handler returned wrong status code: got %v want %v",
			resp.StatusCode, http.StatusNoContent)
		assert.True(t, store.CreateSubscriptionCalled(), "CreateSubscription was not called")
	})

	t.Run("should return 400 if subscription ID is invalid", func(t *testing.T) {
		store := &mock.MockStore{}

		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s/users/2#/subscriptions", server.URL),
			nil,
		)
		assert.NoError(t, err, "create request")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(testMW(1, HandleSubscribe(store)))

		handler.ServeHTTP(rr, req)

		assert.Equalf(t, http.StatusBadRequest, rr.Code,
			"handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusBadRequest)
	})

	t.Run("should return 500 if database error", func(t *testing.T) {
		store.SetError(errors.NewInternalServerError(fmt.Errorf("database error")))

		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s/users/1/subscriptions", server.URL),
			nil,
		)
		assert.NoError(t, err, "create request")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")

		assert.Equalf(t, http.StatusInternalServerError, resp.StatusCode,
			"handler returned wrong status code: got %v want %v",
			resp.StatusCode, http.StatusInternalServerError)
	})
}

func TestHandleUnauthenticatedSubscribe(t *testing.T) {

	store := &mock.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users/{id}/subscriptions",
		middleware.MiddlewareChain(
			HandleSubscribe(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
			middleware.WithJWTAuth(store),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("should return 401 if user is not authenticated", func(t *testing.T) {

		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s/users/1/subscriptions", server.URL),
			nil,
		)
		assert.NoError(t, err, "create request")

		resp, err := http.DefaultClient.Do(req)

		assert.NoError(t, err, "perform the request")
		assert.Equalf(t, http.StatusUnauthorized, resp.StatusCode,
			"handler returned wrong status code: got %v want %v",
			resp.StatusCode, http.StatusUnauthorized)
		assert.False(t, store.CreateSubscriptionCalled(), "CreateSubscription was called")
	})
}

func TestHandleUnsubscribe(t *testing.T) {

	store := &mock.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /users/{id}/subscriptions",
		middleware.MiddlewareChain(
			HandleUnsubscribe(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
			fakeAuth(int64(1)),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("should unsubscribe user", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%d/subscriptions",
			server.URL, 2), nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.True(t, store.DeleteSubscriptionCalled())
	})

	t.Run("should return 404 if subscription not found", func(t *testing.T) {
		store.SetError(errors.NewNotFoundError("subscription not found"))

		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%d/subscriptions",
			server.URL, 1), nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 if subscription ID is invalid", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/users/bob/subscriptions",
			server.URL), nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestHandleOwnersFeedAuthenticated(t *testing.T) {

	store := &mock.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/feed",
		middleware.MiddlewareChain(
			HandleOwnersFeed(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
			fakeAuth(int64(1)),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("should return 200 and feed if authenticated user", func(t *testing.T) {
		row := db.ListUserPostsRow{
			ID:        1,
			AuthorID:  1,
			Content:   "Test Post",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			Likes:     0,
		}
		store.AddListUserPostsRows(row)

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/feed", server.URL), nil)
		assert.NoError(t, err, "create request")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		bodyBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.True(t, len(bodyBytes) > 0)

		respFeed, err := utils.Unmarshal[Feed](bodyBytes)
		assert.NoError(t, err)
		assert.Equal(t, row, respFeed[0])
	})

	t.Run("should return 500 if error getting feed", func(t *testing.T) {
		store.SetError(errors.NewDatabaseError(fmt.Errorf("error retrieve feed!")))

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/feed", server.URL), nil)
		assert.NoError(t, err, "create request")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestHandleOwnersFeedUnauthenticated(t *testing.T) {

	store := &mock.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/feed",
		middleware.MiddlewareChain(
			HandleOwnersFeed(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("should return 401 if no authenticated user", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/feed", server.URL), nil)
		assert.NoError(t, err, "create request")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestHandleUsersFeed(t *testing.T) {

	store := &mock.MockStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}/feed",
		middleware.MiddlewareChain(
			HandleUsersFeed(store),
			middleware.LoggerMW,
			middleware.ErrorHandler,
			fakeAuth(1),
		),
	)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("returns 200 and feed if valid user ID", func(t *testing.T) {
		row := db.ListUserPostsRow{
			ID:        1,
			AuthorID:  1,
			Content:   "Test post 1",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			Likes:     0,
		}
		store.AddListUserPostsRows(row)

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/1/feed", server.URL), nil)
		assert.NoError(t, err, "create request")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		bodyBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.True(t, len(bodyBytes) > 0)

		respFeed, err := utils.Unmarshal[Feed](bodyBytes)
		assert.NoError(t, err)
		assert.Equal(t, row, respFeed[0])
	})

	t.Run("returns 400 if invalid user ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/invalid/feed", server.URL), nil)
		assert.NoError(t, err, "create request")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("returns 500 if error getting feed", func(t *testing.T) {
		store.SetError(errors.NewDatabaseError(fmt.Errorf("error retrieve feed!")))

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/1/feed", server.URL), nil)
		assert.NoError(t, err, "create request")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "perform the request")

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func fakeAuth(id int64) types.Middleware {
	return func(h types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			rCtx := context.WithValue(r.Context(), utils.UserIdKey, id)
			r = r.WithContext(rCtx)

			return h(w, r)
		}
	}
}

func reqBodyOf(content interface{}) io.Reader {
	jsonBytes, _ := json.Marshal(content)

	return bytes.NewBuffer(jsonBytes)
}
