package physical

import (
	"testing"

	"github.com/jacobsimpson/mtsql/ast"
	"github.com/stretchr/testify/assert"
)

func TestNestedLoopJoin(t *testing.T) {
	assert := assert.New(t)
	left := &MockRowReader{
		columns: []string{"lc1", "lc2", "lc3"},
		rows: [][]string{
			[]string{"row1-col1", "left-row1-col2", "left-row1-col3"},
			[]string{"row2-col1", "left-row2-col2", "left-row2-col3"},
		},
	}
	right := &MockRowReader{
		columns: []string{"rc1", "rc2", "rc3"},
		rows: [][]string{
			[]string{"row1-col1", "right-row1-col2"},
			[]string{"row2-col1", "right-row2-col2"},
		},
	}
	on := &ast.EqualColumnCondition{
		Left:  &ast.Attribute{Name: "lc1"},
		Right: &ast.Attribute{Name: "rc1"},
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
