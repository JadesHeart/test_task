package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"test_task/internal/lib/sl"
	"time"
)

type AuthPostgres struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewAuthPostgres(db *sql.DB, logger *slog.Logger) *AuthPostgres {
	return &AuthPostgres{
		db:     db,
		logger: logger,
	}
}

const (
	FailedLoginAttempts  = 5
	findUserPath         = "repository.auth_repo.FindUser"
	checkPassPath        = "repository.auth_repo.CheckPass"
	checkAttemptsPath    = "repository.auth_repo.CheckFailedLoginAttempts"
	checkAddAttemptsPath = "repository.auth_repo.AddingFailedLoginAttempt"
	createSessionPath    = "repository.auth_repo.CreateSession"
)

func (storage *AuthPostgres) FindUser(username string) (bool, int64, error) {

	query := "SELECT EXISTS(SELECT 1 FROM Users WHERE Username = $1), UserID FROM Users WHERE Username = $1"

	var rowExist bool
	var userID int64

	err := storage.db.QueryRow(query, username).Scan(&rowExist, &userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, nil
		}
		return false, 0, fmt.Errorf("%s: %w", findUserPath, err)
	}

	return rowExist, userID, nil
}

func (storage *AuthPostgres) CheckPass(username string, password string) (bool, error) {
	query := `
		SELECT CASE WHEN PasswordHash = md5($1 || Salt) THEN true ELSE false END AS PasswordMatch
		FROM Users
		WHERE Username = $2`

	var passCorrect bool

	err := storage.db.QueryRow(query, password, username).Scan(&passCorrect)
	if err != nil {
		return false, fmt.Errorf("%s: %w", checkPassPath, err)
	}

	return passCorrect, nil
}

func (storage *AuthPostgres) CheckFailedLoginAttempts(username string) (bool, error) {

	query := `
		SELECT CASE WHEN FailedLoginAttempts >= $1 THEN TRUE ELSE FALSE END
		FROM Users
		WHERE Username = $2;
	`

	var isFailed bool

	err := storage.db.QueryRow(query, FailedLoginAttempts, username).Scan(&isFailed)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", checkAttemptsPath, err)
	}

	return isFailed, nil
}

func (storage *AuthPostgres) AddingFailedLoginAttempt(username string) error {
	query := `
       UPDATE Users
       SET FailedLoginAttempts = FailedLoginAttempts + 1
       WHERE Username = $1;
   `
	_, err := storage.db.Exec(query, username)
	if err != nil {
		return fmt.Errorf("%s: %w", checkAddAttemptsPath, err)
	}

	return nil
}

func (storage *AuthPostgres) CreateSession(userID int64, token string) error {
	currentTime := time.Now()

	insertSQL := `
        INSERT INTO Sessions (UserID, Token, ExpirationTime)
        VALUES ($1, $2, $3)
        RETURNING SessionID
    `

	var sessionID int64

	err := storage.db.QueryRow(insertSQL, userID, token, currentTime).Scan(&sessionID)
	if err != nil {
		storage.logger.Error("Failed save session", sl.Err(err))

		return fmt.Errorf("%s: %w", createSessionPath, err)
	}

	return nil
}

func (storage *AuthPostgres) GenerateToken() (string, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		storage.logger.Error("Failed generate token", sl.Err(err))

		return "", err
	}
	return newUUID.String(), nil
}
