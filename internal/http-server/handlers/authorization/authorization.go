package authorization

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	resp "test_task/internal/lib/response"
	"test_task/internal/lib/sl"
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

type Authorizer interface {
	FindUser(username string) (bool, int64, error)
	CheckPass(username string, password string) (bool, error)
	CheckFailedLoginAttempts(username string) (bool, error)
	AddingFailedLoginAttempt(username string) error
	CreateSession(userID int64, token string) error
}

func New(log *slog.Logger, authorizer Authorizer) gin.HandlerFunc {
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

		userExist, userID, err := authorizer.FindUser(req.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Error("User not exist", sl.Err(err))

				c.JSON(resp.StatusError, "User not exist")

				return
			}
			log.Error("Failed check user existence", sl.Err(err))
			c.JSON(resp.StatusError, "Failed check user existence")
			return
		}

		if !userExist {
			c.JSON(resp.StatusError, "User does not exist")

			return
		}

		userBlock, err := authorizer.CheckFailedLoginAttempts(req.Username)
		if err != nil {
			log.Error("Failed check login attempts", sl.Err(err))

			c.JSON(resp.StatusError, "Failed check login attempts")

			return
		}
		if userBlock {
			c.JSON(resp.StatusError, "User block")

			return
		}

		passCorrect, err := authorizer.CheckPass(req.Username, req.Password)
		if err != nil {
			log.Error("Failed check user existence", sl.Err(err))

			c.JSON(resp.StatusError, "Failed compared password")

			return
		}
		if !passCorrect {
			err = authorizer.AddingFailedLoginAttempt(req.Username)
			if err != nil {
				log.Error("Failed add login attempt", sl.Err(err))

				c.JSON(resp.StatusError, "Some error")

				return
			}

			c.JSON(resp.StatusError, "Incorrect password")

			return
		}

		token, err := generateToken()
		if err != nil {
			log.Error("Failed generate token", sl.Err(err))

			c.JSON(resp.StatusError, "Some error")

			return
		}

		fmt.Println(token)

		err = authorizer.CreateSession(userID, token)
		if err != nil {
			log.Error("Failed save session", sl.Err(err))

			c.JSON(resp.StatusError, "Failed save session")

			return
		}

		responseOK(c, token)
	}
}

func responseOK(c *gin.Context, token string) {
	c.JSON(http.StatusOK, Response{
		Status: resp.StatusOK,
		Token:  token,
	})
}

func generateToken() (string, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return newUUID.String(), nil
}
