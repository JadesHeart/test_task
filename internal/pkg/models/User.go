package models

type User struct {
	UserID         int64
	UserName       string
	PasswordHash   string
	Salt           string
	FailedAttempts int64
}
