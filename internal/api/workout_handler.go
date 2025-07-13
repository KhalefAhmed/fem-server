package api

import (
	"encoding/json"
	"fmt"
	"github.com/KhalefAhmed/fem-server/internal/store"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
}

func NewWorkoutHandler(workoutStore store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {

	workoutId, err, done := handleIdParam(w, r)
	if done {
		return
	}

	workout, err := wh.workoutStore.GetWorkout(workoutId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to retrieve the workout :"+err.Error(), http.StatusNotFound)
		return
	}

	if workout == nil {
		http.Error(w, "workout not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workout)

	fmt.Fprintf(w, "this is the workout id %d\n", workoutId)

}

func handleIdParam(w http.ResponseWriter, r *http.Request) (int64, error, bool) {
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return 0, nil, true
	}
	workoutId, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return 0, nil, true
	}
	return workoutId, err, false
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to insert in base workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdWorkout)
}

func (wh *WorkoutHandler) HandleUpdateWorkoutById(w http.ResponseWriter, r *http.Request) {

	workoutId, _, done := handleIdParam(w, r)
	if done {
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkout(workoutId)
	if err != nil {
		http.Error(w, "Failed to fetch workout : "+err.Error(), http.StatusInternalServerError)
		return
	}

	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	err = wh.workoutStore.UpdateWorkout(existingWorkout)

	if err != nil {
		fmt.Println("Update workout error", err)
		http.Error(w, "Failed to update the workout", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingWorkout)
}

var updateWorkoutRequest struct {
	Title           *string              `json:"title"`
	Description     *string              `json:"description"`
	DurationMinutes *int                 `json:"duration_minutes"`
	CaloriesBurned  *int                 `json:"calories_burned"`
	Entries         []store.WorkoutEntry `json:"entries"`
}
