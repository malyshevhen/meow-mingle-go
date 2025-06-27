package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
	"github.com/malyshEvhen/meow_mingle/pkg/logger"
)

func handleSubscribe(subscriptionRepo app.SubscriptionService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := logger.GetLogger().WithComponent("subscription_handler")
		ctx := r.Context()

		id := mux.Vars(r)["id"]

		if err := subscriptionRepo.Subscribe(ctx, id); err != nil {
			logger.WithError(err).Error("Error subscribing")
			return err
		}

		logger.Info("Successfully subscribed")

		return writeJSON(w, http.StatusNoContent, nil)
	}
}

func handleUnsubscribe(subscriptionRepo app.SubscriptionService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := logger.GetLogger().WithComponent("subscription_handler")
		ctx := r.Context()

		id := mux.Vars(r)["id"]

		if err := subscriptionRepo.Unsubscribe(ctx, id); err != nil {
			logger.WithError(err).Error("Error unsubscribing")
			return err
		}

		logger.Info("Successfully unsubscribed")

		return writeJSON(w, http.StatusNoContent, nil)
	}
}
