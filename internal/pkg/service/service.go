package service

import (
	"log/slog"
	"test_task/internal/pkg/repository"
	"time"
)

type AuthorizationBD interface {
	FindUser(username string) (int64, error)
	CheckPass(username string, password string) (bool, error)
	CheckFailedLoginAttempts(username string) (bool, error)
	CreateSession(userID int64, token string) error
	GenerateToken() (string, error)
}

type Ping interface {
	GetSession(token string) (time.Time, error)
	CheckTokenIsAlive(tokenTime time.Time) bool
}

type Service struct {
	AuthorizationBD
	Ping
}

func NewService(repos *repository.Repository, logger *slog.Logger) *Service {
	return &Service{
		AuthorizationBD: NewAuthServices(repos.AuthorizationBD, logger),
		Ping:            NewPingServices(repos.Ping, logger),
	}
}
