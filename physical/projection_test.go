package physical

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockRowReader struct {
	columns []string
	rows    [][]string
	next    int
}

func (m *MockRowReader) Columns() []string { return m.columns }
func (m *MockRowReader) Read() ([]string, error) {
	if m.next >= len(m.rows) {
		return nil, io.EOF
	}
	row := m.rows[m.next]
	m.next++
	return row, nil
}
func (m *MockRowReader) Close() {}
func (m *MockRowReader) Reset() error {
	m.next = 0
	return nil
}
func (m *MockRowReader) PlanDescription() *PlanDescription { return nil }
func (m *MockRowReader) Children() []RowReader             { return nil }

func TestProjectOneColumn(t *testing.T) {
	assert := assert.New(t)
	rowReader := MockRowReader{
		columns: []string{"col1", "col2", "col3"},
		rows: [][]string{
			[]string{"row1-col1", "row1-col2", "row1-col3"},
			[]string{"row2-col1", "row2-col2", "row2-col3"},
		},
	}

	proj, err := NewProjection(&rowReader, []string{"col3"})
	assert.Nil(err)
	assert.NotNil(proj)

	columns := proj.Columns()
	assert.Equal([]string{"col3"}, columns)

	r, err := proj.Read()
	assert.Nil(err)
	assert.Equal([]string{"row1-col3"}, r)
}

func TestProjectTwoColumn(t *testing.T) {
	assert := assert.New(t)
	rowReader := MockRowReader{
		columns: []string{"col1", "col2", "col3"},
		rows: [][]string{
			[]string{"row1-col1", "row1-col2", "row1-col3"},
			[]string{"row2-col1", "row2-col2", "row2-col3"},
		},
	}

	proj, err := NewProjection(&rowReader, []string{"col3", "col1"})
	assert.Nil(err)
	assert.NotNil(proj)

	columns := proj.Columns()
	assert.Equal([]string{"col3", "col1"}, columns)

	r, err := proj.Read()
	assert.Nil(err)
	assert.Equal([]string{"row1-col3", "row1-col1"}, r)
}
