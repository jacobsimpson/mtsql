package physical_test

import (
	"testing"

	"github.com/jacobsimpson/csvsql/physical"
	"github.com/stretchr/testify/assert"
)

type mockRowReader struct {
	columns []string
	rows    [][]string
	next    int
}

func (m *mockRowReader) Columns() []string { return m.columns }
func (m *mockRowReader) Read() ([]string, error) {
	if m.next >= len(m.rows) {
		return nil, nil
	}
	row := m.rows[m.next]
	m.next++
	return row, nil
}
func (m *mockRowReader) Close() {}

func TestProjectOneColumn(t *testing.T) {
	assert := assert.New(t)
	rowReader := mockRowReader{
		columns: []string{"col1", "col2", "col3"},
		rows: [][]string{
			[]string{"row1-col1", "row1-col2", "row1-col3"},
			[]string{"row2-col1", "row2-col2", "row2-col3"},
		},
	}

	proj, err := physical.NewProjection(&rowReader, []string{"col3"})
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
	rowReader := mockRowReader{
		columns: []string{"col1", "col2", "col3"},
		rows: [][]string{
			[]string{"row1-col1", "row1-col2", "row1-col3"},
			[]string{"row2-col1", "row2-col2", "row2-col3"},
		},
	}

	proj, err := physical.NewProjection(&rowReader, []string{"col3", "col1"})
	assert.Nil(err)
	assert.NotNil(proj)

	columns := proj.Columns()
	assert.Equal([]string{"col3", "col1"}, columns)

	r, err := proj.Read()
	assert.Nil(err)
	assert.Equal([]string{"row1-col3", "row1-col1"}, r)
}
