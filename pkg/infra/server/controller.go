package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ragrag/envoy/pkg/engine"
	"github.com/ragrag/envoy/pkg/infra"
	"github.com/ragrag/envoy/pkg/usecase"
)

type Controller struct {
	getRuntimes     usecase.UseCase
	runCode         usecase.UseCase
	judgeCode       usecase.UseCase
	schemaValidator *SchemaValidator
}

func NewController(getRuntimes *usecase.GetRuntimes, runCode *usecase.RunCode, judgeCode *usecase.JudgeCode) *Controller {
	schemaValidator := NewSchemaValidator()
	return &Controller{getRuntimes: getRuntimes, runCode: runCode, judgeCode: judgeCode, schemaValidator: &schemaValidator}
}

func (c *Controller) GetRuntimes(ctx *fiber.Ctx) error {
	runtimes, err := c.getRuntimes.Execute()
	if err != nil {
		return err
	}

	return ctx.JSON(runtimes)
}

func (c *Controller) Run(ctx *fiber.Ctx) error {
	dto := new(RunRequestSchema)

	errors, ok := c.schemaValidator.ValidateStrict(ctx.Body(), dto)
	if !ok {
		return infra.NewValidationError(errors)
	}

	options := dto.Options

	result, err := c.runCode.Execute(&usecase.RunCodeParams{Language: dto.Language, Code: dto.Code, Stdin: dto.Stdin, Options: options})
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}

func (c *Controller) Judge(ctx *fiber.Ctx) error {
	dto := new(JudgeRequestSchema)

	errors, ok := c.schemaValidator.ValidateStrict(ctx.Body(), dto)
	if !ok {
		return infra.NewValidationError(errors)
	}

	if len(dto.TestCases) == 0 {
		return infra.NewValidationError("atleast 1 test case is required")
	}

	testCases := make([]usecase.JudgeTestCase, len(dto.TestCases))
	for i, tc := range dto.TestCases {
		testCases[i] = engine.TestCase{
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
		}
	}

	result, err := c.judgeCode.Execute(&usecase.JudgeCodeParams{Language: dto.Language, Code: dto.Code, TestCases: testCases, Options: dto.Options})
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
