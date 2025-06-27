package api

import (
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/internal/auth"
	"github.com/malyshEvhen/meow_mingle/pkg/logger"
)

func NewServer(
	cfg Config,
	authProvider *auth.Provider,
	profileService app.ProfileService,
	commentService app.CommentService,
	postService app.PostService,
	subscriptionService app.SubscriptionService,
	reactionService app.ReactionService,
) *http.Server {
	appLogger := logger.GetLogger()

	appLogger.WithComponent("service").Info("Business services initialized")

	mux := RegisterRouts(
		authProvider,
		profileService,
		commentService,
		postService,
		subscriptionService,
		reactionService,
	)

	appLogger.WithComponent("api").Info("API routes registered")

	recoveryHandler := handlers.RecoveryHandler()
	corsHandler := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		handlers.AllowCredentials(),
		handlers.ExposedHeaders([]string{"Authorization", "Content-Type", "Content-Encoding", "Content-Length", "Location"}),
	)

	appLogger.WithComponent("middleware").Info("HTTP middleware configured",
		"recovery", true,
		"cors", true,
	)

	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      corsHandler(recoveryHandler(mux)),
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	appLogger.WithComponent("server").Info("HTTP server configured",
		"addr", cfg.Port,
		"read_timeout", "2s",
		"write_timeout", "2s",
	)

	return srv
}
