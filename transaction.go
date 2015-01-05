package main

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
)

type Transaction struct {
	Date        time.Time
	Project     string
	Description string
	Accounts    []Account
}

// Create a transaction from a string
func (t *Transaction) FromString(text string) {
	// Parse the lines of text
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}

		switch i {
		case 0:
			fields := strings.Split(line, "\t")

			date, err := parseDate(fields[0])
			check(err)

			project := fields[1]

			description := ""
			if len(fields) > 2 {
				description = strings.Join(fields[2:], " ")
			}

			t.Date = date
			t.Project = project
			t.Description = description
			t.Accounts = []Account{}
			break

		default:
			var a Account
			err := a.FromString(line)
			check(err)
			t.Accounts = append(t.Accounts, a)
			break
		}
	}
}

// Check the transaction to ensure it is balanced
func (t *Transaction) CheckBalance() error {
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

// Convert a transaction to string format
func (t *Transaction) ToString() string {
	accounts := ""

	for _, account := range t.Accounts {
		accounts += account.ToString()
	}

	err := t.CheckBalance()
	if err != nil {
		accounts += fmt.Sprintf("Error: %s\n", err)
	}

	return fmt.Sprintf("%s\t%s\t%s\n%s\n", t.Date.Format("2006-01-02"), t.Project, t.Description, accounts)
}
