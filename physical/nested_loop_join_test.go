package physical

import (
	"testing"

	"github.com/jacobsimpson/mtsql/metadata"
	"github.com/stretchr/testify/assert"
)

func TestNestedLoopJoin(t *testing.T) {
	assert := assert.New(t)
	left := &MockRowReader{
		columns: []*metadata.Column{
			&metadata.Column{Qualifier: "tb1", Name: "lc1"},
			&metadata.Column{Qualifier: "tb1", Name: "lc2"},
			&metadata.Column{Qualifier: "tb1", Name: "lc3"},
		},
		rows: [][]string{
			[]string{"row1-col1", "left-row1-col2", "left-row1-col3"},
			[]string{"row2-col1", "left-row2-col2", "left-row2-col3"},
		},
	}
	right := &MockRowReader{
		columns: []*metadata.Column{
			&metadata.Column{Qualifier: "tb2", Name: "rc1"},
			&metadata.Column{Qualifier: "tb2", Name: "rc2"},
			&metadata.Column{Qualifier: "tb2", Name: "rc3"},
		},
		rows: [][]string{
			[]string{"row1-col1", "right-row1-col2"},
			[]string{"row2-col1", "right-row2-col2"},
		},
	}
	rr, err := NewNestedLoopJoin(left, right)
	assert.Nil(err)
	assert.NotNil(rr)

	row, err := rr.Read()
	assert.Nil(err)
	assert.Equal([]string{"row1-col1", "left-row1-col2", "left-row1-col3", "row1-col1", "right-row1-col2"}, row)
	row, err = rr.Read()
	assert.Nil(err)
	assert.Equal([]string{"row1-col1", "left-row1-col2", "left-row1-col3", "row2-col1", "right-row2-col2"}, row)
	row, err = rr.Read()
	assert.Nil(err)
	assert.Equal([]string{"row2-col1", "left-row2-col2", "left-row2-col3", "row1-col1", "right-row1-col2"}, row)
	row, err = rr.Read()
	assert.Nil(err)
	assert.Equal([]string{"row2-col1", "left-row2-col2", "left-row2-col3", "row2-col1", "right-row2-col2"}, row)
}
