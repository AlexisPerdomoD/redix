package config

import (
	"os"
	"strings"
	"time"
	"unicode"
)

const DefaultPort = "6379"

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
