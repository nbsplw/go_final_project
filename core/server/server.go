package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
	"main/core/config"
	"main/core/logger"
	"main/internal/router"
	"time"
)

func Run() {
	app := fiber.New(
		fiber.Config{
			BodyLimit:             config.Get().Fiber.BodyLimit,
			IdleTimeout:           time.Duration(config.Get().Fiber.IdleTimeOut) * time.Millisecond,
			StrictRouting:         true,
			DisableStartupMessage: true,
		})
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.Get().Fiber.AllowOrigins,
		AllowMethods: config.Get().Fiber.AllowMethods,
		AllowHeaders: config.Get().Fiber.AllowHeaders,
	}))

	router.SetupRoutes(app)
	if err := app.Listen(config.Get().Listen.BindIP + ":" + config.Get().Listen.Port); err != nil {
		logger.Get().Fatal(`app listen failed`, zap.Error(err))
		return
	}
}
