package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ragrag/envoi/pkg/infra"
	"github.com/sirupsen/logrus"
)

func ProvideServer(cfg *infra.Config, logger *logrus.Logger, controller *Controller) *fiber.App {
	server := fiber.New(fiber.Config{
		ErrorHandler: ErrorMiddleware,
	})

	server.Use(cors.New())

	server.Use(func(c *fiber.Ctx) error {
		if c.Path() != "/healthz" {
			logger.WithFields(logrus.Fields{
				"method": c.Method(),
				"path":   c.Path(),
			}).Info("Incoming request")
		}
		err := c.Next()

		return err
	})

	if cfg.Server.AuthToken != "" {
		server.Use(AuthMiddleware(cfg.Server.AuthToken))
	}

	RegisterRoutes(server, controller)

	return server
}
