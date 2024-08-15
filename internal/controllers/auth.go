package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"main/core/config"
	"main/core/logger"
	"main/core/middleware"
	"main/internal/models/common"
)

func SignIn(c *fiber.Ctx) error {
	var password struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.BodyParser(&password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Error: "invalid request"})
	}
	if password.Password != config.Get().Auth.Password || password.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Error: "incorrect password"})
	}

	token, err := middleware.GenerateToken(password.Password)
	if err != nil {
		logger.Get().Error("failed to generate token", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Error: "failed to generate token"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}
