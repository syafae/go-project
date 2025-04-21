-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS workout_entries (
    id SERIAL PRIMARY KEY,
    workout_id INT NOT NULL,
    excercise_name VARCHAR(255) NOT NULL,
    duration_second INT,
    weights DECIMAL(5, 2),
    sets INT NOT NULL,
    reps INT,
    notes TEXT,
    order_index INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE,
    CONSTRAINT valid_workout_entry CHECK(
        (reps IS NOT NULL OR duration_second IS NOT NULL) AND
        (reps IS NULL OR duration_second IS NULL)
    )
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS workout_entries;
-- +goose StatementEnd