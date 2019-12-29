package physical_test

import (
	"testing"

	"github.com/jacobsimpson/mtsql/metadata"
	"github.com/jacobsimpson/mtsql/physical"
	"github.com/stretchr/testify/assert"
)

func TestReadOneRow(t *testing.T) {
	assert := assert.New(t)

	rowReader, err := physical.NewTableScan("cities", "testdata/cities.csv")
	assert.Nil(err)

	assert.Equal(
		&metadata.Column{Qualifier: "cities", Name: "LatD"},
		rowReader.Columns()[0])
	assert.Equal(
		&metadata.Column{Qualifier: "cities", Name: "City"},
		rowReader.Columns()[8])

	row, err := rowReader.Read()
	assert.Nil(err)
	assert.Equal("41", row[0])
	assert.Equal("5", row[1])
}
