package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartPostgresContainer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	DB, terminate, err := Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}
	defer terminate()
	defer DB.Close()

	assert.NoError(t, err, "Database connection should not return an error")
	assert.NotNil(t, DB, "Database connection should not be nil")

	err = DB.Ping()
	assert.NoError(t, err, "Database ping should not return an error")
}
