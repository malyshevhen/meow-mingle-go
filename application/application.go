package application

import (
	"github.com/malyshEvhen/meow_mingle/api"
	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
)

type Application struct {
	db     *db.ConnectionPool
	store  *db.Store
	server *api.Server
}

func New(database *db.ConnectionPool) (*Application, error) {
	store := db.NewStore(database.DB)

	server := api.NewServer(":8080", store)

	return &Application{
		db:     database,
		store:  store,
		server: server,
	}, nil
}

func (app *Application) Start() error {
	defer app.db.Close()

	return app.server.Serve()
}
