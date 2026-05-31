package config

import (
	"log/slog"
	"os"
	"strings"
	"time"
	"unicode"
)

const DefaultPort = "6379"
const DefaultLogLevel = slog.LevelInfo

func GetPort() string {
	strPort, ok := os.LookupEnv("REDIX_PORT")
	if !ok {
		return DefaultPort
	}

	strPort = strings.TrimSpace(strPort)

	if strPort == "" {
		return DefaultPort
	}

	for _, c := range strPort {
		if !unicode.IsDigit(c) {
			return DefaultPort
		}
	}

	return strPort
}

func GetConnectionIdleTimeout() *time.Duration {
	strTimeout, ok := os.LookupEnv("REDIX_CONNECTION_IDLE_TIMEOUT")
	if !ok {
		return nil
	}

	timeout, err := time.ParseDuration(strTimeout)
	if err != nil {
		return nil
	}

	return &timeout
}

func GetLogLevel() slog.Level {
	strLevel, ok := os.LookupEnv("REDIX_LOG_LEVEL")
	if !ok {
		return slog.LevelInfo
	}

	switch strings.ToUpper(strings.TrimSpace(strLevel)) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
