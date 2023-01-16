package col

import (
	"fmt"

	csvtools "github.com/palsivertsen/csv-tools"
)

type FilterReader struct {
	Reader  csvtools.Reader
	Matcher RowMatcher
}

func (r *FilterReader) Read() ([]string, error) {
	for {
		row, err := r.Reader.Read()
		if err != nil {
			return nil, fmt.Errorf("read row: %w", err)
		}

		match, err := r.Matcher.Match(row)
		if err != nil {
			return nil, fmt.Errorf("match row: %w", err)
		}

		if match {
			return row, nil
		}
	}
}

type RowMatcher interface {
	Match([]string) (bool, error)
}

type AllRowMatcher []RowMatcher

func (m AllRowMatcher) Match(row []string) (bool, error) {
	for i, matcher := range m {
		match, err := matcher.Match(row)
		if err != nil {
			return false, fmt.Errorf("%d matcher error: %w", i, err)
		}

		if !match {
			return false, nil
		}
	}

	return true, nil
}

type Matcher interface {
	Match(string) (bool, error)
}

type StringMatcherFunc func(string) (bool, error)

func (f StringMatcherFunc) Match(s string) (bool, error) { return f(s) }

type IndexMatcher struct {
	Index   int
	Matcher Matcher
}

func (m *IndexMatcher) Match(row []string) (bool, error) {
	if m.Index >= len(row) {
		return false, csvtools.MissingColumnError{Actual: len(row), Expected: m.Index}
	}

	col := row[m.Index]
	match, err := m.Matcher.Match(col)
	if err != nil {
		return false, fmt.Errorf("match column: %w", err)
	}

	return match, nil
}
