package physical

import (
	"fmt"
	"strconv"

	"github.com/jacobsimpson/mtsql/ast"
	"github.com/jacobsimpson/mtsql/metadata"
)

type filter struct {
	rowReader    RowReader
	column       *metadata.Column
	columnNumber int
	value        *ast.Constant
}

func (t *filter) Columns() []*metadata.Column {
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

func (t *filter) Close()       {}
func (t *filter) Reset() error { return t.rowReader.Reset() }

func (t *filter) PlanDescription() *PlanDescription {
	return &PlanDescription{
		Name:        "Filter",
		Description: fmt.Sprintf("%s = %v", t.column.QualifiedName(), t.value.Value),
	}
}

func (t *filter) Children() []RowReader { return []RowReader{t.rowReader} }

func NewFilter(rowReader RowReader, column *metadata.Column, value *ast.Constant) (RowReader, error) {
	n := -1
	for i, c := range rowReader.Columns() {
		if c.Qualifier == column.Qualifier && c.Name == column.Name {
			n = i
		}
	}
	if n < 0 {
		return nil, fmt.Errorf("column %q does not exist in relation", column.QualifiedName())
	}
	return &filter{
		rowReader:    rowReader,
		column:       column,
		columnNumber: n,
		value:        value,
	}, nil
}
