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
	loadDotEnv()
	tryLoad(v.FITZ_FORMPOST_ENV_DIR, ".env."+v.FITZ_FORMPOST_ENV+".local")
	if v.FITZ_FORMPOST_ENV != "test" {
		tryLoad(v.FITZ_FORMPOST_ENV_DIR, ".env.local")
	}
	tryLoad(v.FITZ_FORMPOST_ENV_DIR, ".env."+v.FITZ_FORMPOST_ENV)
	tryLoad(v.FITZ_FORMPOST_ENV_DIR, ".env")
	loadEnv()
	v.print()
}

func tryLoad(path, file string) {
	f := filepath.Join(path, file)
	if err := godotenv.Load(f); err == nil {
		log.Println("Loaded config from " + f)
	}
}

func loadDotEnv() {
	v.FITZ_FORMPOST_ENV = os.Getenv("FITZ_FORMPOST_ENV")
	if len(v.FITZ_FORMPOST_ENV) == 0 {
		v.FITZ_FORMPOST_ENV = "development"
	}
	v.FITZ_FORMPOST_ENV_DIR = os.Getenv("FITZ_FORMPOST_ENV_DIR")
}

func loadEnv() {
	v.FITZ_FORMPOST_TCP_ADDRESS = parseString("FITZ_FORMPOST_TCP_ADDRESS", ":8080")
	v.FITZ_FORMPOST_TEMP_DIR = parseString("FITZ_FORMPOST_TEMP_DIR", os.TempDir())
	v.FITZ_FORMPOST_MAX_REQUEST_SIZE = parseInt64("FITZ_FORMPOST_MAX_REQUEST_SIZE", -1)
	v.FITZ_FORMPOST_MEMORY_BUFFER_SIZE = parseInt64("FITZ_FORMPOST_MEMORY_BUFFER_SIZE", 10<<20)
	v.FITZ_FORMPOST_ENABLE_PROMETHEUS = parseBool("FITZ_FORMPOST_ENABLE_PROMETHEUS", false)
	v.FITZ_FORMPOST_ENABLE_SHUTDOWN_ENDPOINT = parseBool("FITZ_FORMPOST_ENABLE_SHUTDOWN_ENDPOINT", false)
	v.FITZ_FORMPOST_SHUTDOWN_TIMEOUT = parseDuration("FITZ_FORMPOST_SHUTDOWN_TIMEOUT", 0)
	v.FITZ_FORMPOST_PROCESSING_MODE = parseProcessingMode("FITZ_FORMPOST_PROCESSING_MODE", Serialized)
	v.FITZ_FORMPOST_RENDERING_DPI = parseDPI("FITZ_FORMPOST_RENDERING_DPI", 300.0)
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

func parseFloat64(env string, def float64) float64 {
	if f, err := strconv.ParseFloat(os.Getenv(env), 64); err == nil {
		return f
	}
	return def
}

func parseDuration(env string, def time.Duration) time.Duration {
	if d, err := time.ParseDuration(os.Getenv(env)); err == nil {
		return d
	}
	return def
}

func parseProcessingMode(env string, def ProcessingMode) ProcessingMode {
	if pm, ok := processingModeStrings[strings.ToLower(os.Getenv(env))]; ok {
		return pm
	}
	return def
}

func parseDPI(env string, def float64) float64 {
	f := parseFloat64(env, def)
	if f < 1 || math.IsInf(f, 0) || math.IsNaN(f) {
		return def
	}
	return f
}
