package physical

import (
	"testing"

	"github.com/jacobsimpson/mtsql/metadata"
	"github.com/stretchr/testify/assert"
)

func TestProjectOneColumn(t *testing.T) {
	assert := assert.New(t)
	rowReader := memoryScan{
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
	rowReader := memoryScan{
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
