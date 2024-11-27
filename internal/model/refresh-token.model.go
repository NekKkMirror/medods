package model

// RefreshToken represents a refresh token in the system.
type RefreshToken struct {
	UserID    string `db:"user_id"`
	TokenHash string `db:"token_hash"`
	ClientIP  string `db:"client_ip"`
}
