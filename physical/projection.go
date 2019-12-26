package physical

import (
	"fmt"
	"strings"
)

type projection struct {
	rowReader     RowReader
	columnIndexes []int
}

func (t *projection) Columns() []string {
	c := t.rowReader.Columns()
	r := []string{}
	for _, col := range t.columnIndexes {
		r = append(r, c[col])
	}
	return r
}

func (t *projection) Read() ([]string, error) {
	row, err := t.rowReader.Read()
	if err != nil {
		return nil, err
	}
	r := []string{}
	for _, col := range t.columnIndexes {
		r = append(r, row[col])
	}
	return r, nil
}

func (t *projection) Close() {}

func (t *projection) PlanDescription() *PlanDescription {
	return &PlanDescription{
		Name:        "Projection",
		Description: fmt.Sprintf("%v", strings.Join(t.Columns(), ", ")),
	}
}

func (t *projection) Children() []RowReader { return []RowReader{t.rowReader} }

func NewProjection(rowReader RowReader, columns []string) (RowReader, error) {
	columnMap := map[string]int{}
	for i, c := range rowReader.Columns() {
		columnMap[c] = i
	}
	cols := []int{}
	for _, c := range columns {
		cols = append(cols, columnMap[c])
	}
	return &projection{
		rowReader:     rowReader,
		columnIndexes: cols,
	}, nil
}
