package store

import "database/sql"

type Workout struct {
	ID              int            `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	ID              int      `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Reps            *int     `json:"reps"`
	Sets            int      `json:"sets"`
	Weight          *float64 `json:"weight"`
	DurationSeconds *int     `json:"duration_seconds"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}

type WorkoutStore interface {
	GetWorkoutByID(id int64) (*Workout, error)
	CreateWorkout(workout *Workout) (*Workout, error)
}

type postgressWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *postgressWorkoutStore {
	return &postgressWorkoutStore{db: db}
}

//implement the functions of the interface for the postgress database
//the app interact with interface and it doesnot know anything about postgres
//later on we can swap postgres with any relational database system and app does not care as
// he knows only the interface

func (pg *postgressWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	query := `INSERT INTO workouts (title, description, duration_minutes, calories_burned)
	 VALUES($1,$2,$3,$4,$5)
	 returning id
	 `
	err = tx.QueryRow(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)
	if err != nil {
		return nil, err
	}
	// at this point we need to insert the entries
	for _, entry := range workout.Entries {

		query := `INSERT INTO workout_entries (workout_id , excercise_name, sets, reps, duration_seconds, 
				weight, notes, order_index)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
				returning id`
		err = tx.QueryRow(query, entry.ID, entry.ExerciseName, entry.Sets, entry.reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}

	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return workout, nil
}


func (pg *postgressWorkoutStore) GetWorkoutByID(id int64) (*Workout, error) {
	

}