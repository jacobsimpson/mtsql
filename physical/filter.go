package physical

import "fmt"

type filter struct {
	rowReader    RowReader
	columnName   string
	columnNumber int
	value        string
}

func (t *filter) Columns() []string {
	return t.rowReader.Columns()
}

func (t *filter) Read() ([]string, error) {
	for {
		row, err := t.rowReader.Read()
		if err != nil {
			return nil, err
		}
		if row == nil {
			return nil, nil
		}
		if row[t.columnNumber] == t.value {
			return row, nil
		}
	}
}

func (t *filter) Close() {}

func NewFilter(rowReader RowReader, columnName string, value string) (RowReader, error) {
	n := -1
	for i, c := range rowReader.Columns() {
		if c == columnName {
			n = i
		}
	}
	if n < 0 {
		return nil, fmt.Errorf("column %q does not exist in relation", columnName)
	}
	return &filter{
		rowReader:    rowReader,
		columnName:   columnName,
		columnNumber: n,
		value:        value,
	}, nil
}
