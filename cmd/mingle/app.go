package mingle

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/handlers"
	"github.com/malyshEvhen/meow_mingle/internal/api"
	"github.com/malyshEvhen/meow_mingle/internal/app/comment"
	"github.com/malyshEvhen/meow_mingle/internal/app/post"
	"github.com/malyshEvhen/meow_mingle/internal/app/profile"
	"github.com/malyshEvhen/meow_mingle/internal/app/reaction"
	"github.com/malyshEvhen/meow_mingle/internal/app/subscription"
	"github.com/malyshEvhen/meow_mingle/internal/auth"
	"github.com/malyshEvhen/meow_mingle/internal/db"
)

type App struct {
	srv *http.Server

	authProvider *auth.Provider
	session      *gocql.Session
}

func New(ctx context.Context) (mingleApp *App, appError error) {
	cfg, err := initConfig()
	if err != nil {
		return nil, fmt.Errorf("an error occurred when config initializes: %s", err.Error())
	}

	authProvider := auth.NewProvider(nil, cfg.JWTSecret)

	cluster := gocql.NewCluster("localhost:9042")
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("an error occurred when creating session: %s", err.Error())
	}

	profileRepo := db.NewProfileRepository(session)
	commentRepo := db.NewCommentRepository(session)
	postRepo := db.NewPostRepository(session)
	subscriptionRepo := db.NewSubscriptionRepository(session)
	reactionRepo := db.NewReactionRepository(session)

	profileService := profile.NewService(profileRepo)
	commentService := comment.NewService(commentRepo)
	postService := post.NewService(postRepo)
	subscriptionService := subscription.NewService(subscriptionRepo)
	reactionService := reaction.NewService(reactionRepo)

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
		session:      session,
	}, nil
}

func (app *App) Start(ctx context.Context) error {
	if err := app.srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (app *App) Stop(ctx context.Context) error {
	defer app.session.Close()

	if err := app.srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
