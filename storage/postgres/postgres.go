package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type Storage struct {
	db *sql.DB
}

const FailedLoginAttempts = 5

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err)
	}
	return &Storage{db: db}, nil
}

func (storage *Storage) FindUser(username string) (bool, int64, error) {
	const fn = "storage.postgres.FindUser"

	query := "SELECT EXISTS(SELECT 1 FROM Users WHERE Username = $1), UserID FROM Users WHERE Username = $1"

	var rowExist bool
	var userID int64

	err := storage.db.QueryRow(query, username).Scan(&rowExist, &userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, nil
		}
		return false, 0, fmt.Errorf("%s: %w", fn, err)
	}

	return rowExist, userID, nil
}

func (storage *Storage) CheckPass(username string, password string) (bool, error) {
	const fn = "storage.postgres.CheckPass"

	query := `
		SELECT CASE WHEN PasswordHash = md5($1 || Salt) THEN true ELSE false END AS PasswordMatch
		FROM Users
		WHERE Username = $2`

	var passCorrect bool

	err := storage.db.QueryRow(query, password, username).Scan(&passCorrect)
	if err != nil {
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	return passCorrect, nil
}

func (storage *Storage) CheckFailedLoginAttempts(username string) (bool, error) {
	const fn = "storage.postgres.CheckFailedLoginAttempts"

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
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	return isFailed, nil
}

func (storage *Storage) AddingFailedLoginAttempt(username string) error {
	const fn = "storage.postgres.CheckFailedLoginAttempts"

	query := `
        UPDATE Users
        SET FailedLoginAttempts = FailedLoginAttempts + 1
        WHERE Username = $1;
    `
	_, err := storage.db.Exec(query, username)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (storage *Storage) CreateSession(userID int64, token string) error {
	const fn = "storage.postgres.CreateSession"

	currentTime := time.Now()

	insertSQL := `
        INSERT INTO Sessions (UserID, Token, ExpirationTime)
        VALUES ($1, $2, $3)
        RETURNING SessionID
    `

	var sessionID int64

	err := storage.db.QueryRow(insertSQL, userID, token, currentTime).Scan(&sessionID)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (storage *Storage) GetSession(token string) (time.Time, error) {
	const fn = "storage.postgres.GetSession"

	query := `
		SELECT ExpirationTime FROM Sessions
		WHERE Token = $1
	`

	var expirationTime time.Time
	err := storage.db.QueryRow(query, token).Scan(&expirationTime)

	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, fmt.Errorf("session with token %s not found", token)
		} else {
			return time.Time{}, fmt.Errorf("error getting session: %v", err)
		}
	}

	return expirationTime, nil
}
