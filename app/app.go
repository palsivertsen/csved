package app

import (
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
		Name:                 "csv",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name: "column",
				Subcommands: []*cli.Command{
					removeColumnCommand(),
					printHeaderCommand(),
				},
			},
		},
	}
}

func removeColumnCommand() *cli.Command {
	return &cli.Command{
		Name: "remove",
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

func printHeaderCommand() *cli.Command {
	return &cli.Command{
		Name: "header",
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
	}
}
