package col

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
)

func RM(ctx *cli.Context) error {
	in := os.Stdin
	out := os.Stdout
	cols := ctx.IntSlice("columns")
	slices.SortFunc(cols, func(a, b int) bool { return a > b })

	r := &columnSkipper{
		r:               csv.NewReader(in),
		descendingSkips: cols,
	}

	w := csv.NewWriter(out)
	defer w.Flush()

	for {
		row, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("read row: %w", err)
		}

		if err := w.Write(row); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}

	return nil
}

type columnSkipper struct {
	r               *csv.Reader
	descendingSkips []int
}

func (c *columnSkipper) Read() (record []string, err error) {
	row, err := c.r.Read()
	if err != nil {
		return nil, fmt.Errorf("skip column: %w", err)
	}

	for _, i := range c.descendingSkips {
		row = slices.Delete(row, i, i+1)
	}

	return row, nil
}
