package physical

import (
	"fmt"

	"github.com/jacobsimpson/mtsql/metadata"
)

type columnFilter struct {
	rowReader  RowReader
	left       *metadata.Column
	leftIndex  int
	right      *metadata.Column
	rightIndex int
}

func NewColumnFilter(rowReader RowReader, left, right *metadata.Column) (RowReader, error) {
	lIdx, err := findColumn(left, rowReader.Columns())
	if err != nil {
		return nil, err
	}
	rIdx, err := findColumn(right, rowReader.Columns())
	if err != nil {
		return nil, err
	}
	return &columnFilter{
		rowReader:  rowReader,
		left:       left,
		leftIndex:  lIdx,
		right:      right,
		rightIndex: rIdx,
	}, nil
}

func (t *columnFilter) Columns() []*metadata.Column {
	return t.rowReader.Columns()
}

func (t *columnFilter) Read() ([]string, error) {
	for {
		row, err := t.rowReader.Read()
		if err != nil {
			return nil, err
		}
		if row == nil {
			return nil, nil
		}
		if row[t.leftIndex] == row[t.rightIndex] {
			return row, nil
		}
	}
}

func (t *columnFilter) Close()       {}
func (t *columnFilter) Reset() error { return t.rowReader.Reset() }

func (t *columnFilter) PlanDescription() *PlanDescription {
	return &PlanDescription{
		Name:        "ColumnFilter",
		Description: fmt.Sprintf("%s = %s", t.left.QualifiedName(), t.right.QualifiedName()),
	}
}

func (t *columnFilter) Children() []RowReader { return []RowReader{t.rowReader} }

func findColumn(target *metadata.Column, columns []*metadata.Column) (int, error) {
	for i, c := range columns {
		if (target.Qualifier == "" && target.Name == c.Name) ||
			(target.Qualifier == c.Qualifier && target.Name == c.Name) {
			return i, nil
		}
	}
	return 0, fmt.Errorf("column %q does not exist in relation", target.QualifiedName())
}
