package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
)

func handleSubscribe(subscriptionRepo app.SubscriptionService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id := mux.Vars(r)["id"]

		if err := subscriptionRepo.CreateSubscription(ctx, id); err != nil {
			return err
		}

		return writeJson(w, http.StatusNoContent, nil)
	}
}

func handleUnsubscribe(subscriptionRepo app.SubscriptionService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id := mux.Vars(r)["id"]

		if err := subscriptionRepo.DeleteSubscription(ctx, id); err != nil {
			return err
		}

		return writeJson(w, http.StatusNoContent, nil)
	}
}
