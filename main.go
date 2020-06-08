package main

import (
	"os"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
