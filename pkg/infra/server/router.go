package server

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(server *fiber.App, controller *Controller) {
	server.Get("/healthz", func(ctx *fiber.Ctx) error {
		return ctx.SendString("OK")
	})
	server.Get("/runtimes", controller.GetRuntimes)
	server.Post("/run", controller.Run)
	server.Post("/judge", controller.Judge)
}
