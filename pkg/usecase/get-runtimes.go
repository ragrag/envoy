package usecase

import (
	"github.com/ragrag/envoi/pkg/runtimeenv"
)

type GetRuntimes struct {
	runtimeProvider *runtimeenv.RuntimeProvider
}

type GetRuntimesResult struct {
	ID       string `json:"id"`
	Language string `json:"language"`
	Version  string `json:"version"`
}

func NewGetRuntimes(runtimeProvider *runtimeenv.RuntimeProvider) *GetRuntimes {
	return &GetRuntimes{runtimeProvider: runtimeProvider}
}

func (useCase *GetRuntimes) Execute(args ...interface{}) (interface{}, error) {
	availableRuntimes := useCase.runtimeProvider.GetAvailableRuntimes()
	var jsonAvailableRuntimes []GetRuntimesResult

	for _, runtime := range availableRuntimes {
		jsonAvailableRuntimes = append(jsonAvailableRuntimes, GetRuntimesResult{ID: runtime.ID, Language: runtime.Meta.Language, Version: runtime.Meta.Version})
	}

	return jsonAvailableRuntimes, nil
}
