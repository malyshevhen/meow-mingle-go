package api

import (
	"net/http"
	"os"
	"testing"
)

func testMW(userID int64, h Handler) http.HandlerFunc {
	return MiddlewareChain(
		h,
		LoggerMW,
		ErrorHandler,
		fakeAuth(userID),
	)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
