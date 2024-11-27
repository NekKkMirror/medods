package service

import (
	"testing"

	"github.com/NekKkMirror/medods-tz.git/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRefreshTokenRepository используется для имитации RefreshTokenRepository.
type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Save(userID string, token string, clientIP string) error {
	args := m.Called(userID, token, clientIP)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetByToken(token string) (*model.RefreshToken, error) {
	args := m.Called(token)
	return args.Get(0).(*model.RefreshToken), args.Error(1)
}

// MockUserRepository используется для имитации UserRepository.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(id string, email string) error {
	args := m.Called(id, email)
	return args.Error(0)
}

func (m *MockUserRepository) GetById(id string) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func TestIssueTokens(t *testing.T) {
	mockTokenRepo := new(MockRefreshTokenRepository)
	mockUserRepo := new(MockUserRepository)

	authService := NewAuthService(mockTokenRepo, mockUserRepo)

	userID := "testuser"
	clientIP := "127.0.0.1"

	mockTokenRepo.On("Save", userID, mock.AnythingOfType("string"), clientIP).Return(nil)

	accessToken, refreshToken, err := authService.IssueTokens(userID, clientIP)

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	mockTokenRepo.AssertExpectations(t)
}

func TestRefreshToken(t *testing.T) {
	mockTokenRepo := new(MockRefreshTokenRepository)
	mockUserRepo := new(MockUserRepository)

	authService := NewAuthService(mockTokenRepo, mockUserRepo)

	userID := "testuser"
	clientIP := "127.0.0.1"
	token := "somevalidtoken"

	mockTokenRepo.On("GetByToken", token).Return(&model.RefreshToken{
		UserID:    userID,
		TokenHash: "somevalidtokenhash",
		ClientIP:  clientIP,
	}, nil)
	mockTokenRepo.On("Save", userID, mock.AnythingOfType("string"), clientIP).Return(nil)

	accessToken, newRefreshToken, err := authService.RefreshToken(token, clientIP)

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, newRefreshToken)

	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
