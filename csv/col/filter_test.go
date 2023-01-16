package col_test

import (
	"fmt"
	"io"
	"regexp"
	"testing"

	"github.com/davecgh/go-spew/spew"
	csvtools "github.com/palsivertsen/csv-tools"
	"github.com/palsivertsen/csv-tools/csv/col"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterReader_Read(t *testing.T) {
	t.Parallel()
	tests := []struct {
		reader       csvtools.Reader
		expectedRows [][]string
		filter       col.RowMatcher
	}{
		{
			reader:       &sliceReader{rows: [][]string{{"0"}, {"1"}, {"2"}}},
			expectedRows: [][]string{{"0"}, {"2"}},
			filter:       &rowNumberMatcher{rows: []int{0, 2}},
		},
	}
	for testNum, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", testNum), func(t *testing.T) {
			t.Parallel()
			t.Logf("input:\n%s", spew.Sdump(tt))

			unit := col.FilterReader{
				Reader:  tt.reader,
				Matcher: tt.filter,
			}

			for _, expectedRow := range tt.expectedRows {
				row, err := unit.Read()
				require.NoError(t, err)

				assert.Equal(t, expectedRow, row)
			}

			row, err := unit.Read()
			require.ErrorIs(t, err, io.EOF)
			assert.Empty(t, row)
		})
	}
}

func TestColumnMatcher_Match(t *testing.T) {
	t.Parallel()
	rxp := regexp.MustCompile("[[:alpha:]]")

	unit := col.IndexMatcher{
		Index: 1,
		Matcher: col.StringMatcherFunc(func(row string) (bool, error) {
			return rxp.MatchString(row), nil
		}),
	}

	{
		match, _ := unit.Match([]string{"123", "asd"})
		assert.True(t, match)
	}
	{
		match, _ := unit.Match([]string{"asd", "123"})
		assert.False(t, match)
	}
}

// helpers

type sliceReader struct {
	rows [][]string
	next int
}

func (r *sliceReader) Read() ([]string, error) {
	if r.next >= len(r.rows) {
		return nil, io.EOF
	}

	defer func() {
		r.next++
	}()

	return r.rows[r.next], nil
}

type rowNumberMatcher struct {
	rows  []int
	count int
}

func (m *rowNumberMatcher) Match(row []string) (bool, error) {
	defer func() {
		m.count++
	}()

	for _, row := range m.rows {
		if row == m.count {
			return true, nil
		}
	}

	return false, nil
}
