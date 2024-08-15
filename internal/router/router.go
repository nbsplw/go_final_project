package router

import (
	"github.com/gofiber/fiber/v2"
	"main/internal/controllers"
)

func SetupRoutes(main *fiber.App) {
	main.Static("/", "./web")
	api := main.Group("/api")
	{
		api.Get("/nextdate", controllers.NextDate)
		api.Post("/signin", controllers.SignIn)
		authGroup := api.Group("")
		{
			authGroup.Post("/task", controllers.AddTask)
			authGroup.Get("/task", controllers.GetTask)
			authGroup.Put("/task", controllers.UpdateTask)
			authGroup.Delete("/task", controllers.DeleteTask)
			authGroup.Post("/task/done", controllers.DoneTask)
			authGroup.Get("/tasks", controllers.GetTasks)
		}
	}
}
