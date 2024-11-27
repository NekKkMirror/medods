package model

// User represents a user in the system.
type User struct {
	ID    string `db:"id"`
	Email string `db:"email"`
}
