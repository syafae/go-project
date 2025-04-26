package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Envelope map[string]any

func WriteJSON(w http.ResponseWriter, staus int, data Envelope) error {
	json, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return err
	}
	json = append(json, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(staus)
	w.Write(json)
	return nil
}

func ReadIDParam(r *http.Request) (int64, error) {
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		return 0, errors.New("invalid id parameter")
	}
	workoutID, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		return 0, errors.New("invalid ID parameter")
	}

	return workoutID, nil
}
