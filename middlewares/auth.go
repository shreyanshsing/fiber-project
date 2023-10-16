package middlewares

import (
	"errors"
	"fiber-project/config"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	TOKEN_DURATION = time.Hour * 24
	TOKEN_EXPIRED  = "token expired"
	TOKEN_INVALID  = "token invalid"
)

func verifyUserToken(token string, userId string) error {
	// verify token
	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnvoirmentVariable("SECRET")), nil
	})
	if err != nil {
		return errors.New("failed to parse token")
	}
	// check if token has expired
	if !parsedToken.Valid {
		return errors.New(TOKEN_EXPIRED)
	}
	// check if token is valid
	claimUserIdUint64 := uint64(claims["user_id"].(float64))
	userIdUint64, _ := strconv.ParseUint(userId, 10, 64)
	if claimUserIdUint64 == userIdUint64 {
		return nil
	}
	return errors.New(TOKEN_INVALID)
}

func ProtectedRoute() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Missing Authorization header",
			})
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// get user ID from request
		userId := c.Params("id")

		// verify token
		if err := verifyUserToken(token, userId); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err,
			})
		}
		// call next handler
		return c.Next()
	}
}
