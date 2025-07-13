package app

import (
	"database/sql"
	"fmt"
	"github.com/KhalefAhmed/fem-server/internal/api"
	"github.com/KhalefAhmed/fem-server/internal/store"
	"github.com/KhalefAhmed/fem-server/migrations"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	WorkoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	UserHandler := api.NewUserHandler(userStore, logger)
	TokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)

	app := &Application{
		Logger:         logger,
		WorkoutHandler: WorkoutHandler,
		UserHandler:    UserHandler,
		TokenHandler:   TokenHandler,
		DB:             pgDB,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
