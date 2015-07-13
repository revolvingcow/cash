package main

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

type Account struct {
	Name   string
	Debit  bool
	Amount *big.Rat
}

// Convert from a string to an account
func (a *Account) FromString(text string) error {
	parts := strings.Split(text, "\t")
	fields := []string{}
	for _, p := range parts {
		if p != "" {
			fields = append(fields, p)
		}
	}

	if len(fields) != 2 {
		return errors.New("Invalid account format")
	}

	debit := true
	if strings.HasPrefix(fields[1], "-") {
		debit = false
	}

	a.Debit = debit
	a.Name = fields[0]
	a.Amount = new(big.Rat)
	a.Amount.SetString(fields[1][1:])

	return nil
}

// Convert the account to string format
func (a *Account) ToString() string {
	symbol := "-"
	if a.Debit {
		symbol = "+"
	}

	return fmt.Sprintf("\t%s\t%s%s\n", a.Name, symbol, a.Amount.FloatString(2))
}
