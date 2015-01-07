package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

// Command line subcommand for "debit"
var commandDebit = cli.Command{
	Name:      "debit",
	ShortName: "dr",
	Usage:     "",
	Action:    actionDebit,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Value: "general.ledger",
			Usage: "",
		},
	},
}

// Add a debit to the pending transaction
func actionDebit(c *cli.Context) {
	pendingFile := fmt.Sprintf(".%s", c.String("file"))
	ensureFileExists(pendingFile)

	args := c.Args()
	li := len(args) - 1
	name, err := parseAccount(args[:li])
	check(err)
	amount, err := parseAmount(args[li])
	check(err)

	a := Account{
		Name:   name,
		Debit:  true,
		Amount: amount,
	}

	f, err := os.OpenFile(pendingFile, os.O_APPEND|os.O_WRONLY, 0666)
	check(err)
	defer f.Close()
	_, err = f.WriteString(a.ToString())
	check(err)
}
