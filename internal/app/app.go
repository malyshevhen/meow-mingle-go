package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/router"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type closerFunc func(context.Context) error

func Start(ctx context.Context) (cf closerFunc, appError error) {
	var (
		cfg  = config.InitConfig()
		fail = func(fmsg string, a ...any) (closerFunc, error) {
			return nil, fmt.Errorf(fmsg, a)
		}
		driver neo4j.DriverWithContext
		store  db.IStore
		mux    *mux.Router
	)

	driver, err := neo4j.NewDriverWithContext(cfg.DBConnURL, neo4j.BasicAuth(cfg.DBUser, cfg.DBPassword, ""))
	if err != nil {
		return fail("an error occured when neo4j driver creates: %s", err.Error())
	}

	store = db.NewVstore(driver)

	mux = router.RegisterRoutes(store, cfg)
	if err := http.ListenAndServe(cfg.ServerPort, mux); err != nil {
		log.Printf("%-15s ==> Server failed to start: %s\n", "Application", err.Error())
		appError = fmt.Errorf("an error occured while server starts: %s", err.Error())
		return
	}

	cf = func(ctx context.Context) error {
		return driver.Close(ctx)
	}
	return
}
