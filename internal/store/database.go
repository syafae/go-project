package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost user=postgres password= postgres dbname=postgres port= 5432 sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}
	fmt.Println("Connected to the database")
	return db, nil
}

func MigrateFS(db *sql.DB, MigrationFS fs.FS, dir string) error {
	goose.SetBaseFS(MigrationFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	goose.SetDialect("postgres")
	err := goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("db: migrate %w", err)
	}
	fmt.Println("Migrated the database")
	return nil
}
