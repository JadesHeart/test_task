package service

import (
	"database/sql"
	"log/slog"
	"test_task/internal/pkg/lib/sl"
	"test_task/internal/pkg/models"
	"test_task/internal/pkg/repository"
)

type AuthService struct {
	repoBD repository.AuthorizationBD
	logger *slog.Logger
}

func NewAuthServices(repo repository.AuthorizationBD, logger *slog.Logger) *AuthService {
	return &AuthService{
		repoBD: repo,
		logger: logger,
	}
}

func (s *AuthService) FindUser(username string) (*models.User, error) {
	user, err := s.repoBD.FindUser(username)
	if err != nil {
		s.logger.Error("Failed check user existence", sl.Err(err))
		return nil, err
	}
	if user.UserID == 0 {
		s.logger.Info("User does not exist")
		return nil, sql.ErrNoRows
	}
	return user, nil
}

func (s *AuthService) CheckFailedLoginAttempts(username string) (bool, error) {
	userBlock, err := s.repoBD.CheckFailedLoginAttempts(username)
	if err != nil {
		s.logger.Error("Failed check login attempts", sl.Err(err))

		return false, err
	}
	return userBlock, nil
}

func (s *AuthService) CheckPass(username string, password string) (bool, error) {
	passwordCorrect, err := s.repoBD.CheckPass(username, password)
	if err != nil {
		s.logger.Error("Failed compared password", sl.Err(err))

		return false, err
	}
	if !passwordCorrect {
		s.logger.Info("Incorrect password")

		err := s.repoBD.AddingFailedLoginAttempt(username)
		if err != nil {
			return false, err
		}
	}
	return passwordCorrect, nil
}

func (s *AuthService) CreateSession(userID int64, token string) error {
	return s.repoBD.CreateSession(userID, token)
}

func (s *AuthService) GenerateToken() (*models.Session, error) {
	return s.repoBD.GenerateToken()
}
