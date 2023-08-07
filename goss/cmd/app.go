package cmd

import (
	"os"

	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:      "goss",
	Usage:     "static web server of golang (include and build html)",
	UsageText: "goss <command> [options]",
	Commands: []*cli.Command{
		initCmd,
		runCmd,
		buildCmd,
	},
	HideHelpCommand: true,
}

func Execute() error {
	return app.Run(os.Args)
}
