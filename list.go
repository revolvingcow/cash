package main

import "github.com/codegangsta/cli"

var commandList = cli.Command{
	Name:      "list",
	ShortName: "ls",
	Usage:     "",
	Action:    actionList,
}

// List the ledger contents
func actionList(c *cli.Context) {
}
