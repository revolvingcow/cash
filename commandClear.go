package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

// Command line subcommand for "clear"
var commandClear = cli.Command{
	Name:      "clear",
	ShortName: "",
	Usage:     "",
	Action:    actionClear,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Value: "general.ledger",
			Usage: "",
		},
	},
}

// Clear current pending transaction
func actionClear(c *cli.Context) {
	pendingFile := fmt.Sprintf(".%s", c.String("file"))
	ensureFileExists(pendingFile)

	file, err := os.OpenFile(pendingFile, os.O_TRUNC|os.O_WRONLY, 0666)
	check(err)
	defer file.Close()
	_, err = file.Write([]byte{})
	check(err)
}
