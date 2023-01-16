package csvtools

import (
	"errors"
	"fmt"
	"io"
)

type Reader interface {
	Read() ([]string, error)
}

type Writer interface {
	Write([]string) error
}

func Copy(dst Writer, src Reader) (int, error) {
	var count int

	for {
		record, err := src.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return count, fmt.Errorf("read record: %w", err)
		}

		if err := dst.Write(record); err != nil {
			return count, fmt.Errorf("write record: %w", err)
		}
		count++
	}

	return count, nil
}

type RowError struct {
	RowNumber int
	Err       error
}

func (e RowError) Error() string {
	return fmt.Sprintf("row %d: %s", e.RowNumber, e.Err)
}

func (e RowError) Unwrap() error {
	return e.Err
}

type MissingColumnError struct {
	Actual, Expected int
}

func (e MissingColumnError) Error() string {
	return fmt.Sprintf("row too short: %d (expected %d)", e.Actual, e.Expected)
}
