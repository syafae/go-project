package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/syafae/femProject/internal/app"
)

func SetUpRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", app.HealthCheck)
	//workouts
	r.Get("/workout/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkoutByID)
	r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkoutByID)

	//user
	r.Post("/users", app.UserHandler.HandleRegiserUserRequest)
	r.Get("/users/{username}", app.UserHandler.HandleGetUserByName)
	r.Put("/users/{username}", app.UserHandler.HandleUpdateUser)

	return r
}
