package main

import (
	"github.com/ynwcel/gcmds/goss/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
