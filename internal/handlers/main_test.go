package handlers

import (
	"github.com/malyshEvhen/meow_mingle/internal/config"
	"net/http"
	"os"
	"testing"

	"github.com/malyshEvhen/meow_mingle/internal/middleware"
	"github.com/malyshEvhen/meow_mingle/internal/types"
)

var testCfg = config.InitConfig()

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
