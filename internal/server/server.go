package server

import (
	"context"
	"fmt"
	"github.com/Enthreeka/auth-lab3/internal/config"
	controllerHttp "github.com/Enthreeka/auth-lab3/internal/controller/http"
	"github.com/Enthreeka/auth-lab3/internal/repo"
	"github.com/Enthreeka/auth-lab3/internal/usecase"
	"github.com/Enthreeka/auth-lab3/pkg/logger"
	"github.com/Enthreeka/auth-lab3/pkg/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
)

func Run(log *logger.Logger, cfg *config.Config) error {
	psql, err := postgres.New(context.Background(), 3, cfg.Postgres.URL)
	if err != nil {
		log.Fatal("failed to connect PostgreSQL: %v", err)
	}
	defer psql.Close()

	userRepoPG := repo.NewUserRepoPG(psql)
	tokenRepoPG := repo.NewTokenRepoPG(psql)

	userUsecase := usecase.NewUserUsecase(userRepoPG, log)
	tokenUsecase := usecase.NewTokenUsecase(tokenRepoPG, log)

	store := session.New(session.Config{
		CookieSecure:   true,
		CookieHTTPOnly: true,
	})

	userHandler := controllerHttp.NewUserHandler(userUsecase, tokenUsecase, store, log)

	engine := html.New("./template", ".html")

	engine.Debug(true)

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	//app.Static("/template/image", "./template/image")

	v1 := app.Group("/api")

	auth := v1.Group("/auth")

	auth.Use(controllerHttp.AuthMiddleware(tokenUsecase, store))
	auth.Get("/account", userHandler.AccountHandler)

	v1.Post("/login", userHandler.LogInHandler)
	v1.Get("/", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{})
	})
	v1.Post("/signup", userHandler.SignUpHandler)
	v1.Post("/logout", userHandler.LogoutHandler)

	log.Info("Starting http server: %s:%s", cfg.HTTTPServer.TypeServer, cfg.HTTTPServer.Port)
	if err = app.Listen(fmt.Sprintf(":%s", cfg.HTTTPServer.Port)); err != nil {
		log.Fatal("Server listening failed:%s", err)
	}

	return nil
}
