package app

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/NekKkMirror/medods-tz.git/internal/container/postgres"
	"github.com/NekKkMirror/medods-tz.git/internal/handler"
	"github.com/NekKkMirror/medods-tz.git/internal/repository"
	"github.com/NekKkMirror/medods-tz.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("PORT", "8080")
	os.Setenv("JWT_SECRET", "H$jeFJ*#(SJfjJF#*(Wue23")
}

func InitializeApp(db *sqlx.DB) *gin.Engine {
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewRefreshTokenRepository(db)
	authService := service.NewAuthService(tokenRepo, userRepo)

	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()
	r.GET("/issue", authHandler.IssueTokens)
	r.POST("/refresh", authHandler.RefreshTokens)

	return r
}

func TestIssueTokens(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	dbConn, terminate, err := postgres.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}
	defer terminate()

	r := InitializeApp(dbConn)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/issue?userID=3100b0c6-c6cc-4edf-a9a8-444990d0547d", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["accessToken"])
	assert.NotEmpty(t, response["refreshToken"])
}

func TestRefreshTokens(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	dbConn, terminate, err := postgres.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}
	defer terminate()

	r := InitializeApp(dbConn)

	w := httptest.NewRecorder()
	reqIssue, _ := http.NewRequest("GET", "/issue?userID=3100b0c6-c6cc-4edf-a9a8-444990d0547d", nil)
	r.ServeHTTP(w, reqIssue)

	assert.Equal(t, http.StatusOK, w.Code)

	var issueResponse map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &issueResponse)
	assert.NoError(t, err)

	refreshToken := issueResponse["refreshToken"]

	jsonStr := []byte(`refreshToken=` + refreshToken)
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var refreshResponse map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &refreshResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshResponse["accessToken"])
	assert.NotEmpty(t, refreshResponse["refreshToken"])
}
