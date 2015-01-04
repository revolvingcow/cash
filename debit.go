package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

var commandDebit = cli.Command{
	Name:      "debit",
	ShortName: "dr",
	Usage:     "",
	Action:    actionDebit,
}

// Add a debit to the pending transaction
func actionDebit(c *cli.Context) {
	args := c.Args()

	account, err := parseAccount(args)
	check(err)

	value, err := parseValue(args, account)
	check(err)

	f, err := os.OpenFile(PendingFile, os.O_APPEND|os.O_WRONLY, 0666)
	check(err)
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("\t%s\t-%s\n", account, value.FloatString(2)))
	check(err)
}
