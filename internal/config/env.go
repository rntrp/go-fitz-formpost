package config

import "time"

type values struct {
	tcpAddress       string
	tempDir          string
	maxFileSize      int64
	memoryBufferSize int64
	enablePrometheus bool
	enableShutdown   bool
	shutdownTimeout  time.Duration
}

var v values

func GetTCPAddress() string {
	return v.tcpAddress
}

func GetTempDir() string {
	return v.tempDir
}

func GetMaxFileSize() int64 {
	return v.maxFileSize
}

func GetMemoryBufferSize() int64 {
	return v.memoryBufferSize
}

func IsEnablePrometheus() bool {
	return v.enablePrometheus
}

func IsEnableShutdown() bool {
	return v.enableShutdown
}

func GetShutdownTimeout() time.Duration {
	return v.shutdownTimeout
}
