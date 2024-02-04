package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"test_task/internal/lib/sl"
	"time"
)

type PingPostgres struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPingPostgres(db *sql.DB, logger *slog.Logger) *PingPostgres {
	return &PingPostgres{
		db:     db,
		logger: logger,
	}
}

const (
	getSessionPath = "repository.ping_repo.GetSession"
)

func (p *PingPostgres) GetSession(token string) (time.Time, error) {

	query := `
		SELECT ExpirationTime FROM Sessions
		WHERE Token = $1
	`

	var expirationTime time.Time
	err := p.db.QueryRow(query, token).Scan(&expirationTime)

	if err != nil {
		p.logger.Error("Failed get token time from bd", sl.Err(err))

		return time.Time{}, fmt.Errorf("%s: %w", getSessionPath, err)
	}

	return expirationTime, nil
}

func (p *PingPostgres) CheckTokenIsAlive(tokenTime time.Time) bool {
	current := time.Now()
	current = current.In(time.UTC)
	tokenTime = tokenTime.In(time.UTC)
	duration := tokenTime.Sub(current)

	if duration.Abs().Minutes() > 5 {
		return false
	} else {
		return true
	}

}
