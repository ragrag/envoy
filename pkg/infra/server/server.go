package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ragrag/envoi/pkg/infra"
)

func ProvideServer(cfg *infra.Config, controller *Controller) *fiber.App {
	server := fiber.New(fiber.Config{
		ErrorHandler: ErrorMiddleware,
	})

	server.Use(cors.New())

	if cfg.Server.AuthToken != "" {
		server.Use(AuthMiddleware(cfg.Server.AuthToken))
	}

	RegisterRoutes(server, controller)

	return server
}
