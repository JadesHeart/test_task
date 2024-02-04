package ping

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	resp "test_task/internal/pkg/lib/response"
	"test_task/internal/pkg/lib/sl"
	"test_task/internal/pkg/service"
)

type Request struct {
	Token string `json:"token"`
}

type Response struct {
	Pong   string `json:"pong,omitempty"`
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

func New(log *slog.Logger, services *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		op := "handlers.ping.New"

		log.With(
			slog.String("op", op),
			slog.String("request", c.GetString("requestID")),
		)

		var req Request

		//парсим json с юзером, ловим ошибку
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("failed decode request json", sl.Err(err))

			c.JSON(resp.StatusError, "failed decode request json")

			return
		}

		tokenTime, err := services.Ping.GetSession(req.Token)
		if err != nil {
			c.JSON(resp.StatusError, "Failed get token time from bd")

			return
		}

		tokenIsAlive := services.Ping.CheckTokenIsAlive(tokenTime)

		if tokenIsAlive {
			responseOK(c)

		} else {
			c.JSON(resp.StatusError, "Token is not alive")

		}
	}
}

func responseOK(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Status: resp.StatusOK,
		Pong:   "Pong",
	})
}
