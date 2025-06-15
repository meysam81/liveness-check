package config

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/meysam81/x/logging"
)

type StaticHTTPTarget struct {
	HTTPTarget string
}

type K8sPodTarget struct {
	LabelSelectors []string
	Namespace      string
	Scheme         string
	TLSVerify      bool
	Endpoint       string
	Port           int32
}

type Config struct {
	LogLevel   string
	Retries    uint
	Timeout    uint
	StatusCode int

	StaticHTTPTarget *StaticHTTPTarget
	K8sPodTarget     *K8sPodTarget
}

func New() *Config {
	logLevel := "info"
	if level, ok := os.LookupEnv("LOG_LEVEL"); ok {
		logLevel = strings.ToLower(level)
	}

	return &Config{
		LogLevel:         logLevel,
		Timeout:          5,
		StatusCode:       200,
		StaticHTTPTarget: &StaticHTTPTarget{},
		K8sPodTarget:     &K8sPodTarget{},
	}
}

func (c *Config) Validate() error {
	validLogLevels := []string{"debug", "info", "warn", "error", "critical"}
	validSchemes := []string{"http", "https"}
	logLevel := strings.ToLower(c.LogLevel)

	errs := []string{}

	if !slices.Contains(validLogLevels, logLevel) {
		errs = append(errs, fmt.Sprintf("invalid log level %q, valid levels: %s",
			c.LogLevel, strings.Join(validLogLevels, ", ")))
	}

	if len(c.K8sPodTarget.LabelSelectors) == 0 {
		errs = append(errs, "label selector cannot be empty. provide key=value pairs to fix this.")
	} else {
		for _, pair := range c.K8sPodTarget.LabelSelectors {
			keyValue := strings.Split(pair, "=")
			if len(keyValue) != 2 {
				errs = append(errs, fmt.Sprintf("invalid label selector %s. provide selector in key=value format", keyValue))
			}
		}
	}

	if !slices.Contains(validSchemes, c.K8sPodTarget.Scheme) {
		errs = append(errs, fmt.Sprintf("invalid scheme provided: %s. only http & https are supported.", c.K8sPodTarget.Scheme))
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s", fmt.Sprint(strings.Join(errs, "\n")))
	}

	return nil
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
