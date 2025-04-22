package store

import (
	"database/sql"
	"fmt"
)

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
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}

type WorkoutStore interface {
	GetWorkoutByID(id int64) (*Workout, error)
	CreateWorkout(workout *Workout) (*Workout, error)
	UpdateWorkout(*Workout) error
	DeleteWorkout(int64) error
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
	 VALUES($1,$2,$3,$4)
	 returning id
	 `
	err = tx.QueryRow(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)
	if err != nil {
		return nil, err
	}
	// at this point we need to insert the entries
	for _, entry := range workout.Entries {

		query := `INSERT INTO workout_entries (workout_id , exercise_name, sets, 
				reps, duration_seconds, 
				weight, notes, order_index)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
				returning id`
		err = tx.QueryRow(query, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
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
	workout := &Workout{}
	query := `SELECT id, title, description, duration_minutes, calories_burned
				FROM workouts
				WHERE id = $1
			`
	err := pg.db.QueryRow(query, id).Scan(&workout.ID, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		fmt.Println("Mohammed")
		return nil, err
	}

	// let's get the entries
	entryQuery := `SELECT id, exercise_name,sets, reps, duration_seconds,weight, notes, order_index
		FROM workout_entries
		WHERE workout_id = $1
		ORDER BY order_index
	`

	rows, err := pg.db.Query(entryQuery, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry WorkoutEntry
		err = rows.Scan(
			&entry.ID,
			&entry.ExerciseName,
			&entry.Sets,
			&entry.Reps,
			&entry.DurationSeconds,
			&entry.Weight,
			&entry.Notes,
			&entry.OrderIndex,
		)
		if err != nil {
			return nil, err
		}

		workout.Entries = append(workout.Entries, entry)
	}

	return workout, nil
}

func (pg *postgressWorkoutStore) UpdateWorkout(workout *Workout) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE workouts
	SET title = $1, description = $2, duration_minutes = $3, calories_burned = $4
	WHERE id = $5
	`
	result, err := tx.Exec(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned, workout.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	//update the entries
	_, err = tx.Exec(`DELETE FROM workout_entries WHERE workout_id = $1`, workout.ID)
	if err != nil {
		return err
	}

	for _, entry := range workout.Entries {
		query := `
    INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

		_, err := tx.Exec(query,
			workout.ID,
			entry.ExerciseName,
			entry.Sets,
			entry.Reps,
			entry.DurationSeconds,
			entry.Weight,
			entry.Notes,
			entry.OrderIndex,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (pg *postgressWorkoutStore) DeleteWorkout(id int64) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `DELETE from workouts WHERE id = $1`
	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()

}
