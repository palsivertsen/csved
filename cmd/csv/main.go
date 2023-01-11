package main

import (
	"fmt"
	"log"
	"os"

	"csv-tools/csv/col"

	"github.com/urfave/cli/v2"
)

func main() {
	if err := run(); err != nil {
		log.Printf("run: %s", err.Error())
	}
}

func run() error {
	app := cli.App{
		Name:                 "csv",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			col.Command(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return nil
}
