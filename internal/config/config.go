package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func Load() {
	// https://github.com/joho/godotenv#precedence--conventions
	env := os.Getenv("FITZREST_ENV")
	if len(env) == 0 {
		env = "development"
	}
	path := os.Getenv("FITZREST_ENV_DIR")
	tryLoad(path, ".env."+env+".local")
	if env != "test" {
		tryLoad(path, ".env.local")
	}
	tryLoad(path, ".env."+env)
	tryLoad(path, ".env")
	initValues()
}

func tryLoad(path, file string) {
	f := filepath.Join(path, file)
	if err := godotenv.Load(f); err == nil {
		log.Println("Loaded config from " + f)
	}
}

func initValues() {
	v.tcpAddress = parseString("FITZREST_TCP_ADDRESS", ":8080")
	v.tempDir = parseString("FITZREST_TEMP_DIR", "")
	v.maxFileSize = parseInt64("FITZREST_MAX_FILE_SIZE", 1024*1024*512)
	v.memoryBufferSize = parseInt64("FITZREST_MEMORY_BUFFER_SIZE", 1024*64)
	v.enablePrometheus = parseBool("FITZREST_ENABLE_PROMETHEUS", true)
	v.enablePrometheus = parseBool("FITZREST_ENABLE_SHUTDOWN_ENDPOINT", false)
	v.shutdownTimeout = parseDuration("FITZREST_SHUTDOWN_TIMEOUT_SECONDS", 0, time.Second)
}

func parseBool(env string, def bool) bool {
	if b, err := strconv.ParseBool(os.Getenv(env)); err == nil {
		return b
	}
	return def
}

func parseString(env, def string) string {
	if s, ok := os.LookupEnv(env); ok {
		return s
	}
	return def
}

func parseInt64(env string, def int64) int64 {
	if i, err := strconv.ParseInt(os.Getenv(env), 10, 64); err == nil {
		return i
	}
	return def
}

func parseDuration(env string, def int64, unit time.Duration) time.Duration {
	return time.Duration(parseInt64(env, def)) * unit
}
