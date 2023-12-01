package http

import (
	"context"
	"github.com/Enthreeka/auth-lab3/internal/apperror"
	"github.com/Enthreeka/auth-lab3/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func AuthMiddleware(token usecase.TokenUsecase, store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.SendStatus(fiber.StatusForbidden)
		}

		sessionToken := sess.Get("session_token")

		if sessionToken == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		t, err := token.GetToken(context.Background(), sessionToken.(string))
		if err != nil {
			if err == apperror.ErrNoFoundRows {
				return c.SendStatus(fiber.StatusUnauthorized)
			}
			return c.Status(fiber.StatusUnauthorized).JSON(err)
		}

		c.Locals("token", t)

		return c.Next()
	}
}
