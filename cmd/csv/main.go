package main

import (
	"fmt"
	"log"
	"os"

	"github.com/palsivertsen/csv-tools/app"
)

func main() {
	if err := run(); err != nil {
		log.Printf("run: %s", err.Error())
	}
}

func run() error {
	app := app.App()

	if err := app.Run(os.Args); err != nil {
		return fmt.Errorf("app run: %w", err)
	}

	return nil
}
