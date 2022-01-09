package config

import (
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func Load() {
	// https://github.com/joho/godotenv#precedence--conventions
	v.FITZREST_ENV = os.Getenv("FITZREST_ENV")
	if len(v.FITZREST_ENV) == 0 {
		v.FITZREST_ENV = "development"
	}
	v.FITZREST_ENV_DIR = os.Getenv("FITZREST_ENV_DIR")
	tryLoad(v.FITZREST_ENV_DIR, ".env."+v.FITZREST_ENV+".local")
	if v.FITZREST_ENV != "test" {
		tryLoad(v.FITZREST_ENV_DIR, ".env.local")
	}
	tryLoad(v.FITZREST_ENV_DIR, ".env."+v.FITZREST_ENV)
	tryLoad(v.FITZREST_ENV_DIR, ".env")
	loadNonSecrets()
}

func tryLoad(path, file string) {
	f := filepath.Join(path, file)
	if err := godotenv.Load(f); err == nil {
		log.Println("Loaded config from " + f)
	}
}

func loadNonSecrets() {
	v.FITZREST_TCP_ADDRESS = parseString("FITZREST_TCP_ADDRESS", ":8080")
	v.FITZREST_TEMP_DIR = parseString("FITZREST_TEMP_DIR", "")
	v.FITZREST_MAX_FILE_SIZE = parseInt64("FITZREST_MAX_FILE_SIZE", math.MaxInt64)
	v.FITZREST_MEMORY_BUFFER_SIZE = parseInt64("FITZREST_MEMORY_BUFFER_SIZE", 1<<16)
	v.FITZREST_ENABLE_PROMETHEUS = parseBool("FITZREST_ENABLE_PROMETHEUS", false)
	v.FITZREST_ENABLE_SHUTDOWN_ENDPOINT = parseBool("FITZREST_ENABLE_SHUTDOWN_ENDPOINT", false)
	v.FITZREST_SHUTDOWN_TIMEOUT_SECONDS = parseDuration("FITZREST_SHUTDOWN_TIMEOUT_SECONDS", 0, time.Second)
	v.FITZREST_PROCESSING_MODE = parseProcessingMode("FITZREST_PROCESSING_MODE", Serialized)
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
	return time.Duration(parseInt64(os.Getenv(env), def)) * unit
}

func parseProcessingMode(env string, def ProcessingMode) ProcessingMode {
	if pm, ok := processingModeStrings[strings.ToLower(os.Getenv(env))]; ok {
		return pm
	}
	return def
}
