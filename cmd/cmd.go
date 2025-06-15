package cmd

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

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

const (
	CMD_STATIC_CHECK = "static-check"
	CMD_K8S_POD      = "k8s-pod"
	COPYRIGHT        = "(c) Meysam Azad"
)

func createLogger(cfg *config.Config) *logging.Logger {
	logger := logging.NewLogger(logging.WithLogLevel(cfg.GetLogLevel()))
	return &logger
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
			a.createK8sCheckCommand(),
		},
		Flags:          a.createGlobalFlags(),
		Action:         a.rootAction,
		DefaultCommand: CMD_K8S_POD,
		Copyright:      COPYRIGHT,
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
		&cli.IntFlag{
			Name:        "status-code",
			Aliases:     []string{"c"},
			Usage:       "The status to check for when sending HTTP request",
			Value:       http.StatusOK,
			Destination: &a.Config.StatusCode,
			Sources:     cli.EnvVars("STATUS_CODE"),
		},
		&cli.IntFlag{
			Name:        "jitter-min-seconds",
			Usage:       "The min seconds when picking a random time for backoff",
			Value:       5,
			Destination: &a.Config.JitterMin,
			Sources:     cli.EnvVars("JITTER_MIN"),
		},
		&cli.IntFlag{
			Name:        "jitter-max-seconds",
			Usage:       "The max seconds when picking a random time for backoff",
			Value:       10,
			Destination: &a.Config.JitterMax,
			Sources:     cli.EnvVars("JITTER_MAX"),
		},
	}
}

func (a *app) createCheckCommand() *cli.Command {
	return &cli.Command{
		Name:  CMD_STATIC_CHECK,
		Usage: "Perform HTTP check on a static uptream target",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "http-target",
				Aliases:     []string{"u"},
				Usage:       "The http target/upstream in the format http://my-service.com/healthz",
				Required:    true,
				Destination: &a.Config.StaticHTTPTarget.HTTPTarget,
				Sources:     cli.EnvVars("HTTP_TARGET"),
			},
		},
		Action: a.staticHTTPCheck,
	}
}

func (a *app) createK8sCheckCommand() *cli.Command {
	return &cli.Command{
		Name:  CMD_K8S_POD,
		Usage: "Find the most recent pod deployed with a set of label selectors and perform HTTP health check on the given endpoint.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "namespace",
				Aliases:     []string{"n"},
				Destination: &a.Config.K8sPodTarget.Namespace,
				Usage:       "The namespace of the pod. Provide a value for a more efficient search.",
				Required:    false,
				DefaultText: "All namespaces",
				Sources:     cli.EnvVars("NAMESPACE"),
			},
			&cli.StringSliceFlag{
				Name:        "labels",
				Aliases:     []string{"l"},
				Required:    true,
				Usage:       "The key=value pair of label selectors to match the pod against. Specify multiple times or comma-separated.",
				Destination: &a.Config.K8sPodTarget.LabelSelectors,
				Sources:     cli.EnvVars("LABEL_SELECTORS"),
			},
			&cli.StringFlag{
				Name:        "scheme",
				Aliases:     []string{"s"},
				Destination: &a.Config.K8sPodTarget.Scheme,
				Usage:       "The scheme/protocol of the HTTP check (http, https)",
				Value:       "http",
				Sources:     cli.EnvVars("SCHEME"),
			},
			&cli.Int32Flag{
				Name:        "port",
				Aliases:     []string{"p"},
				Destination: &a.Config.K8sPodTarget.Port,
				Usage:       "The port of the pod to send healthcheck request",
				Required:    false,
				DefaultText: "First port of the container",
				Sources:     cli.EnvVars("PORT"),
			},
			&cli.StringFlag{
				Name:        "endpoint",
				Aliases:     []string{"e"},
				Destination: &a.Config.K8sPodTarget.Endpoint,
				Usage:       "The URI of the pod to check against.",
				Value:       "/healthz",
				Sources:     cli.EnvVars("ENDPOINT"),
			},
			&cli.BoolFlag{
				Name:        "tls-verify",
				Value:       false,
				Destination: &a.Config.K8sPodTarget.TLSVerify,
				Usage:       "Whether or not to verify the TLS certificate whene scheme is set to https",
				Sources:     cli.EnvVars("TLS_VERIFY"),
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			k8sPodChecker := &checker.K8sPodChecker{
				Namespace:     a.Config.K8sPodTarget.Namespace,
				Scheme:        a.Config.K8sPodTarget.Scheme,
				LabelSelector: a.Config.K8sPodTarget.LabelSelectors,
				Port:          a.Config.K8sPodTarget.Port,
				Endpoint:      a.Config.K8sPodTarget.Endpoint,
				TLSVerify:     a.Config.K8sPodTarget.TLSVerify,
				Common: &checker.HTTPCommon{
					HTTPClient: &http.Client{
						Timeout: time.Duration(a.Config.Timeout) * time.Second,
					},
					Retries:    a.Config.Retries,
					StatusCode: a.Config.StatusCode,
					Logger:     a.Logger,
					JitterMin:  a.Config.JitterMin,
					JitterMax:  a.Config.JitterMax,
				},
			}

			if err := a.Config.Validate(); err != nil {
				return err
			}

			return k8sPodChecker.Check(ctx)
		},
	}
}

func (a *app) rootAction(ctx context.Context, c *cli.Command) error {
	if err := a.Config.Validate(); err != nil {
		return err
	}

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

func (a *app) staticHTTPCheck(ctx context.Context, c *cli.Command) error {
	httpChecker := &checker.StaticHTTPChecker{
		Upstream: a.Config.StaticHTTPTarget.HTTPTarget,
		Common: &checker.HTTPCommon{
			HTTPClient: &http.Client{
				Timeout: time.Duration(a.Config.Timeout) * time.Second,
			},
			Retries:    a.Config.Retries,
			StatusCode: a.Config.StatusCode,
			Logger:     a.Logger,
			JitterMin:  a.Config.JitterMin,
			JitterMax:  a.Config.JitterMax,
		},
	}

	return httpChecker.Check(ctx)
}
