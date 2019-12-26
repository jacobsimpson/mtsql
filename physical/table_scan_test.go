package physical_test

import (
	"testing"

	"github.com/jacobsimpson/csvsql/physical"
	"github.com/stretchr/testify/assert"
)

func TestReadOneRow(t *testing.T) {
	assert := assert.New(t)

	rowReader, err := physical.NewTableScan("testdata/cities.csv")
	assert.Nil(err)

	assert.Equal(rowReader.Columns()[0], "LatD")
	assert.Equal(rowReader.Columns()[8], "City")

	row, err := rowReader.Read()
	assert.Nil(err)
	assert.Equal(row[0], "41")
	assert.Equal(row[1], "5")
	//41,5,59,"N",80,39,0,"W","Youngstown",OH
}
