package physical

import (
	"io"
	"testing"

	"github.com/jacobsimpson/mtsql/metadata"
	"github.com/stretchr/testify/assert"
)

type MockRowReader struct {
	columns []*metadata.Column
	rows    [][]string
	next    int
}

func (m *MockRowReader) Columns() []*metadata.Column { return m.columns }
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
		columns: []*metadata.Column{
			&metadata.Column{Qualifier: "tb1", Name: "col1"},
			&metadata.Column{Qualifier: "tb1", Name: "col2"},
			&metadata.Column{Qualifier: "tb1", Name: "col3"},
		},
		rows: [][]string{
			[]string{"row1-col1", "row1-col2", "row1-col3"},
			[]string{"row2-col1", "row2-col2", "row2-col3"},
		},
	}

	proj, err := NewProjection(
		&rowReader,
		[]*metadata.Column{&metadata.Column{Qualifier: "tb1", Name: "col3"}})
	assert.Nil(err)
	assert.NotNil(proj)

	columns := proj.Columns()
	assert.Equal([]*metadata.Column{
		&metadata.Column{Qualifier: "tb1", Name: "col3"},
	}, columns)

	r, err := proj.Read()
	assert.Nil(err)
	assert.Equal([]string{"row1-col3"}, r)
}

func TestProjectTwoColumn(t *testing.T) {
	assert := assert.New(t)
	rowReader := MockRowReader{
		columns: []*metadata.Column{
			&metadata.Column{Qualifier: "tb1", Name: "col1"},
			&metadata.Column{Qualifier: "tb1", Name: "col2"},
			&metadata.Column{Qualifier: "tb1", Name: "col3"},
		},
		rows: [][]string{
			[]string{"row1-col1", "row1-col2", "row1-col3"},
			[]string{"row2-col1", "row2-col2", "row2-col3"},
		},
	}

	proj, err := NewProjection(&rowReader, []*metadata.Column{
		&metadata.Column{Qualifier: "tb1", Name: "col3"},
		&metadata.Column{Qualifier: "tb1", Name: "col1"},
	})
	assert.Nil(err)
	assert.NotNil(proj)

	columns := proj.Columns()
	assert.Equal([]*metadata.Column{
		&metadata.Column{Qualifier: "tb1", Name: "col3"},
		&metadata.Column{Qualifier: "tb1", Name: "col1"},
	}, columns)

	r, err := proj.Read()
	assert.Nil(err)
	assert.Equal([]string{"row1-col3", "row1-col1"}, r)
}
