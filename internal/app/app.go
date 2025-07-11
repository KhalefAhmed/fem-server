package app

import (
	"fmt"
	"github.com/KhalefAhmed/fem-server/internal/api"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler api.WorkoutHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	WorkoutHandler := api.WorkoutHandler{}

	app := &Application{
		Logger:         logger,
		WorkoutHandler: WorkoutHandler,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "Status is available\n")
	if err != nil {
	}
	return
}
