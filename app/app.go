package app

import (
	_ "embed"
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	csvtools "github.com/palsivertsen/csv-tools"
	"github.com/palsivertsen/csv-tools/csv/col"
	"github.com/urfave/cli/v2"
)

func App() *cli.App {
	return &cli.App{
		Usage:                "manipulate csv streams",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "column",
				Usage: "edit columns",
				Subcommands: []*cli.Command{
					removeColumnCommand(),
				},
			},
			printCommand(),
			completionHelper(),
		},
	}
}

func removeColumnCommand() *cli.Command {
	return &cli.Command{
		Name:  "remove",
		Usage: "remove columns",
		Flags: []cli.Flag{
			&cli.IntSliceFlag{
				Name: "columns",
			},
		},
		Action: func(ctx *cli.Context) error {
			reader := col.NewRemoverReader(
				csv.NewReader(ctx.App.Reader),
				ctx.IntSlice("columns")...,
			)

			writer := csv.NewWriter(ctx.App.Writer)
			defer writer.Flush()

			if rowNumber, err := csvtools.Copy(writer, reader); err != nil {
				return csvtools.RowError{
					RowNumber: rowNumber,
					Err:       err,
				}
			}

			return nil
		},
	}
}

func printCommand() *cli.Command {
	return &cli.Command{
		Name: "print",
		Subcommands: []*cli.Command{
			{
				Name:  "header",
				Usage: "print header",
				Action: func(ctx *cli.Context) error {
					reader := csv.NewReader(ctx.App.Reader)

					row, err := reader.Read()
					if err != nil {
						if errors.Is(err, io.EOF) {
							return nil
						}
						return fmt.Errorf("read first row: %w", err)
					}

					for i, col := range row {
						if _, err := fmt.Fprintf(ctx.App.Writer, "%d: %s\n", i, col); err != nil {
							return fmt.Errorf("write header: %w", err)
						}
					}

					return nil
				},
			},
		},
	}
}

//go:generate curl -O https://raw.githubusercontent.com/urfave/cli/v2.23.7/autocomplete/bash_autocomplete
//go:embed bash_autocomplete
var script string

func completionHelper() *cli.Command {
	return &cli.Command{
		Name:  "completion",
		Usage: "shell completion scripts",
		Description: `For bash run:
	PROG=csv eval "$(csv completion --shell=bash)"`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "shell",
				Required: true,
				Action: func(ctx *cli.Context, shell string) error {
					switch shell {
					case "bash":
						return nil
					}

					return fmt.Errorf("shell %q not supported", shell)
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			if _, err := fmt.Fprint(ctx.App.Writer, script); err != nil {
				return fmt.Errorf("print script: %w", err)
			}

			return nil
		},
	}
}
