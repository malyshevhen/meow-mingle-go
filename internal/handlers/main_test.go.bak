package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/utils"

	"github.com/malyshEvhen/meow_mingle/internal/middleware"
	"github.com/malyshEvhen/meow_mingle/internal/types"
)

var testCfg = config.InitConfig()

func newAuthRequest(method, url string, body io.Reader, userID int64) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	secret := []byte(testCfg.JWTSecret)
	token, err := utils.CreateJwt(secret, userID)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     utils.TOKEN_COOKIE_KEY,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Duration(TOKEN_EXPIRATION_TIME) * time.Hour),
		Secure:   true,
		HttpOnly: true,
	}
	req.AddCookie(cookie)

	return req, nil
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

func testMW(userID int64, h types.Handler) http.HandlerFunc {
	return middleware.MiddlewareChain(
		h,
		middleware.LoggerMW,
		middleware.ErrorHandler,
		fakeAuth(userID),
	)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
