package store

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

// User represents a user in the system.
type User struct {
	ID           int       `json:"id"`
	UserName     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (p *password) Set(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	p.plainText = &password
	p.hash = hash
	return hash, nil
}

func (p *password) Matches(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil

}

// UserStore provides an interface for user-related DB operations.
type UserStore interface {
	CreateUser(user *User) error
	GetUserByName(username string) (*User, error)
	UpdateUser(user *User) error
}

type postgresUserStore struct {
	db *sql.DB
}

// NewPostgresUserStore returns a new instance of the Postgres based user store.
func NewPostgresUserStore(db *sql.DB) *postgresUserStore {
	return &postgresUserStore{db: db}
}

func (pg *postgresUserStore) CreateUser(user *User) error {
	query := `INSERT INTO users (username, email, password_hash, bio)
	          VALUES($1, $2, $3, $4)
			  RETURNING id, created_at, updated_at`
	err := pg.db.QueryRow(query, user.UserName, user.Email, user.PasswordHash.hash, user.Bio).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (pg *postgresUserStore) GetUserByName(username string) (*User, error) {
	query := `SELECT id, username, email, password_hash, bio, created_at, updated_at
	          FROM users
			  WHERE username = $1`
	user := &User{
		PasswordHash: password{},
	}
	var hash []byte
	err := pg.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user.PasswordHash.hash = hash
	return user, nil
}

func (pg *postgresUserStore) UpdateUser(user *User) error {
	query := `UPDATE users 
	          SET username = $1, email = $2, bio = $3, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = $4
			  RETURNING updated_at`
	result, err := pg.db.Exec(query, user.UserName, user.Email, user.Bio, user.ID)
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
	return nil
}
