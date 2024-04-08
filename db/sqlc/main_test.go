package db

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}
