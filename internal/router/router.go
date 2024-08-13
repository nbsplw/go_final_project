package router

import (
	"github.com/gofiber/fiber/v2"
	"main/internal/controllers"
)

func SetupRoutes(main *fiber.App) {
	main.Static("/", "./web")
	api := main.Group("/api")
	{
		api.Post("/task", controllers.AddTask)
		api.Get("/task", controllers.GetTask)
		api.Put("/task", controllers.UpdateTask)
		api.Delete("/task", controllers.DeleteTask)
		api.Post("/task/done", controllers.DoneTask)
		api.Get("/tasks", controllers.GetTasks)
		api.Get("/nextdate", controllers.NextDate)
	}
}
