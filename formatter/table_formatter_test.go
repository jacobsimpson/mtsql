package formatter

import (
	"strings"
	"testing"

	"github.com/jacobsimpson/mtsql/metadata"
	"github.com/jacobsimpson/mtsql/physical"
	"github.com/stretchr/testify/assert"
)

func TestTableFormatter(t *testing.T) {
	assert := assert.New(t)
	rowReader := physical.NewMemoryScan(
		[]*metadata.Column{
			&metadata.Column{Qualifier: "tb1", Name: "col1"},
			&metadata.Column{Qualifier: "tb1", Name: "col2"},
			&metadata.Column{Qualifier: "tb1", Name: "col3"},
		},
		[][]string{
			[]string{"row1-col1", "row1-col2", "row1-col3"},
			[]string{"row2-col1", "row2-col2", "row2-col3"},
		},
	)
	formatter := NewTableFormatter(rowReader)
	var builder strings.Builder
	formatter.Print(&builder)

	// The output contains a separater between the column names and rows that
	// varies depending on terminal width, so exclude that from the testing
	// result.
	result := strings.Split(builder.String(), "\n")
	columnNames := result[0]
	rows := result[2:]
	assert.Equal("tb1.col1 tb1.col2 tb1.col3 ", columnNames)
	assert.Equal(
		[]string{
			"row1-col1 row1-col2 row1-col3 ",
			"row2-col1 row2-col2 row2-col3 ",
			"",
			"",
			"EOF",
			"",
			"",
		},
		rows)
}
