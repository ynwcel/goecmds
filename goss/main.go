package main

import (
	"github.com/ynwcel/goecmds/goss/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
