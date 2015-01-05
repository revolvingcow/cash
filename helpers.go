package main

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

// Helper function to check for fatal errors
func check(e error) {
	if e != nil {
		log.Fatal(fmt.Sprintf("Error: %s", e))
	}
}

func ensureFileExists(fileName string) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		_, err = os.Create(fileName)
		check(err)
	}
}

// Format the ledger so it is human readable
func formatLedger() {
}

// Determines if there is currently a pending transaction in the ledger
func hasPendingTransaction(pendingFile string) bool {
	file, err := os.Open(pendingFile)
	check(err)
	defer file.Close()

	info, err := file.Stat()
	check(err)

	return info.Size() > 0
}

// Parse a given string to extract an account name
func parseAccount(fields []string) (string, error) {
	if len(fields) > 0 {
		return strings.Join(fields, " "), nil
	}

	return "", errors.New("No account information found")
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

// Parse the value from the arguments
func parseAmount(text string) (*big.Rat, error) {
	amount := new(big.Rat)

	if r, _ := amount.SetString(text); r != nil {
		return amount, nil
	}

	return nil, errors.New("No value found")
}
