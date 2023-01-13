package col_test

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/palsivertsen/csv-tools/csv/col"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoverReader_Read(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in      string
		columns []int
		out     string
	}{
		{},
		{
			in:  "a,b,c",
			out: "a,b,c",
		},
		{
			in:      "a,b,c",
			columns: []int{1},
			out:     "a,c",
		},
	}
	for testNum, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", testNum), func(t *testing.T) {
			t.Parallel()
			t.Logf("input:\n%s", spew.Sdump(tt))

			unit := col.NewRemoverReader(
				csv.NewReader(strings.NewReader(tt.in)),
				tt.columns...,
			)

			var out strings.Builder
			writer := csv.NewWriter(&out)
			for {
				row, err := unit.Read()
				if err != nil {
					require.ErrorIs(t, err, io.EOF)
					break
				}

				require.NoError(t, writer.Write(row))
			}

			writer.Flush()

			assert.Equal(t, tt.out, strings.TrimSpace(out.String()))
		})
	}
}
