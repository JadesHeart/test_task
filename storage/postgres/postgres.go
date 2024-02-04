package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
}

const (
	FailedLoginAttempts  = 5
	findUserPath         = "storage.postgres.FindUser"
	checkPassPath        = "storage.postgres.CheckPass"
	checkAttemptsPath    = "storage.postgres.CheckFailedLoginAttempts"
	checkAddAttemptsPath = "storage.postgres.AddingFailedLoginAttempt"
	createSessionPath    = "storage.postgres.CreateSession"
	getSessionPath       = "storage.postgres.GetSession"
)

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err)
	}
	return &Storage{DB: db}, nil
}
