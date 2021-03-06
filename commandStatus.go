package main

import (
	"fmt"
	"io/ioutil"

	"github.com/codegangsta/cli"
)

// Command line subcommand for "status"
var commandStatus = cli.Command{
	Name:      "status",
	ShortName: "stat",
	Usage:     "",
	Action:    actionStatus,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Value: "general.ledger",
			Usage: "The ledger file",
		},
	},
}

// Display the current status of the ledger
func actionStatus(c *cli.Context) {
	pendingFile := fmt.Sprintf(".%s", c.String("file"))
	ensureFileExists(pendingFile)

	if hasPendingTransaction(pendingFile) {
		pending, err := ioutil.ReadFile(pendingFile)
		check(err)
		fmt.Println(string(pending))
	} else {
		fmt.Println("No pending transactions")
	}
}
