package physical

import (
	"fmt"
	"strings"

	"github.com/jacobsimpson/mtsql/metadata"
)

type projection struct {
	rowReader     RowReader
	columnIndexes []int
}

func (t *projection) Columns() []*metadata.Column {
	c := t.rowReader.Columns()
	r := []*metadata.Column{}
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

func (t *projection) Close()       {}
func (t *projection) Reset() error { return t.rowReader.Reset() }

func (t *projection) PlanDescription() *PlanDescription {
	columnNames := []string{}
	for _, c := range t.Columns() {
		columnNames = append(columnNames, c.QualifiedName())
	}
	return &PlanDescription{
		Name:        "Projection",
		Description: fmt.Sprintf("%v", strings.Join(columnNames, ", ")),
	}
}

func (t *projection) Children() []RowReader { return []RowReader{t.rowReader} }

func NewProjection(rowReader RowReader, columns []*metadata.Column) (RowReader, error) {
	columnMap := map[metadata.Column]int{}
	for i, c := range rowReader.Columns() {
		columnMap[*c] = i
	}
	cols := []int{}
	for _, c := range columns {
		cols = append(cols, columnMap[*c])
	}
	return &projection{
		rowReader:     rowReader,
		columnIndexes: cols,
	}, nil
}
