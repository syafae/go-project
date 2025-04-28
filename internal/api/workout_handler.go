package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/syafae/femProject/internal/middleware"
	"github.com/syafae/femProject/internal/store"
	"github.com/syafae/femProject/internal/utils"
)

type WorkoutHandler struct {
	WorkoutStore store.WorkoutStore // the apis know only about the interface only to decouple the database from thr api
	Logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		WorkoutStore: workoutStore,
		Logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.Logger.Printf("ERRR:ReadIDParam %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}
	workout, err := wh.WorkoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.Logger.Printf("ERRR:GetWorkoutByID %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.Logger.Printf("ERRR:decodingCreateWorkout %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}
	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser.IsAnonymous() {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "unauthorized"})
		return
	}
	workout.UserID = currentUser.ID
	createdWorkout, err := wh.WorkoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.Logger.Printf("ERRR:CreateWorkout %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"createdWorkout": createdWorkout})

}
func (wh *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.Logger.Printf("ERRR:ReadIDParam %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	existingWorkout, err := wh.WorkoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.Logger.Printf("ERRR:GetWorkoutByID %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if existingWorkout == nil {
		http.NotFound(w, r)
		//utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}
	var updatedWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}
	err = json.NewDecoder(r.Body).Decode(&updatedWorkoutRequest)
	if err != nil {
		wh.Logger.Printf("ERRR:decodingUpdateRequest %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	if updatedWorkoutRequest.Title != nil {
		existingWorkout.Title = *updatedWorkoutRequest.Title
	}
	if updatedWorkoutRequest.Description != nil {
		existingWorkout.Description = *updatedWorkoutRequest.Description
	}
	if updatedWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updatedWorkoutRequest.DurationMinutes
	}
	if updatedWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updatedWorkoutRequest.CaloriesBurned
	}
	if updatedWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updatedWorkoutRequest.Entries
	}

	currentUser := middleware.GetUser(r)
	// Check if the user is authenticated
	// and if the workout belongs to the user

	if currentUser == nil || currentUser.IsAnonymous() {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "unauthorized"})
		return
	}
	workoutOwner, err := wh.WorkoutStore.GetWorkoutOwnerID(workoutID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
	}
	if workoutOwner != currentUser.ID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "forbidden"})
		return
	}
	err = wh.WorkoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		wh.Logger.Printf("ERRR:UpdateWorkout %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})

}

func (wh *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.Logger.Printf("ERRR:ReadIDParam %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}
	currentUser := middleware.GetUser(r)
	// Check if the user is authenticated
	// and if the workout belongs to the user

	if currentUser == nil || currentUser.IsAnonymous() {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "unauthorized"})
		return
	}
	workoutOwner, err := wh.WorkoutStore.GetWorkoutOwnerID(workoutID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
	}
	if workoutOwner != currentUser.ID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "forbidden"})
		return
	}
	err = wh.WorkoutStore.DeleteWorkout(workoutID)
	if err == sql.ErrNoRows {
		wh.Logger.Printf("ERRR:DeleteWorkout %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}
	if err != nil {
		wh.Logger.Printf("ERRR:DeleteWorkout %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": "workout deleted"})
}
