package mingle

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/malyshEvhen/meow_mingle/internal/api"
	"github.com/malyshEvhen/meow_mingle/internal/app/comment"
	"github.com/malyshEvhen/meow_mingle/internal/app/post"
	"github.com/malyshEvhen/meow_mingle/internal/app/profile"
	"github.com/malyshEvhen/meow_mingle/internal/app/reaction"
	"github.com/malyshEvhen/meow_mingle/internal/app/subscription"
	"github.com/malyshEvhen/meow_mingle/internal/graph"
	"github.com/malyshEvhen/meow_mingle/pkg/auth"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type App struct {
	srv *http.Server

	authProvider *auth.Provider
	driver       neo4j.DriverWithContext
}

func New(ctx context.Context) (mingleApp *App, appError error) {
	cfg, err := initConfig()
	if err != nil {
		return nil, fmt.Errorf("an error occurred when config initializes: %s", err.Error())
	}

	driver, err := neo4j.NewDriverWithContext(cfg.DBConnURL, neo4j.BasicAuth(cfg.DBUser, cfg.DBPassword, ""))
	if err != nil {
		return nil, fmt.Errorf("an error occurred when neo4j driver creates: %s", err.Error())
	}

	authProvider := auth.NewProvider(nil, cfg.JWTSecret)

	profileRepo := graph.NewProfileRepository(driver)
	commentRepo := graph.NewCommentRepository(driver)
	postRepo := graph.NewPostRepository(driver)

	profileService := profile.NewService(profileRepo)
	commentService := comment.NewService(commentRepo)
	postService := post.NewService(postRepo)
	subscriptionService := subscription.NewService(nil)
	reactionService := reaction.NewService(nil)

	mux := api.RegisterRouts(
		authProvider,
		profileService,
		commentService,
		postService,
		subscriptionService,
		reactionService,
	)

	recoveryHandler := handlers.RecoveryHandler()
	corsHandler := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		handlers.AllowCredentials(),
		handlers.ExposedHeaders([]string{"Authorization", "Content-Type", "Content-Encoding", "Content-Length", "Location"}),
	)

	srv := &http.Server{
		Addr:         cfg.ServerPort,
		Handler:      corsHandler(recoveryHandler(mux)),
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	return &App{
		srv:          srv,
		authProvider: authProvider,
		driver:       driver,
	}, nil
}

func (app *App) Start(ctx context.Context) error {
	if err := app.srv.ListenAndServe(); err != nil {
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
