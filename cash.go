package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/cli"
)

const (
	// Application name
	APP_NAME = "cash"

	// Application usage description
	APP_USAGE = "counting coins"

	// Application version
	APP_VER = "0.0.0"
)

var (
	// Ledger file name
	Ledger            = "example.ledger"
	TransactionFormat = "%s\t%s\t%s"
	AccountFormat     = "\t%s\t%s"
)

// Application entry point
func main() {
	// Flags pertaining to a transaction action
	transactionFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "date",
			Value: time.Now().UTC().Format("2006-01-02"),
			Usage: "",
		},
	}

	app := cli.NewApp()
	app.Name = APP_NAME
	app.Usage = APP_USAGE
	app.Version = APP_VER
	app.Commands = []cli.Command{
		{
			Name:      "credit",
			ShortName: "cr",
			Usage:     "",
			Action:    actionCredit,
		},
		{
			Name:      "debit",
			ShortName: "dr",
			Usage:     "",
			Action:    actionDebit,
		},
		{
			Name:      "status",
			ShortName: "stat",
			Usage:     "",
			Action:    actionStatus,
		},
		{
			Name:      "commit",
			ShortName: "c",
			Usage:     "",
			Action:    actionCommit,
			Flags:     transactionFlags,
		},
		{
			Name:      "list",
			ShortName: "ls",
			Usage:     "",
			Action:    actionList,
		},
	}

	app.Run(os.Args)
}

// Helper function to check for fatal errors
func check(e error) {
	if e != nil {
		log.Fatal(fmt.Sprintf("Error: %s", e))
	}
}

// Add a credit to the pending transaction
func actionCredit(c *cli.Context) {
	addPendingTransaction()

	// Format: /t{account}/t-{value}
}

// Add a debit to the pending transaction
func actionDebit(c *cli.Context) {
	addPendingTransaction()

	// Format: /t{account}/t+{value}
}

// Display the current status of the ledger
func actionStatus(c *cli.Context) {
	log.Println(hasPendingTransaction())
}

// Commit the pending transaction
func actionCommit(c *cli.Context) {
	date := parseDate(c.String("date"))
	if date == "" {
		log.Fatal("Invalid transaction date")
	}

	args := c.Args()
	project := parseProject(args)
	description := parseDescription(args, project)

	err := writeTransaction(date, project, description)
	if err != nil {
		log.Fatal(err)
	}
}

// List the ledger contents
func actionList(c *cli.Context) {
}

// Format the ledger so it is human readable
func formatLedger() {
}

// Determines if there is currently a pending transaction in the ledger
func hasPendingTransaction() bool {
	file, err := os.Open(Ledger)
	check(err)
	defer file.Close()

	info, err := file.Stat()
	check(err)
	size := info.Size()
	if size > 1024 {
		size = 1024
	}

	_, err = file.Seek(size*-1, 2)
	check(err)

	buffer := make([]byte, size)
	_, err = file.Read(buffer)
	check(err)

	return strings.Contains(string(buffer), "@pending")
}

// Adds a pending transaction if one is not already present
func addPendingTransaction() {
	if !hasPendingTransaction() {
		file, err := os.OpenFile(Ledger, os.O_APPEND|os.O_WRONLY, 0666)
		check(err)
		defer file.Close()

		_, err = file.WriteString("@pending")
		check(err)
	}
}

// Parse the given string to extract a proper date
func parseDate(in string) string {
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
			return d.Format(formats[0])
		}
	}

	return ""
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
func writeTransaction(date, project, description string) error {
	if !hasPendingTransaction() {
		return errors.New("No pending transaction to write")
	}

	file, err := os.OpenFile(Ledger, os.O_RDWR, 0666)
	check(err)
	defer file.Close()

	info, err := file.Stat()
	check(err)
	size := info.Size()
	if size > 1024 {
		size = 1024
	}

	_, err = file.Seek(size*-1, 2)
	check(err)

	buffer := make([]byte, size)
	_, err = file.Read(buffer)
	check(err)

	// Find the line containing @pending and replace it with our transaction
	line := fmt.Sprintf("%s\t%s\t%s", date, project, description)
	buffer = bytes.Replace(buffer, []byte("@pending"), []byte(line), -1)

	offset := info.Size() - size
	_, err = file.WriteAt(buffer, offset)
	check(err)

	return nil
}
