package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

// Command line subcommand for "credit"
var commandCredit = cli.Command{
	Name:      "credit",
	ShortName: "cr",
	Usage:     "",
	Action:    actionCredit,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Value: "general.ledger",
			Usage: "",
		},
	},
}

// Add a credit to the pending transaction
func actionCredit(c *cli.Context) {
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
		Debit:  false,
		Amount: amount,
	}

	f, err := os.OpenFile(pendingFile, os.O_APPEND|os.O_WRONLY, 0666)
	check(err)
	defer f.Close()
	_, err = f.WriteString(a.ToString())
	check(err)
}
