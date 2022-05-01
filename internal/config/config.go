package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func Load() {
	// https://github.com/joho/godotenv#precedence--conventions
	v.FITZ_FORMPOST_ENV = os.Getenv("FITZ_FORMPOST_ENV")
	if len(v.FITZ_FORMPOST_ENV) == 0 {
		v.FITZ_FORMPOST_ENV = "development"
	}
	v.FITZ_FORMPOST_ENV_DIR = os.Getenv("FITZ_FORMPOST_ENV_DIR")
	tryLoad(v.FITZ_FORMPOST_ENV_DIR, ".env."+v.FITZ_FORMPOST_ENV+".local")
	if v.FITZ_FORMPOST_ENV != "test" {
		tryLoad(v.FITZ_FORMPOST_ENV_DIR, ".env.local")
	}
	tryLoad(v.FITZ_FORMPOST_ENV_DIR, ".env."+v.FITZ_FORMPOST_ENV)
	tryLoad(v.FITZ_FORMPOST_ENV_DIR, ".env")
	loadNonSecrets()
	if j, err := json.MarshalIndent(&v, "", "  "); err == nil {
		log.Println("Environment: " + string(j))
	} else {
		log.Panicln(err)
	}
}

func tryLoad(path, file string) {
	f := filepath.Join(path, file)
	if err := godotenv.Load(f); err == nil {
		log.Println("Loaded config from " + f)
	}
}

func loadNonSecrets() {
	v.FITZ_FORMPOST_TCP_ADDRESS = parseString("FITZ_FORMPOST_TCP_ADDRESS", ":8080")
	v.FITZ_FORMPOST_TEMP_DIR = parseString("FITZ_FORMPOST_TEMP_DIR", "")
	v.FITZ_FORMPOST_MAX_REQUEST_SIZE = parseInt64("FITZ_FORMPOST_MAX_REQUEST_SIZE", -1)
	v.FITZ_FORMPOST_MEMORY_BUFFER_SIZE = parseInt64("FITZ_FORMPOST_MEMORY_BUFFER_SIZE", 10<<20)
	v.FITZ_FORMPOST_ENABLE_PROMETHEUS = parseBool("FITZ_FORMPOST_ENABLE_PROMETHEUS", false)
	v.FITZ_FORMPOST_ENABLE_SHUTDOWN_ENDPOINT = parseBool("FITZ_FORMPOST_ENABLE_SHUTDOWN_ENDPOINT", false)
	v.FITZ_FORMPOST_SHUTDOWN_TIMEOUT_SECONDS = parseDuration("FITZ_FORMPOST_SHUTDOWN_TIMEOUT_SECONDS", 0, time.Second)
	v.FITZ_FORMPOST_PROCESSING_MODE = parseProcessingMode("FITZ_FORMPOST_PROCESSING_MODE", Serialized)
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
