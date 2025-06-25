package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/malyshEvhen/meow_mingle/internal/api"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type App struct {
	s *http.Server

	userRepo     db.IUserRepository
	commentRepo  db.ICommentRepository
	postRepo     db.IPostRepository
	authProvider *api.AuthProvider

	driver neo4j.DriverWithContext
}

func New(ctx context.Context) (app *App, appError error) {
	cfg, err := initConfig()
	if err != nil {
		return nil, fmt.Errorf("an error occurred when config initializes: %s", err.Error())
	}

	driver, err := neo4j.NewDriverWithContext(cfg.DBConnURL, neo4j.BasicAuth(cfg.DBUser, cfg.DBPassword, ""))
	if err != nil {
		return nil, fmt.Errorf("an error occurred when neo4j driver creates: %s", err.Error())
	}

	app = &App{}
	app.driver = driver
	app.userRepo = db.NewUserRepository(driver)
	app.commentRepo = db.NewCommentRepository(driver)
	app.postRepo = db.NewPostRepository(driver)
	app.authProvider = api.NewAuthProvider(app.userRepo, cfg.JWTSecret)

	mux := api.RegisterRouts(
		ctx,
		app.authProvider,
		app.userRepo,
		app.commentRepo,
		app.postRepo,
		cfg.JWTSecret,
	)

	recoveryHandler := handlers.RecoveryHandler()
	corsHandler := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		handlers.AllowCredentials(),
		handlers.ExposedHeaders([]string{"Authorization", "Content-Type", "Content-Encoding", "Content-Length", "Location"}),
	)

	app.s = &http.Server{
		Addr:         cfg.ServerPort,
		Handler:      corsHandler(recoveryHandler(mux)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	return app, nil
}

func (app *App) Start(ctx context.Context) error {
	if err := app.s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (app *App) Stop(ctx context.Context) error {
	if err := app.driver.Close(ctx); err != nil {
		return err
	}
	return nil
}
