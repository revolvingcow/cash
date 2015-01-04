package main

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
)

// Helper function to check for fatal errors
func check(e error) {
	if e != nil {
		log.Fatal(fmt.Sprintf("Error: %s", e))
	}
}

// Format the ledger so it is human readable
func formatLedger() {
}

// Determines if there is currently a pending transaction in the ledger
func hasPendingTransaction() bool {
	file, err := os.Open(PendingFile)
	check(err)
	defer file.Close()

	info, err := file.Stat()
	check(err)

	return info.Size() > 0
}

// Parse a given string to extract an account name
func parseAccount(fields []string) (string, error) {
	for i := 0; i < len(fields); i++ {
		if strings.HasPrefix(fields[i], "#") {
			return fields[i], nil
		}
	}

	return fields[0], nil
}

// Parse the value from the arguments
func parseValue(fields []string, account string) (*big.Rat, error) {
	r := new(big.Rat)

	for i := 0; i < len(fields); i++ {
		if fields[i] != account {
			r.SetString(fields[i])
			return r, nil
		}
	}

	return r, errors.New("No value found")
}
