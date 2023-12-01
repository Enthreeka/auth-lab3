package http

import (
	"context"
	"errors"
	"github.com/Enthreeka/auth-lab3/internal/apperror"
	"github.com/Enthreeka/auth-lab3/internal/entity"
	"github.com/Enthreeka/auth-lab3/internal/usecase"
	"github.com/Enthreeka/auth-lab3/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"time"
)

type userHandler struct {
	userUsecase  usecase.UserUsecase
	tokenUsecase usecase.TokenUsecase

	store *session.Store
	log   *logger.Logger
}

func NewUserHandler(userUsecase usecase.UserUsecase, tokenUsecase usecase.TokenUsecase, store *session.Store, log *logger.Logger) *userHandler {
	return &userHandler{
		userUsecase:  userUsecase,
		tokenUsecase: tokenUsecase,
		store:        store,
		log:          log,
	}
}

func (u *userHandler) LogInHandler(c *fiber.Ctx) error {
	u.log.Info("Starting LogInHandler")

	login := c.FormValue("login")
	password := c.FormValue("password")
	u.log.Info("%s,%s", login, password)

	user, err := u.userUsecase.GetUserByLogin(context.Background(), login, password)
	if err != nil {
		u.log.Error("%v", err)
		if errors.Is(err, apperror.ErrNoFoundRows) {
			return c.Status(fiber.StatusNotFound).JSON(apperror.ErrNoFoundRows)
		} else if errors.Is(err, apperror.ErrHashPasswordsNotEqual) {
			return c.Status(fiber.StatusUnauthorized).JSON(apperror.ErrHashPasswordsNotEqual)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(err)
	}

	newToken := uuid.New().String()
	status := updateToken(c, u.store, newToken)
	if status != 0 {
		return c.SendStatus(status)
	}

	// Желательно сделать функцию для обновления старого токена
	err = u.tokenUsecase.CreateToken(context.Background(), newToken, user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	u.log.Info("Login successfully")
	return c.Redirect("/api/auth/account", fiber.StatusSeeOther)
}

func (u *userHandler) SignUpHandler(c *fiber.Ctx) error {
	u.log.Info("starting SignUpHandler")

	login := c.FormValue("login")
	password := c.FormValue("password")
	u.log.Info("%s,%s", login, password)

	user := &entity.User{
		Login:    login,
		Password: password,
	}

	userID := uuid.New().String()
	err := u.userUsecase.CreateUser(context.Background(), user, userID)
	if err != nil {
		u.log.Error("%v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	token := uuid.New().String()
	status := updateToken(c, u.store, token)
	if status != 0 {
		return c.SendStatus(status)
	}

	err = u.tokenUsecase.CreateToken(context.Background(), token, userID)
	if err != nil {
		u.log.Error("%v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	u.log.Info("SignUpHandler successfully")
	return c.Redirect("/api/auth/account", fiber.StatusSeeOther)
}

func (u *userHandler) AccountHandler(c *fiber.Ctx) error {
	u.log.Info("starting AccountHandler")

	t := c.Locals("token")

	token := t.(*entity.Token)
	user, err := u.userUsecase.GetUserByID(context.Background(), token.UserID)
	if err != nil {
		u.log.Error("%v", err)
		if errors.Is(err, apperror.ErrNoFoundRows) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	u.log.Info("AccountHandler successfully")

	return c.Render("account", fiber.Map{
		"user": user,
	})
}

func (u *userHandler) LogoutHandler(c *fiber.Ctx) error {
	u.log.Info("starting LogoutHandler")

	sess, err := u.store.Get(c)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	err = sess.Destroy()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	u.log.Info("LogoutHandler successfully")
	return c.Redirect("/api", fiber.StatusSeeOther)
}

func updateToken(c *fiber.Ctx, store *session.Store, token string) int {
	sess, err := store.Get(c)
	if err != nil {
		return fiber.StatusForbidden
	}

	sess.Set("session_token", token)
	sess.SetExpiry(1 * time.Hour)

	err = sess.Save()
	if err != nil {
		return fiber.StatusInternalServerError
	}

	return 0
}
