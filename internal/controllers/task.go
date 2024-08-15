package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"main/core/database/sqlite"
	"main/core/logger"
	"main/internal/models/common"
	"main/internal/models/tasks"
	"main/pkg"
	"strconv"
	"time"
)

func AddTask(c *fiber.Ctx) error {
	var body common.AddTask
	if err := c.BodyParser(&body); err != nil {
		logger.Get().Info("cannot parse body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Error: "incorrect request"})
	}
	if err := body.CheckTask(); err != nil {
		logger.Get().Info("internal check failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Error: err.Error()})
	}
	id, err := sqlite.Get().AddTaskDB(tasks.Task{
		Date:    body.Date,
		Title:   body.Title,
		Comment: body.Comment,
		Repeat:  body.Repeat,
	})
	if err != nil {
		logger.Get().Info("cannot add task", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Error: "cannot add task"})
	}
	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{Id: int(id)})
}

func GetTasks(c *fiber.Ctx) error {
	var offset int
	if c.Query("offset") != "" && c.Query("offset") != "" {
		if temp, err := strconv.Atoi(c.Query("offset")); err == nil {
			offset = temp
		}
	}
	if c.Query("search") != "" {
		if parsedDate, err := time.Parse("02.1.2006", c.Query("search")); err == nil {
			resultTasks, err := sqlite.Get().TasksByDate(parsedDate.Format(common.TimeFormat))
			if err != nil {
				logger.Get().Info("cannot get tasks", zap.Error(err))
				return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Error: "cannot get tasks"})
			}
			return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{Tasks: resultTasks})
		}
		resultTasks, err := sqlite.Get().SearchTasks(c.Query("search"), offset)
		if err != nil {
			return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{Tasks: resultTasks})
		}
	}
	resultTasks, err := sqlite.Get().Tasks(offset)
	if err != nil {
		logger.Get().Info("cannot get resultTasks", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Error: "cannot get resultTasks"})
	}
	if resultTasks == nil {
		resultTasks = []tasks.Task{}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"tasks": resultTasks})
}

func GetTask(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Error: "id required"})
	}
	task, err := sqlite.Get().FindTask(id)
	if err != nil {
		if errors.Is(err, sqlite.ErrNoSuchTask) {
			logger.Get().Info("cannot find task", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Error: "cannot find task"})
		}
		logger.Get().Error("cannot get task", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Error: "cannot get task"})
	}
	if task == nil {
		task = &tasks.Task{}
	}
	return c.Status(fiber.StatusOK).JSON(task)
}

func UpdateTask(c *fiber.Ctx) error {
	var body tasks.Task
	if err := c.BodyParser(&body); err != nil {
		logger.Get().Info("cannot parse body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Error: "incorrect request"})
	}
	if body.ID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Error: "id required"})
	}
	if _, err := strconv.Atoi(body.ID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Error: "id required"})
	}
	validate := common.AddTask{
		Date:    body.Date,
		Title:   body.Title,
		Comment: body.Comment,
		Repeat:  body.Repeat,
	}
	if err := validate.CheckTask(); err != nil {
		logger.Get().Info("internal check failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Error: err.Error()})
	}
	if _, err := sqlite.Get().FindTask(body.ID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Error: "cannot find task"})
	}
	if err := sqlite.Get().UpdateTask(body); err != nil {
		if errors.Is(err, sqlite.ErrNoSuchTask) {
			logger.Get().Info("no such task", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Error: "no such task"})
		}
		logger.Get().Error("cannot update task", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Error: "cannot update task"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func DoneTask(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Error: "invalid id"})
	}
	if err := sqlite.Get().DoneTask(id); err != nil {
		if errors.Is(err, sqlite.ErrNoSuchTask) {
			logger.Get().Info("no such task", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Error: "no such task"})
		}
		logger.Get().Error("cannot done task", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Error: "cannot done task"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func DeleteTask(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Error: "invalid id"})
	}
	if err := sqlite.Get().DeleteTask(id); err != nil {
		if errors.Is(err, sqlite.ErrNoSuchTask) {
			logger.Get().Info("no such task", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Error: "no such task"})
		}
		logger.Get().Error("cannot delete task", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Error: "cannot delete task"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func NextDate(c *fiber.Ctx) error {
	now, err := time.Parse(common.TimeFormat, c.Query("now"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Error: err.Error()})
	}
	next, err := pkg.NextDate(now, c.Query("date"), c.Query("repeat"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Error: err.Error()})
	}
	return c.SendString(next)
}
