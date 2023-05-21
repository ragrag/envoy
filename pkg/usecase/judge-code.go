package usecase

import (
	"github.com/ragrag/envoi/pkg/engine"
)

type JudgeCode struct {
	engine *engine.Engine
}

type JudgeTestCase = engine.TestCase
type JudgeCodeParams = engine.JudgeParams
type JudgeCodeOptions = engine.RunOptions

func NewJudgeCode(engine *engine.Engine) *JudgeCode {
	return &JudgeCode{engine: engine}
}

func (useCase *JudgeCode) Execute(args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		panic("expected 1 argument")
	}

	params := args[0].(*engine.JudgeParams)

	result, err := useCase.engine.JudgeCode(params)
	if err != nil {
		return nil, err
	}

	return result, nil
}
