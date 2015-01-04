package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
)

const (
	APP_NAME  = "cash"
	APP_USAGE = "counting coins"
	APP_VER   = "0.0.0"
)

var (
	Ledger            = "general.ledger"
	TransactionFormat = "%s\t%s\t%s"
	AccountFormat     = "\t%s\t%s"
	PendingFile       = filepath.Join(os.TempDir(), "pending.ledger")
)

// Initialize the application
func init() {
	if _, err := os.Stat(PendingFile); os.IsNotExist(err) {
		_, err = os.Create(PendingFile)
		check(err)
		log.Println("Created pending file at", PendingFile)
	}
}

// Application entry point
func main() {
	app := cli.NewApp()
	app.Name = APP_NAME
	app.Usage = APP_USAGE
	app.Version = APP_VER
	app.Commands = []cli.Command{
		commandCredit,
		commandDebit,
		commandStatus,
		commandCommit,
		commandList,
		commandClear,
	}

	app.Run(os.Args)
}
