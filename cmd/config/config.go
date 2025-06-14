package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/meysam81/x/logging"
)

type Config struct {
	LogLevel   string
	Retries    uint
	Timeout    uint
	StatusCode uint
	HTTPTarget string
}

func New() *Config {
	logLevel := "info"
	if level, ok := os.LookupEnv("LOG_LEVEL"); ok {
		logLevel = level
	}

	return &Config{
		LogLevel:   logLevel,
		Timeout:    5,
		StatusCode: 200,
	}
}

func (c *Config) Validate() error {
	validLogLevels := []string{"debug", "info", "warn", "error", "critical"}
	logLevel := strings.ToLower(c.LogLevel)

	for _, level := range validLogLevels {
		if level == logLevel {
			return nil
		}
	}

	return fmt.Errorf("invalid log level %q, valid levels: %s",
		c.LogLevel, strings.Join(validLogLevels, ", "))
}

func (c *Config) GetLogLevel() logging.LogLevel {
	logLevels := map[string]logging.LogLevel{
		"debug":    logging.DEBUG,
		"info":     logging.INFO,
		"warn":     logging.WARN,
		"error":    logging.ERROR,
		"critical": logging.CRITICAL,
	}

	if level, ok := logLevels[strings.ToLower(c.LogLevel)]; ok {
		return level
	}

	return logging.INFO
}
