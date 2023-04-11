package app_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/palsivertsen/csv-tools/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestApp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args []string
		in   string
		out  string
	}{
		// remove columns
		{
			args: strings.Split("csv column remove --columns 1", " "),
			in:   "a,b,c",
			out:  "a,c\n",
		},
		{
			args: strings.Split("csv column remove --columns 1,0", " "),
			in:   "a,b,c",
			out:  "c\n",
		},
		{
			args: strings.Split("csv column remove --columns 0 --columns 1", " "),
			in:   "a,b,c",
			out:  "c\n",
		},
		// print header
		{
			args: strings.Split("csv print header", " "),
			in:   "a,b,c\n1,2,3",
			out:  "0: a\n1: b\n2: c\n",
		},
		{
			args: strings.Split("csv print header", " "),
		},
		// filter
		{
			args: []string{"csv", "filter", "--column-regexp", "2:^ "},
			in:   "a,b,c\n1,2, 3",
			out:  "1,2, 3\n",
		},
		{
			args: []string{"csv", "filter", "--column-regexp", "0:[[:alpha:]]"},
			in:   "a,b,c\n1,2,3\n4,5,6",
			out:  "a,b,c\n",
		},
		// pick columns
		{
			args: strings.Split("csv column pick --columns 1,2", " "),
			in:   "a,b,c",
			out:  "b,c\n",
		},
	}
	for testNum, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", testNum), func(t *testing.T) {
			t.Parallel()
			t.Logf("input:\n%s", spew.Sdump(tt))

			var out bytes.Buffer

			unit := app.App()
			unit.Reader = strings.NewReader(tt.in)
			unit.Writer = &out
			unit.ExitErrHandler = func(*cli.Context, error) {}

			require.NoError(t, unit.Run(tt.args))
			assert.Equal(t, tt.out, out.String())
		})
	}
}
