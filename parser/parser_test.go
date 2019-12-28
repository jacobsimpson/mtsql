package parser

import (
	"strings"
	"testing"

	"github.com/jacobsimpson/mtsql/ast"
	"github.com/jacobsimpson/mtsql/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParseEmptyQuery(t *testing.T) {
	assert := assert.New(t)

	q, err := Parse(lexer.NewFilterWhitespace(strings.NewReader("")))

	assert.Equal(err.Error(), `expected SELECT, found ""`)
	assert.Nil(q)
}

func TestParseString(t *testing.T) {
	assert := assert.New(t)

	q, err := Parse(lexer.NewFilterWhitespace(strings.NewReader("'sql string'")))

	assert.Equal(err.Error(), `expected SELECT, found "'sql string'"`)
	assert.Nil(q)
}

func TestParseMissingColumns(t *testing.T) {
	assert := assert.New(t)

	q, err := Parse(lexer.NewFilterWhitespace(strings.NewReader("SELECT ,,,")))

	assert.Equal(err.Error(), `expected column name, found ","`)
	assert.Nil(q)
}

func TestParseMissingFromClause(t *testing.T) {
	assert := assert.New(t)

	q, err := Parse(lexer.NewFilterWhitespace(strings.NewReader("SELECT column_name_1, c2")))

	assert.Equal(err.Error(), `expected FROM, found ""`)
	assert.Nil(q)
}

func TestParseSelectAllRows(t *testing.T) {
	assert := assert.New(t)

	q, err := Parse(lexer.NewFilterWhitespace(strings.NewReader("SELECT col1 FROM tablename")))

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

	q, err := Parse(lexer.NewFilterWhitespace(strings.NewReader("SELECT col1 FROM tablename WHERE col1 = 'abcd'")))

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

func TestParseSelectStar(t *testing.T) {
	assert := assert.New(t)

	q, err := Parse(lexer.NewFilterWhitespace(strings.NewReader("SELECT * FROM table_name")))

	assert.Nil(err)
	assert.NotNil(q)

	swf, ok := q.(*ast.SFW)
	assert.True(ok)

	assert.NotNil(swf.SelList)
	assert.Equal(len(swf.SelList.Attributes), 1)
	assert.Equal(swf.SelList.Attributes[0].Name, "*")

	assert.NotNil(swf.From)
	rel, ok := swf.From.(*ast.Relation)
	assert.True(ok)
	assert.Equal(rel.Name, "table_name")
}

func TestParseSelectTableQualifiedFieldNames(t *testing.T) {
	assert := assert.New(t)

	q, err := Parse(lexer.NewFilterWhitespace(strings.NewReader("SELECT table_name.id, name FROM table_name")))

	assert.Nil(err)
	assert.NotNil(q)

	swf, ok := q.(*ast.SFW)
	assert.True(ok)

	assert.NotNil(swf.SelList)
	assert.Equal(len(swf.SelList.Attributes), 2)
	assert.Equal(swf.SelList.Attributes[0], &ast.Attribute{Qualifier: "table_name", Name: "id"})
	assert.Equal(swf.SelList.Attributes[1], &ast.Attribute{Qualifier: "", Name: "name"})

	assert.NotNil(swf.From)
	rel, ok := swf.From.(*ast.Relation)
	assert.True(ok)
	assert.Equal(rel.Name, "table_name")
}

func TestParseSelectColumnAlias(t *testing.T) {
	assert := assert.New(t)

	q, err := Parse(lexer.NewFilterWhitespace(strings.NewReader("SELECT id the_special_id, name FROM table_name")))

	assert.Nil(err)
	assert.NotNil(q)

	swf, ok := q.(*ast.SFW)
	assert.True(ok)

	assert.NotNil(swf.SelList)
	assert.Equal(len(swf.SelList.Attributes), 2)
	assert.Equal(swf.SelList.Attributes[0], &ast.Attribute{Name: "id", Alias: "the_special_id"})
	assert.Equal(swf.SelList.Attributes[1], &ast.Attribute{Name: "name"})

	assert.NotNil(swf.From)
	rel, ok := swf.From.(*ast.Relation)
	assert.True(ok)
	assert.Equal(rel.Name, "table_name")
}

func TestAttribute(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name     string
		input    string
		expected *ast.Attribute
		err      error
	}{
		{
			name:     "single simple attribute",
			input:    "a_name",
			expected: &ast.Attribute{Name: "a_name"},
		},
		{
			name:     "single qualified attribute",
			input:    "qual.a_name",
			expected: &ast.Attribute{Qualifier: "qual", Name: "a_name"},
		},
		{
			name:     "single qualified attribute with alias",
			input:    "qual.a_name al",
			expected: &ast.Attribute{Qualifier: "qual", Name: "a_name", Alias: "al"},
		},
		{
			name:     "multiple simple attribute",
			input:    "a_name, b_name",
			expected: &ast.Attribute{Name: "a_name"},
		},
		{
			name:     "multiple qualified attribute",
			input:    "qual.a_name, qual.b_name",
			expected: &ast.Attribute{Qualifier: "qual", Name: "a_name"},
		},
		{
			name:     "multiple qualified attribute with alias",
			input:    "qual.a_name alias, qual.b_name",
			expected: &ast.Attribute{Qualifier: "qual", Name: "a_name", Alias: "alias"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a, err := attribute(lexer.NewFilterWhitespace(strings.NewReader(test.input)))
			assert.Nil(err)
			assert.Equal(a, test.expected)
		})
	}
}
