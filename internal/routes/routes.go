package routes
import (
	"github.com/syafae/femProject/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetUpRoutes(app *app.Application) *chi.Mux  {
	r := chi.NewRouter()
	r.Get("/health", app.HealthCheck)
	r.Get("/workout/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	return r
}


