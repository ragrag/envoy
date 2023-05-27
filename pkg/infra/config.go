package infra

import (
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port      int
	AuthToken string
}

type EngineConfig struct {
	WorkerCount      int
	MaxProcesses     int
	TimeLimitSeconds float64
	MemoryLimitKB    int
	UseControlGroups bool
}

type Config struct {
	LogLevel string
	Server   ServerConfig
	Engine   EngineConfig
}

var defaults = struct {
	logLevel         string
	port             int
	workerCount      int
	timeLimitSeconds float64
	memoryLimitKB    int
	useControlGroups bool
}{
	logLevel:         "info",
	port:             4000,
	workerCount:      10,
	timeLimitSeconds: 2,
	memoryLimitKB:    128_000,
	useControlGroups: true,
}

func ProvideConfig() *Config {
	viper.AutomaticEnv()

	logLevel := getString("LOG_LEVEL", defaults.logLevel)

	port := getInt("SERVER_PORT", defaults.port)
	serverAuthToken := viper.GetString("SERVER_AUTH_TOKEN")

	workerCount := getInt("ENGINE_WORKER_COUNT", defaults.workerCount)
	timeLimitSeconds := getFloat("ENGINE_TIME_LIMIT_SECONDS", defaults.timeLimitSeconds)
	memoryLimitKB := getInt("ENGINE_MEMORY_LIMIT_KB", defaults.memoryLimitKB)

	return &Config{
		LogLevel: logLevel,
		Server:   ServerConfig{Port: port, AuthToken: serverAuthToken},
		Engine:   EngineConfig{WorkerCount: workerCount, TimeLimitSeconds: timeLimitSeconds, MemoryLimitKB: memoryLimitKB, UseControlGroups: defaults.useControlGroups},
	}
}

func getString(key string, def ...string) string {
	if !viper.IsSet(key) {
		if len(def) > 0 {
			return def[0]
		}
		panic("missing config variable ")
	}
	val := viper.GetString(key)
	return val
}

func getInt(key string, def ...int) int {
	if !viper.IsSet(key) {
		if len(def) > 0 {
			return def[0]
		}
		panic("missing config variable ")
	}
	val := viper.GetInt(key)
	return val
}

func getFloat(key string, def ...float64) float64 {
	if !viper.IsSet(key) {
		if len(def) > 0 {
			return def[0]
		}
		panic("missing config variable ")
	}
	val := viper.GetFloat64(key)
	return val
}
