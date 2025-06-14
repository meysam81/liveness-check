package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
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
	c *Config
	h *HttpCheck
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

func newLogger(level string) *logging.Logger {
	logger := logging.NewLogger(logging.WithLogLevel(translateLogLevel(level)))
	return &logger
}

func (a *AppState) performHttpCheck(ctx context.Context) error {
	client := &http.Client{
		Timeout: time.Duration(a.h.timeout) * time.Second,
	}

	var tries uint = 0
	for {
		if a.h.retries > 0 && tries >= a.h.retries {
			a.l.Error().Msgf("halting check with max retries reached: %d", a.h.retries)
			return nil
		}

		start := time.Now()

		resp, err := client.Get(a.h.upstream)
		if err != nil {
			tries += 1
			jitter := randomBetween(5, 10)

			a.l.Info().Err(err).Msgf("[%d] upstream check unsuccessful, retrying in %ds", tries, jitter)

			t := time.NewTicker(time.Duration(jitter) * time.Second)
			defer t.Stop()

			select {
			case <-ctx.Done():
				a.l.Info().Msg("shutdown signal received. stopping...")
				return nil
			case <-t.C:
				continue
			}
		}

		elapsed := time.Since(start)

		if resp.StatusCode == int(a.h.statusCode) {
			a.l.Info().Msgf("took %s for http check successful with status code: %s", elapsed.Round(time.Millisecond).String(), resp.Status)
			return nil
		}
	}
}

func (a *AppState) createCheckCommand() *cli.Command {
	return &cli.Command{
		Name:  "check",
		Usage: "Perform the HTTP check",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "http-target",
				Aliases:     []string{"u"},
				Usage:       "The http target/upstream in the format http://my-service.com/healthz",
				Required:    true,
				Destination: &a.h.upstream,
				Sources:     cli.EnvVars("HTTP_TARGET"),
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			return a.performHttpCheck(ctx)
		},
	}
}

func main() {
	ctx := context.Background()

	config := &Config{logLevel: "info"}
	logger := newLogger(config.logLevel)
	app := AppState{l: logger, c: config, h: &HttpCheck{}}

	cli.VersionPrinter = func(c *cli.Command) {
		fmt.Printf("Version:    %s\n", version)
		fmt.Printf("Commit:     %s\n", commit)
		fmt.Printf("Built:      %s\n", date)
		fmt.Printf("Built by:   %s\n", builtBy)
		fmt.Printf("Go version: %s\n", runtime.Version())
		fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
	}

	cmd := &cli.Command{
		Name:                  "liveness-check",
		Usage:                 "Perform liveness check on a given URL and exit afterwards, with configurable retries.",
		Suggest:               true,
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			app.createCheckCommand(),
		},
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:        "retries",
				Aliases:     []string{"r"},
				Usage:       "How many times to retry the healthcheck. Pass 0 for infinite tries.",
				Value:       0,
				Destination: &app.h.retries,
				Sources:     cli.EnvVars("RETRIES"),
			},
			&cli.UintFlag{
				Name:        "timeout",
				Aliases:     []string{"t"},
				Usage:       "Seconds to wait for each check before considering it a failure.",
				Value:       5,
				Destination: &app.h.timeout,
				Sources:     cli.EnvVars("TIMEOUT"),
			},
			&cli.UintFlag{
				Name:        "status-code",
				Aliases:     []string{"c"},
				Usage:       "The status to check for when sending HTTP request",
				Value:       http.StatusOK,
				Destination: &app.h.statusCode,
				Sources:     cli.EnvVars("STATUS_CODE"),
			},
			&cli.StringFlag{
				Name:        "log-level",
				Aliases:     []string{"l"},
				Usage:       "The verbosity of the logs (debug, info, warn, error, critical)",
				Value:       "info",
				Destination: &app.c.logLevel,
				Sources:     cli.EnvVars("LOG_LEVEL"),
			},
		},
		Version: version,
		Action: func(ctx context.Context, c *cli.Command) error {
			logLevel := strings.ToLower(app.c.logLevel)
			switch logLevel {
			case "debug":
			case "info":
			case "warn":
			case "error":
			case "critical":
				app.l = newLogger(app.c.logLevel)
			default:
				return fmt.Errorf("unknown log level provided: %s, accepted log levels are debug, info, warn, error, critical", app.c.logLevel)
			}

			if c.NArg() == 0 {
				availableCommands := []string{}
				for _, subcommand := range c.Commands {
					availableCommands = append(availableCommands, subcommand.Name)
				}
				return fmt.Errorf("no command provided. available subcommands: %s", strings.Join(availableCommands, ", "))
			}

			return nil
		},
	}

	ctxT, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGHUP)
	defer stop()

	app.l.Debug().Msg("Starting the app...")

	err := cmd.Run(ctxT, os.Args)
	if err != nil {
		fmt.Println(err)
	}

	app.l.Info().Msg("it's been a pleasure. goodbye till next time.")
}
