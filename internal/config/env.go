package config

import "time"

type values struct {
	FITZ_FORMPOST_ENV                      string
	FITZ_FORMPOST_ENV_DIR                  string
	FITZ_FORMPOST_TCP_ADDRESS              string
	FITZ_FORMPOST_TEMP_DIR                 string
	FITZ_FORMPOST_MAX_REQUEST_SIZE         int64
	FITZ_FORMPOST_MEMORY_BUFFER_SIZE       int64
	FITZ_FORMPOST_ENABLE_PROMETHEUS        bool
	FITZ_FORMPOST_ENABLE_SHUTDOWN_ENDPOINT bool
	FITZ_FORMPOST_SHUTDOWN_TIMEOUT_SECONDS time.Duration
	FITZ_FORMPOST_PROCESSING_MODE          ProcessingMode
}

var v values

func GetEnv() string {
	return v.FITZ_FORMPOST_ENV
}

func GetEnvDir() string {
	return v.FITZ_FORMPOST_ENV_DIR
}

func GetTCPAddress() string {
	return v.FITZ_FORMPOST_TCP_ADDRESS
}

func GetTempDir() string {
	return v.FITZ_FORMPOST_TEMP_DIR
}

func GetMaxRequestSize() int64 {
	return v.FITZ_FORMPOST_MAX_REQUEST_SIZE
}

func GetMemoryBufferSize() int64 {
	return v.FITZ_FORMPOST_MEMORY_BUFFER_SIZE
}

func IsEnablePrometheus() bool {
	return v.FITZ_FORMPOST_ENABLE_PROMETHEUS
}

func IsEnableShutdown() bool {
	return v.FITZ_FORMPOST_ENABLE_SHUTDOWN_ENDPOINT
}

func GetShutdownTimeout() time.Duration {
	return v.FITZ_FORMPOST_SHUTDOWN_TIMEOUT_SECONDS
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
	return v.FITZ_FORMPOST_PROCESSING_MODE
}
