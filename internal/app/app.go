package app

import (
	"database/sql"
	"fmt"
	"github.com/KhalefAhmed/fem-server/internal/api"
	"github.com/KhalefAhmed/fem-server/internal/store"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	WorkoutHandler := api.NewWorkoutHandler()

	app := &Application{
		Logger:         logger,
		WorkoutHandler: WorkoutHandler,
		DB:             pgDB,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "Status is available\n")
	if err != nil {
	}
	return
}
