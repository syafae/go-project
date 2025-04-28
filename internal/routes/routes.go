package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/syafae/femProject/internal/app"
)

func SetUpRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenicate)
		// Add any middleware that should be applied to all routes here
		//workouts
		r.Get("/workout/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleGetWorkoutByID))
		r.Post("/workouts", app.Middleware.RequireUser(app.WorkoutHandler.HandleCreateWorkout))
		r.Put("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleUpdateWorkoutByID))
		r.Delete("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleDeleteWorkoutByID))

		//users
		r.Get("/users/{username}", app.Middleware.RequireUser(app.UserHandler.HandleGetUserByName))
		r.Put("/users/{username}", app.Middleware.RequireUser(app.UserHandler.HandleUpdateUser))

	})
	// Add routes that don't require authentication here

	r.Get("/health", app.HealthCheck)

	//user
	r.Post("/users", app.UserHandler.HandleRegiserUserRequest)

	//tokens
	r.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)

	return r
}
