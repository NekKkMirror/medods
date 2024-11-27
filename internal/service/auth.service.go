package service

import (
	"errors"
	"log"

	"github.com/NekKkMirror/medods-tz.git/internal/repository"
	"github.com/NekKkMirror/medods-tz.git/internal/util"
)

type AuthService struct {
	refreshTokenRepo *repository.RefreshTokenRepository
	userRepo         *repository.UserRepository
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(refreshTokenRepo *repository.RefreshTokenRepository, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{refreshTokenRepo: refreshTokenRepo, userRepo: userRepo}
}

// IssueTokens generates and saves access and refresh tokens for a given user ID and client IP.
func (s *AuthService) IssueTokens(userID string, clientIP string) (string, string, error) {
	accessToken, err := util.GenerateAccessToken(clientIP)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := util.GenerateRandomString(32)
	if err != nil {
		return "", "", err
	}
	refreshTokenHash, err := util.HashRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	err = s.refreshTokenRepo.Save(userID, refreshTokenHash, clientIP)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshTokenHash, nil
}

// RefreshToken generates a new access and refresh tokens for a given refresh token and client IP.
func (s *AuthService) RefreshToken(refreshToken string, clientIP string) (string, string, error) {
	storedRefreshToken, err := s.refreshTokenRepo.GetByToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	userID := storedRefreshToken.UserID

	if storedRefreshToken.ClientIP != clientIP {
		user, err := s.userRepo.GetById(userID)
		if err != nil {
			return "", "", err
		}
		util.SendSecurityAlertEmail(user.Email, clientIP)
		return "", "", errors.New("IP address mismatch")
	}

	accessToken, err := util.GenerateAccessToken(clientIP)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := util.GenerateRandomString(32)
	if err != nil {
		return "", "", err
	}
	newRefreshTokenHash, err := util.HashRefreshToken(newRefreshToken)
	if err != nil {
		return "", "", err
	}

	err = s.refreshTokenRepo.Save(userID, newRefreshTokenHash, clientIP)
	if err != nil {
		log.Printf("Failed to save new refresh token: %v", err)
	}

	return accessToken, newRefreshToken, nil
}
