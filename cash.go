package main

import (
	"os"

	"github.com/codegangsta/cli"
)

const (
	APP_NAME  = "cash"
	APP_USAGE = "counting coins"
	APP_VER   = "0.0.0"
)

// Application entry point
func main() {
	app := cli.NewApp()
	app.Name = APP_NAME
	app.Usage = APP_USAGE
	app.Version = APP_VER
	app.Commands = []cli.Command{
		commandCredit,
		commandDebit,
		commandStatus,
		commandCommit,
		commandList,
		commandClear,
	}

	app.Run(os.Args)
}
