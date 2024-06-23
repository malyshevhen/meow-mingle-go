package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/handlers"
	"github.com/malyshEvhen/meow_mingle/internal/middleware"
	"github.com/malyshEvhen/meow_mingle/internal/router"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type closerFunc func(context.Context) error

func Start(ctx context.Context) (cf closerFunc, appError error) {
	var (
		cfg    = config.InitConfig()
		driver neo4j.DriverWithContext
		muxer  *mux.Router
	)

	fail := func(format string, a ...any) (closerFunc, error) {
		return nil, fmt.Errorf(format, a)
	}

	driver, err := initDB(cfg)
	if err != nil {
		return fail(err.Error())
	}

	muxer = initRouter(ctx, cfg, driver)

	if err := listenAndServe(cfg, muxer); err != nil {
		return fail(err.Error())
	}

	cf = func(ctx context.Context) error {
		return driver.Close(ctx)
	}
	return
}

func initDB(cfg config.Config) (neo4j.DriverWithContext, error) {
	driver, err := neo4j.NewDriverWithContext(cfg.DBConnURL, neo4j.BasicAuth(cfg.DBUser, cfg.DBPassword, ""))
	if err != nil {
		return nil, fmt.Errorf("an error occured when neo4j driver creates: %s", err.Error())
	}

	return driver, nil
}

func initRouter(ctx context.Context, cfg config.Config, driver neo4j.DriverWithContext) *mux.Router {
	userRepo := db.NewUserReposiory(driver)
	postRepo := db.NewPostRepository(driver)
	commentRepo := db.NewCommentRepository(driver)

	authMW := middleware.NewAuthProvider(userRepo, cfg)

	muxer := mux.NewRouter()
	apiMux := muxer.PathPrefix("/api/v1").Subrouter()

	userHandler := handlers.NewUserHandler(userRepo)
	postHandler := handlers.NewPostHandler(postRepo)
	commentHandler := handlers.NewCommentHandler(commentRepo)

	userRouret := router.NewUserRouter(authMW, userHandler, postHandler)
	postRouter := router.NewPostRouter(authMW, postHandler, commentHandler, userRepo)
	commentRouter := router.NewCommentRouter(authMW, userRepo, commentHandler)

	userRouret.RegisterRouts(ctx, apiMux, cfg)
	postRouter.RegisterRouts(ctx, apiMux, cfg)
	commentRouter.RegisterRouts(ctx, apiMux, cfg)

	return muxer
}

func listenAndServe(cfg config.Config, router *mux.Router) error {
	if err := http.ListenAndServe(cfg.ServerPort, router); err != nil {
		log.Printf("%-15s ==> Server failed to start: %s\n", "Application", err.Error())
		return fmt.Errorf("an error occured while server starts: %s", err.Error())
	}
	return nil
}
