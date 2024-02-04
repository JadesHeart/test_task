package service

import (
	"log/slog"
	"test_task/internal/pkg/repository"
	"time"
)

type PingService struct {
	repoBD repository.Ping
	logger *slog.Logger
}

func NewPingServices(repo repository.Ping, logger *slog.Logger) *PingService {
	return &PingService{
		repoBD: repo,
		logger: logger,
	}
}

func (p *PingService) GetSession(token string) (time.Time, error) {
	return p.repoBD.GetSession(token)
}

func (p *PingService) CheckTokenIsAlive(tokenTime time.Time) bool {
	return p.repoBD.CheckTokenIsAlive(tokenTime)
}
