package main

import (
	"os"

	"github.com/alexjx/sd-companion/cmd"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v3"
)

var (
	// Version is the version of the application
	Version = "0.0.1"
	// BuildNumber is the build number of the application
	BuildNumber = ""
	// GitCommit is the commit hash of the application
	GitCommit = ""
	// BuildTime is the date of the build
	BuildTime = ""
)

func main() {
	app := cli.App{
		Name:                  "sd-companion",
		Usage:                 "sd-companion is a tool to broswer images",
		EnableShellCompletion: true,
		ExitErrHandler: func(cctx *cli.Context, err error) {
			cli.HandleExitCoder(err)
		},
		Version: Version,
		Commands: []*cli.Command{
			cmd.ServeCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatalf("Error: %v", err)
	}
}
