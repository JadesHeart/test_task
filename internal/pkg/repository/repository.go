package repository

import (
	"log/slog"
	"test_task/storage/postgres"
	"time"
)

type AuthorizationBD interface {
	FindUser(username string) (bool, int64, error)
	CheckPass(username string, password string) (bool, error)
	CheckFailedLoginAttempts(username string) (bool, error)
	AddingFailedLoginAttempt(username string) error
	CreateSession(userID int64, token string) error
	GenerateToken() (string, error)
}

type Ping interface {
	GetSession(token string) (time.Time, error)
	CheckTokenIsAlive(tokenTime time.Time) bool
}

type Repository struct {
	AuthorizationBD
	Ping
}

func NewRepository(storage *postgres.Storage, logger *slog.Logger) *Repository {
	return &Repository{
		AuthorizationBD: NewAuthPostgres(storage.DB, logger),
		Ping:            NewPingPostgres(storage.DB, logger),
	}
}
