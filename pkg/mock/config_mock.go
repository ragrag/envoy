package mock

import "github.com/ragrag/envoi/pkg/infra"

var ConfigMock = &infra.Config{
	LogLevel: "panic",
	Engine:   infra.EngineConfig{WorkerCount: 50, TimeLimitSeconds: 3, MemoryLimitKB: 256_000, UseControlGroups: true},
}
