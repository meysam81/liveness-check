package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/meysam81/liveness-check/cmd"
)

var (
	version = "dev"
	commit  = "HEAD"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGHUP)
	defer stop()

	app := cmd.NewApp()
	command := app.CreateCommand(version, commit, date, builtBy)

	app.Logger.Debug().Msg("Starting the app...")

	if err := command.Run(ctx, os.Args); err != nil {
		app.Logger.Fatal().Err(err).Send()
	}
}
