package server

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(server *fiber.App, controller *Controller) {
	server.Get("/", controller.Run)
	server.Get("/runtimes", controller.GetRuntimes)
	server.Post("/run", controller.Run)
	server.Post("/judge", controller.Judge)
}
