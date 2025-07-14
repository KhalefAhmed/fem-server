package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/KhalefAhmed/fem-server/internal/middleware"
	"github.com/KhalefAhmed/fem-server/internal/store"
	"github.com/KhalefAhmed/fem-server/internal/utils"
	"log"
	"net/http"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {

	workoutId, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: ReadIDParam: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	workout, err := wh.workoutStore.GetWorkout(workoutId)
	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutById: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: decoding create workout: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}
	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser.IsAnonymous() {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged"})
	}

	workout.UserID = currentUser.ID
	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: create workout : %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create workout"})
		return
	}
	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})
}

func (wh *WorkoutHandler) HandleUpdateWorkoutById(w http.ResponseWriter, r *http.Request) {

	workoutId, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: read Id param : %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}
	existingWorkout, err := wh.workoutStore.GetWorkout(workoutId)
	if err != nil {
		wh.logger.Printf("ERROR: get workout by id : %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal Server Error"})
		return
	}

	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}

	var updateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)

	if err != nil {
		wh.logger.Printf("ERROR:  error decoding update workout request: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser.IsAnonymous() {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged"})
		return
	}

	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(workoutId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJson(w, http.StatusNotFound, utils.Envelope{"error": "workout doesn't exist"})
			return
		}
	}

	err = wh.workoutStore.UpdateWorkout(existingWorkout)

	if err != nil {
		wh.logger.Printf("ERROR: updating workout : %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if workoutOwner != currentUser.ID {
		utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "you are not authorized to update this workout"})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}

func (wh *WorkoutHandler) HandleDeleteWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: read Id param : %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	err = wh.workoutStore.DeleteWorkout(workoutId)

	if errors.Is(err, sql.ErrNoRows) {
		wh.logger.Printf("ERROR: deleting workout : %v", err)
		utils.WriteJson(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}

	if err != nil {
		wh.logger.Printf("ERROR: deleting workout : %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJson(w, http.StatusNoContent, utils.Envelope{"workout": workoutId})
}
