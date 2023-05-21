package usecase

import (
	"github.com/ragrag/envoi/pkg/engine"
)

type RunCode struct {
	engine *engine.Engine
}

type RunCodeParams = engine.RunParams
type RunCodeOptions = engine.RunOptions

func NewRunCode(engine *engine.Engine) *RunCode {
	return &RunCode{engine: engine}
}

func (useCase *RunCode) Execute(args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		panic("expected 1 argument")
	}

	params := args[0].(*engine.RunParams)

	result, err := useCase.engine.RunCode(params)
	if err != nil {
		return nil, err
	}

	return result, nil
}
