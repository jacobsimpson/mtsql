package parser_test

import (
	"strings"
	"testing"

	"github.com/jacobsimpson/mtsql/ast"
	"github.com/jacobsimpson/mtsql/lexer"
	"github.com/jacobsimpson/mtsql/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseEmptyQuery(t *testing.T) {
	assert := assert.New(t)

	q, err := parser.Parse(lexer.NewFilterWhitespace(strings.NewReader("")))

	assert.Equal(err.Error(), `expected SELECT, found ""`)
	assert.Nil(q)
}

func TestParseString(t *testing.T) {
	assert := assert.New(t)

	q, err := parser.Parse(lexer.NewFilterWhitespace(strings.NewReader("'sql string'")))

	assert.Equal(err.Error(), `expected SELECT, found "'sql string'"`)
	assert.Nil(q)
}

func TestParseMissingColumns(t *testing.T) {
	assert := assert.New(t)

	q, err := parser.Parse(lexer.NewFilterWhitespace(strings.NewReader("SELECT ,,,")))

	assert.Equal(err.Error(), `expected column name, found ","`)
	assert.Nil(q)
}

func TestParseMissingFromClause(t *testing.T) {
	assert := assert.New(t)

	q, err := parser.Parse(lexer.NewFilterWhitespace(strings.NewReader("SELECT column_name_1, c2")))

	assert.Equal(err.Error(), `expected FROM, found ""`)
	assert.Nil(q)
}

func TestParseSelectAllRows(t *testing.T) {
	assert := assert.New(t)

	q, err := parser.Parse(lexer.NewFilterWhitespace(strings.NewReader("SELECT col1 FROM tablename")))

	assert.Nil(err)
	assert.NotNil(q)

	swf, ok := q.(*ast.SFW)
	assert.True(ok)

	assert.NotNil(swf.SelList)
	assert.Equal(len(swf.SelList.Attributes), 1)
	assert.Equal(swf.SelList.Attributes[0].Name, "col1")

	assert.NotNil(swf.From)
	rel, ok := swf.From.(*ast.Relation)
	assert.True(ok)
	assert.Equal(rel.Name, "tablename")

	assert.Nil(swf.Condition)
}

func TestParseSelectWhereColumnEqual(t *testing.T) {
	assert := assert.New(t)

	q, err := parser.Parse(lexer.NewFilterWhitespace(strings.NewReader("SELECT col1 FROM tablename WHERE col1 = 'abcd'")))

	assert.Nil(err)
	assert.NotNil(q)

	swf, ok := q.(*ast.SFW)
	assert.True(ok)

	assert.NotNil(swf.SelList)
	assert.Equal(len(swf.SelList.Attributes), 1)
	assert.Equal(swf.SelList.Attributes[0].Name, "col1")

	assert.NotNil(swf.From)
	rel, ok := swf.From.(*ast.Relation)
	assert.True(ok)
	assert.Equal(rel.Name, "tablename")

	assert.NotNil(swf.Condition)
	eq, ok := swf.Condition.(*ast.EqualCondition)
	assert.True(ok)
	assert.Equal(eq.LHS.Name, "col1")
	assert.Equal(eq.RHS.Raw, "'abcd'")
}
