package col

import (
	"fmt"

	csvtools "github.com/palsivertsen/csv-tools"
)

type TransformFunc func([]string) ([]string, error)

type TransformReader struct {
	Reader    csvtools.Reader
	Transform TransformFunc
}

func (t *TransformReader) Read() ([]string, error) {
	row, err := t.Reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	transformed, err := t.Transform(row)
	if err != nil {
		return nil, fmt.Errorf("transform: %w", err)
	}

	return transformed, nil
}

func NewPickReader(r csvtools.Reader, columnIndexes ...int) *TransformReader {
	return &TransformReader{
		Reader: r,
		Transform: func(rowIn []string) ([]string, error) {
			rowOut := make([]string, 0, len(columnIndexes))
			for _, columnIndex := range columnIndexes {
				rowOut = append(rowOut, rowIn[columnIndex])
			}
			return rowOut, nil
		},
	}
}
