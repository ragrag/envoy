package server

import "github.com/ragrag/envoy/pkg/usecase"

type RunOptionsSchema struct {
	TimeLimit   *float64 `json:"timeLimit,omitempty"`
	MemoryLimit *int     `json:"memoryLimit,omitempty"`
}

type RunRequestSchema struct {
	Language string                  `validate:"required" json:"language"`
	Code     string                  `validate:"required" json:"code"`
	Stdin    string                  `json:"stdin"`
	Options  *usecase.RunCodeOptions `json:"options,omitempty"`
}

type JudgeRequestSchema struct {
	Language  string `validate:"required" json:"language"`
	Code      string `validate:"required" json:"code"`
	TestCases []struct {
		Input          string `validate:"required" json:"input"`
		ExpectedOutput string `validate:"required" json:"expectedOutput"`
	} `validate:"required" json:"testCases"`
	Options *usecase.JudgeCodeOptions `json:"options,omitempty"`
}
