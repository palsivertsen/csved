package app

import (
	_ "embed"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

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
					pickColumnCommand(),
				},
			},
			printCommand(),
			completionHelper(),
			filterCommand(),
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

func pickColumnCommand() *cli.Command {
	return &cli.Command{
		Name:  "pick",
		Usage: "pick columns",
		Flags: []cli.Flag{
			&cli.IntSliceFlag{
				Name: "columns",
			},
		},
		Action: func(ctx *cli.Context) error {
			reader := col.NewPickReader(
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

const flagColumnRegexP = "column-regexp"

func filterCommand() *cli.Command {
	return &cli.Command{
		Name:  "filter",
		Usage: "filter row by column(s)",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  flagColumnRegexP,
				Usage: "column-index:regex-pattern",
			},
		},
		Action: func(ctx *cli.Context) error {
			rules := ctx.StringSlice(flagColumnRegexP)
			matcher := make(col.AllRowMatcher, 0, len(rules))

			for i, cr := range ctx.StringSlice(flagColumnRegexP) {
				sp := strings.SplitN(cr, ":", 2)
				if len(sp) != 2 {
					return ParamError{Name: flagColumnRegexP}
				}

				colIndexString, regexPattern := sp[0], sp[1]
				colIndex, err := strconv.Atoi(colIndexString)
				if err != nil {
					return ParamError{Name: flagColumnRegexP, Err: fmt.Errorf("param #%d: parse column index %q: %w", i, colIndexString, err)}
				}

				comp, err := regexp.Compile(regexPattern)
				if err != nil {
					return ParamError{Name: flagColumnRegexP, Err: fmt.Errorf("param #%d: parse pattern: %w", i, err)}
				}

				matcher = append(
					matcher,
					&col.IndexMatcher{
						Index: colIndex,
						Matcher: col.StringMatcherFunc(func(s string) (bool, error) {
							return comp.MatchString(s), nil
						}),
					},
				)
			}

			reader := col.FilterReader{
				Reader:  csv.NewReader(ctx.App.Reader),
				Matcher: matcher,
			}
			writer := csv.NewWriter(ctx.App.Writer)
			defer writer.Flush()

			if _, err := csvtools.Copy(writer, &reader); err != nil {
				return fmt.Errorf("copy: %w", err)
			}

			return nil
		},
	}
}

type ParamError struct {
	Name string
	Err  error
}

func (e ParamError) Error() string { return fmt.Sprintf("malformed parameter %q", e.Name) }

func (e ParamError) Unwrap() error { return e.Err }
