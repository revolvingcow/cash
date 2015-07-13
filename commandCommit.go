package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/cli"
)

// Command line subcommand for "commit"
var commandCommit = cli.Command{
	Name:      "commit",
	ShortName: "c",
	Usage:     "",
	Action:    actionCommit,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "date",
			Value: time.Now().UTC().Format("2006-01-02"),
			Usage: "The transaction date",
		},
		cli.StringFlag{
			Name:  "file",
			Value: "general.ledger",
			Usage: "The ledger file to store the transaction",
		},
	},
}

// Commit the pending transaction
func actionCommit(c *cli.Context) {
	ledgerFile := c.String("file")
	ensureFileExists(ledgerFile)

	pendingFile := fmt.Sprintf(".%s", c.String("file"))
	ensureFileExists(pendingFile)

	date, err := parseDate(c.String("date"))
	check(err)

	args := c.Args()
	project := parseProject(args)
	description := parseDescription(args, project)

	writeTransaction(ledgerFile, pendingFile, project, description, date)
	actionClear(c)
}

// Parse a given string to extract a project name
func parseProject(fields []string) string {
	project := "@general"

	for i := 0; i < len(fields); i++ {
		if strings.HasPrefix(fields[i], "@") {
			project = fields[i]
			break
		}
	}

	return project
}

// Parse the description from the arguments
func parseDescription(fields []string, project string) string {
	for i := 0; i < len(fields); i++ {
		if fields[i] == project {
			fields[i] = ""
			break
		}
	}

	return strings.Replace(strings.Join(fields, " "), "  ", " ", -1)
}

// Write a transaction line where there is a pending transaction
func writeTransaction(ledgerFile, pendingFile, project, description string, date time.Time) {
	if !hasPendingTransaction(pendingFile) {
		check(errors.New("No pending transaction to write"))
	}

	pending, err := ioutil.ReadFile(pendingFile)
	check(err)

	t := Transaction{
		Date:        date,
		Project:     project,
		Description: description,
		Accounts:    []Account{},
	}

	lines := strings.Split(strings.TrimRight(string(pending), "\n"), "\n")
	for _, line := range lines {
		var a Account
		err = a.FromString(line)
		check(err)

		t.Accounts = append(t.Accounts, a)
	}

	err = t.CheckBalance()
	check(err)

	file, err := os.OpenFile(ledgerFile, os.O_APPEND|os.O_WRONLY, 0666)
	check(err)
	defer file.Close()
	_, err = file.WriteString(t.ToString())
	check(err)
}
