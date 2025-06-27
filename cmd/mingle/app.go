package mingle

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/malyshEvhen/meow_mingle/internal/api"
	"github.com/malyshEvhen/meow_mingle/internal/app/comment"
	"github.com/malyshEvhen/meow_mingle/internal/app/post"
	"github.com/malyshEvhen/meow_mingle/internal/app/profile"
	"github.com/malyshEvhen/meow_mingle/internal/app/reaction"
	"github.com/malyshEvhen/meow_mingle/internal/app/subscription"
	"github.com/malyshEvhen/meow_mingle/internal/auth"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/pkg/logger"
)

type App struct {
	srv    *http.Server
	logger *logger.Logger

	authProvider *auth.Provider
	session      *gocql.Session
}

func New(ctx context.Context, cfg Config) (mingleApp *App, appError error) {
	appLogger := logger.GetLogger()

	authProvider := auth.NewProvider(nil)
	appLogger.WithComponent("auth").Info("Authentication provider initialized")

	cluster := gocql.NewCluster(cfg.Database.URL)
	appLogger.WithComponent("database").Info("Connecting to Cassandra", "host", cfg.Database.URL)

	session, err := cluster.CreateSession()
	if err != nil {
		appLogger.WithComponent("database").Error("Failed to create database session", "error", err.Error())
		return nil, fmt.Errorf("an error occurred when creating session: %s", err.Error())
	}

	appLogger.WithComponent("database").Info("Database session created successfully")

	profileRepo := db.NewProfileRepository(session)
	commentRepo := db.NewCommentRepository(session)
	postRepo := db.NewPostRepository(session)
	subscriptionRepo := db.NewSubscriptionRepository(session)
	reactionRepo := db.NewReactionRepository(session)

	appLogger.WithComponent("repository").Info("Database repositories initialized")

	profileService := profile.NewService(profileRepo)
	commentService := comment.NewService(commentRepo)
	postService := post.NewService(postRepo)
	subscriptionService := subscription.NewService(subscriptionRepo)
	reactionService := reaction.NewService(reactionRepo)

	srv := api.NewServer(
		cfg.Server,
		authProvider,
		profileService,
		commentService,
		postService,
		subscriptionService,
		reactionService,
	)

	return &App{
		srv:          srv,
		logger:       appLogger,
		authProvider: authProvider,
		session:      session,
	}, nil
}

func (app *App) Start(ctx context.Context) error {
	app.logger.WithComponent("server").Info("Starting HTTP server", "addr", app.srv.Addr)

	if err := app.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		app.logger.WithComponent("server").Error("HTTP server failed to start", "error", err.Error())
		return err
	}

	return nil
}

func (app *App) Stop(ctx context.Context) error {
	app.logger.WithComponent("server").Info("Shutting down HTTP server")

	if err := app.srv.Shutdown(ctx); err != nil {
		app.logger.WithComponent("server").Error("Failed to shutdown HTTP server", "error", err.Error())
		defer app.session.Close()
		return err
	}

	app.logger.WithComponent("database").Info("Closing database session")
	app.session.Close()
	app.logger.WithComponent("database").Info("Database session closed")

	return nil
}
