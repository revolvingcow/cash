package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
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

	writeTransaction(date.Format("2006-01-02"), project, description)
}

// Parse the given string to extract a proper date
func parseDate(in string) (time.Time, error) {
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
			return d, nil
		}
	}

	return time.Now().UTC(), errors.New("No valid date provided")
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

type Account struct {
	Name   string
	Debit  bool
	Amount *big.Rat
}

type Transaction struct {
	Date        time.Time
	Project     string
	Description string
	Accounts    []Account
}

func (t *Transaction) FromString(text string) error {
	// Parse the lines of text
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		fields := strings.Split(line, "\t")

		switch i {
		case 0:
			date, err := parseDate(fields[0])
			check(err)
			project := fields[1]
			description := ""
			if len(fields) > 2 {
				description = strings.Join(fields[2:], " ")
			}

			t = &Transaction{
				Date:        date,
				Project:     project,
				Description: description,
				Accounts:    []Account{},
			}
			break

		default:
			if len(fields) != 3 {
				break
			}

			account := fields[1]
			debit := true

			if strings.HasPrefix(fields[2], "-") {
				debit = false
			}
			value := new(big.Rat)
			value.SetString(fields[2][1:])

			t.Accounts = append(
				t.Accounts,
				Account{
					Name:   account,
					Debit:  debit,
					Amount: value,
				})

			break
		}
	}

	if len(t.Accounts) == 0 {
		return errors.New("Transaction does not have any accounts")
	}

	// Check that they balance
	balance := new(big.Rat)
	for _, a := range t.Accounts {
		if a.Debit {
			balance.Add(balance, a.Amount)
		} else {
			balance.Sub(balance, a.Amount)
		}
	}
	if balance.FloatString(2) != "0.00" {
		return errors.New("Transaction does not balance")
	}

	return nil
}

// Write a transaction line where there is a pending transaction
func writeTransaction(date, project, description string) {
	if !hasPendingTransaction() {
		check(errors.New("No pending transaction to write"))
	}

	pending, err := ioutil.ReadFile(PendingFile)
	check(err)

	// Find the line containing @pending and replace it with our transaction
	s := fmt.Sprintf("%s\t%s\t%s\n%s\n", date, project, description, string(pending))
	var t Transaction
	err = t.FromString(s)
	check(err)

	file, err := os.OpenFile(Ledger, os.O_APPEND|os.O_WRONLY, 0666)
	check(err)
	defer file.Close()
	_, err = file.WriteString(s)
	check(err)
}
