package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
)

// Command line subcommand for "list"
var commandList = cli.Command{
	Name:      "list",
	ShortName: "ls",
	Usage:     "",
	Action:    actionList,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Value: "general.ledger",
			Usage: "The ledger file",
		},
		cli.StringFlag{
			Name:  "project",
			Value: "",
			Usage: "The project name",
		},
		cli.StringFlag{
			Name:  "sort",
			Value: "account",
			Usage: "The sort field",
		},
		cli.BoolFlag{
			Name:  "asc",
			Usage: "The sort direction",
		},
	},
}

// List the ledger contents
func actionList(c *cli.Context) {
	ledgerFile := c.String("file")
	ensureFileExists(ledgerFile)

	f, err := os.Open(ledgerFile)
	check(err)
	defer f.Close()

	l := Ledger{}
	scanner := bufio.NewScanner(f)
	scanner.Split(ScanTransactions)
	for scanner.Scan() {
		text := scanner.Text()

		if strings.Index(text, "#") != 0 {
			t := Transaction{}
			t.FromString(text)
			l.Transactions = append(l.Transactions, t)
		}
	}

	fmt.Print(l.ToString())
}

// ScanTransactions is a split function for a Scanner that returns each line of
// text, stripped of any trailing end-of-line marker. The returned line may be
// empty. The end-of-line marker is one optional carriage return followed
// by one mandatory newline. In regular expression notation, it is `\r\n`.
// The last non-empty line of input will be returned even if it has no newline.
//
// source: https://golang.org/src/bufio/scan.go
func ScanTransactions(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.Index(data, []byte("\n\n")); i >= 0 {
		// We have a double newline terminated line
		return i + 2, dropCR(data[0:i]), nil
	}

	// If we're at EOF, we have a final, non-terminated line. Return it
	if atEOF {
		return len(data), dropCR(data), nil
	}

	// Request more data
	return 0, nil, nil
}

// dropCR drops a terminal \r from the data
// source: https://golang.org/src/bufio/scan.go
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}

	return data
}
