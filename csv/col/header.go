package col

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
)

func CMDPrintHeader(ctx *cli.Context) error {
	in := os.Stdin
	out := os.Stdout
	r := csv.NewReader(in)

	row, err := r.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return fmt.Errorf("read row: %w", err)
	}

	for i, col := range row {
		fmt.Fprintf(out, "%d: %s\n", i, col)
	}

	return nil
}
