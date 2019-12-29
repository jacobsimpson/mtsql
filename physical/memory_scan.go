package physical

import (
	"io"

	"github.com/jacobsimpson/mtsql/metadata"
)

type memoryScan struct {
	columns []*metadata.Column
	rows    [][]string
	next    int
}

func NewMemoryScan(columns []*metadata.Column, rows [][]string) RowReader {
	return &memoryScan{
		columns: columns,
		rows:    rows,
	}
}

func (m *memoryScan) Columns() []*metadata.Column {
	return m.columns
}

func (m *memoryScan) Read() ([]string, error) {
	if m.next >= len(m.rows) {
		return nil, io.EOF
	}
	row := m.rows[m.next]
	m.next++
	return row, nil
}

func (m *memoryScan) Close() {}

func (m *memoryScan) Reset() error {
	m.next = 0
	return nil
}

func (m *memoryScan) PlanDescription() *PlanDescription {
	return nil
}

func (m *memoryScan) Children() []RowReader {
	return nil
}
