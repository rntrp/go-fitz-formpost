package config

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

type values struct {
	FITZ_FORMPOST_ENV                      string
	FITZ_FORMPOST_ENV_DIR                  string
	FITZ_FORMPOST_TCP_ADDRESS              string
	FITZ_FORMPOST_TEMP_DIR                 string
	FITZ_FORMPOST_MAX_REQUEST_SIZE         int64
	FITZ_FORMPOST_MEMORY_BUFFER_SIZE       int64
	FITZ_FORMPOST_ENABLE_PROMETHEUS        bool
	FITZ_FORMPOST_ENABLE_SHUTDOWN_ENDPOINT bool
	FITZ_FORMPOST_SHUTDOWN_TIMEOUT         time.Duration
	FITZ_FORMPOST_PROCESSING_MODE          ProcessingMode
	FITZ_FORMPOST_RENDERING_DPI            float64
}

func (v *values) print() {
	buf := new(strings.Builder)
	buf.WriteString("Environment has been resolved to:\n")
	val := reflect.Indirect(reflect.ValueOf(v))
	valType := val.Type()
	for i := range val.NumField() {
		a := valType.Field(i).Name
		b := val.Field(i).Interface()
		fmt.Fprintf(buf, "%-40s= %v\n", a, b)
	}
	log.Print(buf.String())
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
	return v.FITZ_FORMPOST_SHUTDOWN_TIMEOUT
}

type ProcessingMode int

const (
	Serialized ProcessingMode = iota
	Interleaved
	InMemory
)

func (pm ProcessingMode) String() string {
	switch pm {
	case Serialized:
		return "serialized"
	case Interleaved:
		return "interleaved"
	case InMemory:
		return "inmemory"
	default:
		return ""
	}
}

var processingModeStrings = map[string]ProcessingMode{
	Serialized.String():  Serialized,
	Interleaved.String(): Interleaved,
	InMemory.String():    InMemory,
}

func GetProcessingMode() ProcessingMode {
	return v.FITZ_FORMPOST_PROCESSING_MODE
}

func GetRenderingDPI() float64 {
	return v.FITZ_FORMPOST_RENDERING_DPI
}
