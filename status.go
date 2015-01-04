package main

import (
	"fmt"
	"io/ioutil"

	"github.com/codegangsta/cli"
)

var commandStatus = cli.Command{
	Name:      "status",
	ShortName: "stat",
	Usage:     "",
	Action:    actionStatus,
}

// Display the current status of the ledger
func actionStatus(c *cli.Context) {
	if hasPendingTransaction() {
		pending, err := ioutil.ReadFile(PendingFile)
		check(err)
		fmt.Println(string(pending))
	} else {
		fmt.Println("No pending transactions")
	}
}
