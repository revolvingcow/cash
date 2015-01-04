package main

import (
	"os"

	"github.com/codegangsta/cli"
)

var commandClear = cli.Command{
	Name:      "clear",
	ShortName: "",
	Usage:     "",
	Action:    actionClear,
}

// Clear current pending transaction
func actionClear(c *cli.Context) {
	file, err := os.OpenFile(PendingFile, os.O_TRUNC|os.O_WRONLY, 0666)
	check(err)
	defer file.Close()

	_, err = file.Write([]byte{})
	check(err)
}
