package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ragrag/envoi/pkg/engine"
	"github.com/ragrag/envoi/pkg/infra"
	"github.com/ragrag/envoi/pkg/infra/server"
	"github.com/ragrag/envoi/pkg/runtimeenv"
	"github.com/ragrag/envoi/pkg/sandbox"
	"github.com/ragrag/envoi/pkg/usecase"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

func main() {
	c := dig.New()

	c.Provide(infra.ProvideConfig)
	c.Provide(infra.NewLogger)

	c.Provide(runtimeenv.NewRuntimeProvider)
	c.Provide(sandbox.NewSanboxManager)
	c.Provide(engine.NewEngine)

	c.Provide(usecase.NewGetRuntimes)
	c.Provide(usecase.NewRunCode)
	c.Provide(usecase.NewJudgeCode)

	c.Provide(server.NewController)
	c.Provide(server.ProvideServer)

	err := c.Invoke(func(logger *logrus.Logger, config *infra.Config, runtimeProvider *runtimeenv.RuntimeProvider, engine *engine.Engine, server *fiber.App) error {
		e := runtimeProvider.Load()
		if e != nil {
			return e
		}

		e = engine.Ignite()
		if e != nil {
			return e
		}

		e = server.Listen(fmt.Sprintf(":%d", config.Server.Port))
		if e != nil {
			return e
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
