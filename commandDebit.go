package main

import (
	"os"
	"path/filepath"

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
	pendingFile := filepath.Join(os.TempDir(), c.String("file"))
	ensureFileExists(pendingFile)

	args := c.Args()
	name, err := parseAccount(args)
	check(err)
	amount, err := parseValue(args, name)
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
