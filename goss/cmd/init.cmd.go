package cmd

import (
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/urfave/cli/v2"
)

var initCmd = &cli.Command{
	Name:  "init",
	Usage: "init a web site directory",
	Action: func(ctx *cli.Context) error {
		dirs := []string{
			"./assets",
			"./assets/css",
			"./assets/images",
			"./assets/scripts",
			"./includes",
		}
		for _, dir := range dirs {
			if !gfile.IsDir(dir) {
				if err := gfile.Mkdir(dir); err != nil {
					return err
				}
			}
		}
		return nil
	},
}
