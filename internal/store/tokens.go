package store

import (
	"database/sql"
	"time"

	"github.com/syafae/femProject/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

type TokenStore interface {
	InsertToken(token *tokens.Token) error
	CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int, scope string) error
	DeleteToken(token *tokens.Token) error
	GetTokenByHash(hash []byte) (*tokens.Token, error)
}

func (p *PostgresTokenStore) InsertToken(token *tokens.Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope) VALUES ($1, $2, $3, $4)`
	_, err := p.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	return err
}

func (p *PostgresTokenStore) CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = p.InsertToken(token)
	return token, err
}

func (p *PostgresTokenStore) DeleteAllTokensForUser(userID int, scope string) error {
	query := `DELETE FROM tokens WHERE user_id = $1 AND scope = $2`
	_, err := p.db.Exec(query, userID, scope)
	return err
}

func (p *PostgresTokenStore) DeleteToken(token *tokens.Token) error {
	query := `DELETE FROM tokens WHERE hash = $1`
	_, err := p.db.Exec(query, token.Hash)
	return err
}

func (p *PostgresTokenStore) GetTokenByHash(hash []byte) (*tokens.Token, error) {
	query := `SELECT user_id, expiry, scope FROM tokens WHERE hash = $1`
	token := &tokens.Token{}
	err := p.db.QueryRow(query, hash).Scan(&token.UserID, &token.Expiry, &token.Scope)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	token.Hash = hash
	return token, nil
}
