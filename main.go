package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/meysam81/x/logging"
)

var (
	version = "dev"
	commit  = "HEAD"
	date    = "unknown"
	builtBy = "unknown"
)

type AppState struct {
	l *logging.Logger
}

type Config struct {
	logLevel string
}

type HttpCheck struct {
	upstream   string
	timeout    uint
	retries    uint
	statusCode uint
}

func translateLogLevel(level string) logging.LogLevel {
	logLevels := map[string]logging.LogLevel{
		"debug":    logging.DEBUG,
		"info":     logging.INFO,
		"warn":     logging.WARN,
		"error":    logging.ERROR,
		"critical": logging.CRITICAL,
	}

	logLevel := logging.INFO

	if mappedLevel, ok := logLevels[strings.ToLower(level)]; ok {
		logLevel = mappedLevel
	}

	return logLevel

}

func randomBetween(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func newLogger(c *Config) *logging.Logger {
	logger := logging.NewLogger(logging.WithLogLevel(translateLogLevel(c.logLevel)))
	return &logger
}

func (a *AppState) performHttpCheck(check *HttpCheck) error {
	client := &http.Client{
		Timeout: time.Duration(check.timeout) * time.Second,
	}

	var tries uint = 0
	for {
		if tries >= check.retries {
			a.l.Error().Msgf("halting check with max retries reached: %d", check.retries)
			return nil
		}

		resp, err := client.Get(check.upstream)
		if err != nil {
			tries += 1
			jitter := randomBetween(5, 10)
			a.l.Info().Err(err).Msgf("[%d] upstream check unsuccessful, retrying in %ds", tries, jitter)
			time.Sleep(time.Duration(jitter) * time.Second)
			continue
		}

		if resp.StatusCode == int(check.statusCode) {
			a.l.Info().Msgf("http check successful with status code: %s", resp.Status)
			break
		}
	}

	return nil
}

func createVersionCommand() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "show version information",
		Action: func(ctx context.Context, c *cli.Command) error {
			fmt.Printf("Version:    %s\n", version)
			fmt.Printf("Commit:     %s\n", commit)
			fmt.Printf("Built:      %s\n", date)
			fmt.Printf("Built by:   %s\n", builtBy)
			fmt.Printf("Go version: %s\n", runtime.Version())
			fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
			return nil
		},
	}
}

func (a *AppState) createCheckCommand(config *Config, httpCheck *HttpCheck) *cli.Command {
	return &cli.Command{
		Name:  "check",
		Usage: "Perform the HTTP check",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "http-target",
				Aliases:     []string{"u"},
				Usage:       "The http target/upstream in the format http://my-service.com/healthz",
				Required:    true,
				Destination: &httpCheck.upstream,
				Sources:     cli.EnvVars("HTTP_TARGET"),
			},
			&cli.UintFlag{
				Name:        "retries",
				Aliases:     []string{"r"},
				Usage:       "How many times to retry the healthcheck. Pass 0 for infinite tries.",
				Value:       0,
				Destination: &httpCheck.retries,
				Sources:     cli.EnvVars("RETRIES"),
			},
			&cli.UintFlag{
				Name:        "timeout",
				Aliases:     []string{"t"},
				Usage:       "Seconds to wait for each check before considering it a failure.",
				Value:       5,
				Destination: &httpCheck.timeout,
				Sources:     cli.EnvVars("TIMEOUT"),
			},
			&cli.UintFlag{
				Name:        "status-code",
				Aliases:     []string{"c"},
				Usage:       "The status to check for when sending HTTP request",
				Value:       http.StatusOK,
				Destination: &httpCheck.statusCode,
				Sources:     cli.EnvVars("STATUS_CODE"),
			},
			&cli.StringFlag{
				Name:        "log-level",
				Aliases:     []string{"l"},
				Usage:       "The verbosity of the logs (debug, info, warn, error, critical)",
				Value:       "info",
				Destination: &config.logLevel,
				Sources:     cli.EnvVars("LOG_LEVEL"),
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			logLevel := strings.ToLower(config.logLevel)
			switch logLevel {
			case "debug":
			case "info":
			case "warn":
			case "error":
			case "critical":
				a = &AppState{l: newLogger(config)}
			default:
				return fmt.Errorf("unknown log level provided: %s, accepted log levels are debug, info, warn, error, critical", config.logLevel)
			}

			return a.performHttpCheck(httpCheck)
		},
	}
}

func main() {
	ctx := context.Background()

	httpCheck := &HttpCheck{}
	config := &Config{logLevel: "info"}
	logger := newLogger(config)
	app := AppState{l: logger}

	cmd := &cli.Command{
		Name:                  "liveness-check",
		Usage:                 "Perform liveness check on a given URL and exit afterwards, with configurable retries.",
		Suggest:               true,
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			app.createCheckCommand(config, httpCheck),
			createVersionCommand(),
		},
	}

	err := cmd.Run(ctx, os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
