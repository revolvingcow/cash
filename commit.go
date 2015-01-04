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

var commandCommit = cli.Command{
	Name:      "commit",
	ShortName: "c",
	Usage:     "",
	Action:    actionCommit,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "date",
			Value: time.Now().UTC().Format("2006-01-02"),
			Usage: "",
		},
	},
}

// Commit the pending transaction
func actionCommit(c *cli.Context) {
	date, err := parseDate(c.String("date"))
	check(err)

	args := c.Args()
	project := parseProject(args)
	description := parseDescription(args, project)

	writeTransaction(date, project, description)
}

// Parse the given string to extract a proper date
func parseDate(in string) (string, error) {
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"2006-1-2",
		"2006/1/2",
		"01-02-2006",
		"01/02/2006",
		"1-2-2006",
		"1/2/2006",
		"Jan 2, 2006",
		"Jan 02, 2006",
		"2 Jan 2006",
		"02 Jan 2006",
	}

	for _, f := range formats {
		d, err := time.Parse(f, in)
		if err == nil {
			return d.Format(formats[0]), nil
		}
	}

	return "", errors.New("No valid date provided")
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
func writeTransaction(date, project, description string) {
	if !hasPendingTransaction() {
		check(errors.New("No pending transaction to write"))
	}

	pending, err := ioutil.ReadFile(PendingFile)
	check(err)

	file, err := os.OpenFile(Ledger, os.O_APPEND|os.O_WRONLY, 0666)
	check(err)
	defer file.Close()

	// Find the line containing @pending and replace it with our transaction
	transaction := fmt.Sprintf("%s\t%s\t%s\n%s\n", date, project, description, pending)

	_, err = file.WriteString(transaction)
	check(err)
}
