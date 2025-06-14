package cmd

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/meysam81/x/logging"

	"github.com/meysam81/liveness-check/cmd/checker"
	"github.com/meysam81/liveness-check/cmd/config"
)

type app struct {
	Config *config.Config
	Logger *logging.Logger
}

func NewApp() *app {
	cfg := config.New()
	logger := createLogger(cfg)

	return &app{
		Config: cfg,
		Logger: logger,
	}
}

func createLogger(cfg *config.Config) *logging.Logger {
	logger := logging.NewLogger(logging.WithLogLevel(cfg.GetLogLevel()))
	return &logger
}

func (a *app) updateLogger() {
	a.Logger = createLogger(a.Config)
}

func (a *app) CreateCommand(version, commit, date, builtBy string) *cli.Command {
	cli.VersionPrinter = func(c *cli.Command) {
		fmt.Printf("Version:    %s\n", version)
		fmt.Printf("Commit:     %s\n", commit)
		fmt.Printf("Built:      %s\n", date)
		fmt.Printf("Built by:   %s\n", builtBy)
		fmt.Printf("Go version: %s\n", runtime.Version())
		fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
	}

	return &cli.Command{
		Name:                  "liveness-check",
		Usage:                 "Perform liveness check on a given URL and exit afterwards, with configurable retries.",
		Suggest:               true,
		EnableShellCompletion: true,
		Version:               version,
		Commands: []*cli.Command{
			a.createCheckCommand(),
		},
		Flags:  a.createGlobalFlags(),
		Action: a.rootAction,
	}
}

func (a *app) createGlobalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.UintFlag{
			Name:        "retries",
			Aliases:     []string{"r"},
			Usage:       "How many times to retry the healthcheck. Pass 0 for infinite tries.",
			Value:       0,
			Destination: &a.Config.Retries,
			Sources:     cli.EnvVars("RETRIES"),
		},
		&cli.UintFlag{
			Name:        "timeout",
			Aliases:     []string{"t"},
			Usage:       "Seconds to wait for each check before considering it a failure.",
			Value:       5,
			Destination: &a.Config.Timeout,
			Sources:     cli.EnvVars("TIMEOUT"),
		},
		&cli.UintFlag{
			Name:        "status-code",
			Aliases:     []string{"c"},
			Usage:       "The status to check for when sending HTTP request",
			Value:       http.StatusOK,
			Destination: &a.Config.StatusCode,
			Sources:     cli.EnvVars("STATUS_CODE"),
		},
		&cli.StringFlag{
			Name:        "log-level",
			Aliases:     []string{"l"},
			Usage:       "The verbosity of the logs (debug, info, warn, error, critical)",
			Value:       "info",
			Destination: &a.Config.LogLevel,
			Sources:     cli.EnvVars("LOG_LEVEL"),
		},
	}
}

func (a *app) createCheckCommand() *cli.Command {
	return &cli.Command{
		Name:  "check",
		Usage: "Perform the HTTP check",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "http-target",
				Aliases:     []string{"u"},
				Usage:       "The http target/upstream in the format http://my-service.com/healthz",
				Required:    true,
				Destination: &a.Config.HTTPTarget,
				Sources:     cli.EnvVars("HTTP_TARGET"),
			},
		},
		Action: a.checkAction,
	}
}

func (a *app) rootAction(ctx context.Context, c *cli.Command) error {
	if err := a.Config.Validate(); err != nil {
		return err
	}

	a.updateLogger()

	if c.NArg() == 0 {
		var availableCommands []string
		for _, subcommand := range c.Commands {
			availableCommands = append(availableCommands, subcommand.Name)
		}
		return fmt.Errorf("no command provided. available subcommands: %s",
			strings.Join(availableCommands, ", "))
	}

	return nil
}

func (a *app) checkAction(ctx context.Context, c *cli.Command) error {
	httpChecker := checker.NewHTTPChecker(
		a.Config.HTTPTarget,
		a.Config.Timeout,
		a.Config.Retries,
		a.Config.StatusCode,
		a.Logger,
	)

	return httpChecker.Check(ctx)
}
