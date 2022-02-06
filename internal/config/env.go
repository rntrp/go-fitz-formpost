package config

import "time"

type values struct {
	FITZREST_ENV                      string
	FITZREST_ENV_DIR                  string
	FITZREST_TCP_ADDRESS              string
	FITZREST_TEMP_DIR                 string
	FITZREST_MAX_REQUEST_SIZE         int64
	FITZREST_MEMORY_BUFFER_SIZE       int64
	FITZREST_ENABLE_PROMETHEUS        bool
	FITZREST_ENABLE_SHUTDOWN_ENDPOINT bool
	FITZREST_SHUTDOWN_TIMEOUT_SECONDS time.Duration
	FITZREST_PROCESSING_MODE          ProcessingMode
}

var v values

func GetEnv() string {
	return v.FITZREST_ENV
}

func GetEnvDir() string {
	return v.FITZREST_ENV_DIR
}

func GetTCPAddress() string {
	return v.FITZREST_TCP_ADDRESS
}

func GetTempDir() string {
	return v.FITZREST_TEMP_DIR
}

func GetMaxRequestSize() int64 {
	return v.FITZREST_MAX_REQUEST_SIZE
}

func GetMemoryBufferSize() int64 {
	return v.FITZREST_MEMORY_BUFFER_SIZE
}

func IsEnablePrometheus() bool {
	return v.FITZREST_ENABLE_PROMETHEUS
}

func IsEnableShutdown() bool {
	return v.FITZREST_ENABLE_SHUTDOWN_ENDPOINT
}

func GetShutdownTimeout() time.Duration {
	return v.FITZREST_SHUTDOWN_TIMEOUT_SECONDS
}

type ProcessingMode int

const (
	Serialized ProcessingMode = iota
	Interleaved
	InMemory
)

var processingModeStrings = map[string]ProcessingMode{
	"serialized":  Serialized,
	"interleaved": Interleaved,
	"inmemory":    InMemory,
}

func GetProcessingMode() ProcessingMode {
	return v.FITZREST_PROCESSING_MODE
}
