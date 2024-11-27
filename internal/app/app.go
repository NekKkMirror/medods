package app

import (
	"github.com/NekKkMirror/medods-tz.git/config"
	"github.com/NekKkMirror/medods-tz.git/internal/handler"
	"github.com/NekKkMirror/medods-tz.git/internal/repository"
	"github.com/NekKkMirror/medods-tz.git/internal/service"
	"github.com/gin-gonic/gin"
)

// Initialize initializes the application and returns a Gin engine instance and an error.
func Initialize() (*gin.Engine, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}

	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	userRepo := repository.NewUserRepository(db)

	authService := service.NewAuthService(refreshTokenRepo, userRepo)
	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()
	r.GET("/issue", authHandler.IssueTokens)
	r.POST("/refresh", authHandler.RefreshTokens)

	return r, nil
}
