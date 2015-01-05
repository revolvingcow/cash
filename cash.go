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
	LedgerFile  = "general.ledger"
	PendingFile = filepath.Join(os.TempDir(), "pending.ledger")
)

// Initialize the application
func init() {
	if _, err := os.Stat(LedgerFile); os.IsNotExist(err) {
		_, err = os.Create(LedgerFile)
		check(err)
		log.Println("Created ledger file at", LedgerFile)
	}

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
