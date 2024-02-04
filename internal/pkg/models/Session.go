package models

import "time"

type Session struct {
	SessionID      int64
	UserID         int64
	Token          string
	ExpirationTime time.Time
}
