package physical

import (
	"fmt"
	"strconv"

	"github.com/jacobsimpson/mtsql/ast"
)

type filter struct {
	rowReader    RowReader
	columnName   string
	columnNumber int
	value        *ast.Constant
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
		switch t.value.Type {
		case ast.StringType:
			if row[t.columnNumber] == t.value.Value.(string) {
				return row, nil
			}
		case ast.IntegerType:
			i, err := strconv.Atoi(row[t.columnNumber])
			if err == nil && i == t.value.Value.(int) {
				return row, nil
			}
		}
	}
}

func (t *filter) Close() {}

func NewFilter(rowReader RowReader, columnName string, value *ast.Constant) (RowReader, error) {
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
