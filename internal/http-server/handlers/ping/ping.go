package ping

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	resp "test_task/internal/lib/response"
	"test_task/internal/lib/sl"
	"time"
)

type Request struct {
	Token string `json:"token"`
}

type Response struct {
	Pong   string `json:"pong,omitempty"`
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type SessionStorage interface {
	GetSession(token string) (time.Time, error)
}

func New(log *slog.Logger, sessionStorage SessionStorage) gin.HandlerFunc {
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

		tokenTime, err := sessionStorage.GetSession(req.Token)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {

				log.Error("Token not exist", sl.Err(err))

				c.JSON(resp.StatusError, "Bad token")

				return
			} else {
				log.Error("Failed get token time from bd", sl.Err(err))

				c.JSON(resp.StatusError, "Failed get token time from bd")

				return
			}
		}

		current := time.Now()
		current = current.In(time.UTC)
		tokenTime = tokenTime.In(time.UTC)

		fmt.Println(current)
		fmt.Println(tokenTime)

		// Вычисляем разницу между временами
		duration := tokenTime.Sub(current)

		fmt.Println("Разница времени:", duration.Abs().Minutes())

		if duration.Abs().Minutes() > 5 {

			c.JSON(resp.StatusError, "Token is not alive")

		} else {

			responseOK(c)

		}
	}
}

func responseOK(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Status: resp.StatusOK,
		Pong:   "Pong",
	})
}
