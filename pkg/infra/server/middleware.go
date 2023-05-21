package server

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ragrag/envoi/pkg/infra"
)

func AuthMiddleware(authToken string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()

		authHeader, ok := headers["Authorization"]

		if !ok {
			return c.Status(401).Send([]byte("unauthorized"))
		}

		authHeaderSplit := strings.Split(authHeader, "Bearer ")

		if len(authHeaderSplit) != 2 {
			return c.Status(401).Send([]byte("missing bearer auth token"))
		}

		authTokenValue := authHeaderSplit[1]

		if authTokenValue != authToken {
			return c.Status(401).Send([]byte("unauthorized"))
		}

		return c.Next()
	}
}

func ErrorMiddleware(c *fiber.Ctx, err error) error {
	var statusCode int
	errCode := infra.INTERNAL_ERR
	message := err.Error()

	switch e := err.(type) {
	case *infra.ValidationError:
		statusCode = fiber.StatusUnprocessableEntity
		errCode = infra.VALIDATION_ERR
	case *fiber.Error:
		statusCode = e.Code
	default:
		statusCode = fiber.StatusInternalServerError
	}

	return c.Status(statusCode).JSON(fiber.Map{"code": errCode, "message": message})
}
