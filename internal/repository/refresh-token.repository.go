package repository

import (
	"github.com/NekKkMirror/medods-tz.git/internal/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type RefreshTokenRepository struct {
	db *sqlx.DB
}

// NewRefreshTokenRepository creates a new instance of RefreshTokenRepository.
func NewRefreshTokenRepository(db *sqlx.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Save saves a new refresh token for a given user.
func (repo *RefreshTokenRepository) Save(userID string, token string, clientIP string) error {
	_, err := repo.db.Exec("INSERT INTO refresh_tokens (user_id, token_hash, client_ip) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET token_hash = EXCLUDED.token_hash, client_ip = EXCLUDED.client_ip", userID, token, clientIP)
	return err
}

// GetByToken retrieves the refresh token associated with a given token hash from the database.
func (repo *RefreshTokenRepository) GetByToken(token string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	err := repo.db.Get(&refreshToken, "SELECT * FROM refresh_tokens WHERE token_hash = $1", token)
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}
