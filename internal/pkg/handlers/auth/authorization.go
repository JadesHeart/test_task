package auth

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	resp "test_task/internal/pkg/lib/response"
	"test_task/internal/pkg/lib/sl"
	"test_task/internal/pkg/service"
)

type Request struct {
	Password string `json:"password"`
	Username string `json:"login" validator:"required"`
}

type Response struct {
	Token  string `json:"token,omitempty"`
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

func New(log *slog.Logger, services *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		op := "handlers.authorization.New"

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
		// валидация данных, ловим ошибку
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("failed validate", sl.Err(err))

			c.JSON(resp.StatusError, sl.Err(validateErr))

			return
		}

		user, err := services.AuthorizationBD.FindUser(req.Username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(resp.StatusError, "User does not exist")
				return
			}
			c.JSON(http.StatusInternalServerError, "Failed check user existence")
			return
		}

		isBlocked, err := services.AuthorizationBD.CheckFailedLoginAttempts(req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Failed check login attempts")

			return
		}
		if isBlocked {
			c.JSON(resp.StatusError, "User is blocked")

			return
		}

		passCorrect, err := services.AuthorizationBD.CheckPass(req.Username, req.Password)
		if err != nil {
			c.JSON(resp.StatusError, "Failed compared password")

			return
		}
		if !passCorrect {
			c.JSON(resp.StatusError, "Incorrect password")

			return
		}

		tokenModel, err := services.AuthorizationBD.GenerateToken()
		if err != nil {
			c.JSON(resp.StatusError, "Some error")

			return
		}

		err = services.AuthorizationBD.CreateSession(user.UserID, tokenModel.Token)
		if err != nil {
			c.JSON(resp.StatusError, "Failed save session")

			return
		}

		responseOK(c, tokenModel.Token)
	}
}

func responseOK(c *gin.Context, token string) {
	c.JSON(http.StatusOK, Response{
		Status: resp.StatusOK,
		Token:  token,
	})
}
