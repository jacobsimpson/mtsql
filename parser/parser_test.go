package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jacobsimpson/mtsql/ast"
	"github.com/jacobsimpson/mtsql/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParseEmptyQuery(t *testing.T) {
	assert := assert.New(t)

	q, err := Parse(lexer.NewFilterWhitespace(strings.NewReader("")))

	assert.Equal(`expected SELECT or PROFILE`, err.Error())
	assert.Nil(q)
}

func TestParseString(t *testing.T) {
	assert := assert.New(t)

	q, err := Parse(lexer.NewFilterWhitespace(strings.NewReader("'sql string'")))

	assert.Equal(`expected SELECT or PROFILE`, err.Error())
	assert.Nil(q)
}

func TestParseMissingColumns(t *testing.T) {
	assert := assert.New(t)

	q, err := Parse(lexer.NewFilterWhitespace(strings.NewReader("SELECT ,,,")))

	assert.Equal(`expected column name, found ","`, err.Error())
	assert.Nil(q)
}

func TestParseMissingFromClause(t *testing.T) {
	assert := assert.New(t)

	q, err := Parse(lexer.NewFilterWhitespace(strings.NewReader("SELECT column_name_1, c2")))

	assert.Equal(`expected FROM clause`, err.Error())
	assert.Nil(q)
}

func TestParseSelectAllRows(t *testing.T) {
	tests := []struct {
		query    string
		expected ast.Query
		err      error
	}{
		{
			query: "SELECT col1 FROM tablename",
			expected: &ast.SFW{
				SelList: &ast.SelList{
					Attributes: []*ast.Attribute{
						{Name: "col1"},
					},
				},
				From: &ast.Relation{Name: "tablename"},
			},
		},
		{
			query: "SELECT col1 FROM tablename WHERE col1 = 'abcd'",
			expected: &ast.SFW{
				SelList: &ast.SelList{
					Attributes: []*ast.Attribute{
						{Name: "col1"},
					},
				},
				From: &ast.Relation{Name: "tablename"},
				Where: &ast.EqualCondition{
					LHS: &ast.Attribute{Name: "col1"},
					RHS: &ast.Constant{Type: ast.StringType, Value: "abcd", Raw: "'abcd'"},
				},
			},
		},
		{
			query: "SELECT * FROM books INNER JOIN ratings ON books.book_id = ratings.book_id",
			expected: &ast.SFW{
				SelList: &ast.SelList{
					Attributes: []*ast.Attribute{
						{Name: "*"},
					},
				},
				From: &ast.InnerJoin{
					Left:  &ast.Relation{Name: "books"},
					Right: &ast.Relation{Name: "ratings"},
					On: &ast.EqualColumnCondition{
						Left:  &ast.Attribute{Qualifier: "books", Name: "book_id"},
						Right: &ast.Attribute{Qualifier: "ratings", Name: "book_id"},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.query, func(t *testing.T) {
			assert := assert.New(t)

			q, err := Parse(lexer.NewFilterWhitespace(strings.NewReader(test.query)))

			assert.Equal(test.err, err)
			assert.Equal(test.expected, q)
		})
	}
}

func TestParseSelectStar(t *testing.T) {
	assert := assert.New(t)

	q, err := Parse(lexer.NewFilterWhitespace(strings.NewReader("SELECT * FROM table_name")))

	assert.Nil(err)
	assert.NotNil(q)

	swf, ok := q.(*ast.SFW)
	assert.True(ok)

	assert.NotNil(swf.SelList)
	assert.Equal(1, len(swf.SelList.Attributes))
	assert.Equal("*", swf.SelList.Attributes[0].Name)

	assert.NotNil(swf.From)
	rel, ok := swf.From.(*ast.Relation)
	assert.True(ok)
	assert.Equal("table_name", rel.Name)
}

func TestParseSelectTableQualifiedFieldNames(t *testing.T) {
	assert := assert.New(t)

	q, err := Parse(lexer.NewFilterWhitespace(strings.NewReader("SELECT table_name.id, name FROM table_name")))

	assert.Nil(err)
	assert.NotNil(q)

	swf, ok := q.(*ast.SFW)
	assert.True(ok)

	assert.NotNil(swf.SelList)
	assert.Equal(2, len(swf.SelList.Attributes))
	assert.Equal(&ast.Attribute{Qualifier: "table_name", Name: "id"}, swf.SelList.Attributes[0])
	assert.Equal(&ast.Attribute{Qualifier: "", Name: "name"}, swf.SelList.Attributes[1])

	assert.NotNil(swf.From)
	rel, ok := swf.From.(*ast.Relation)
	assert.True(ok)
	assert.Equal("table_name", rel.Name)
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
	assert.Equal("table_name", rel.Name)
}

func TestAttribute(t *testing.T) {
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
			assert := assert.New(t)
			a, err := attribute(lexer.NewFilterWhitespace(strings.NewReader(test.input)))
			assert.Nil(err)
			assert.Equal(test.expected, a)
		})
	}
}

func TestFrom(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ast.From
		err      error
	}{
		{
			name:     "simple table",
			input:    "from a_name",
			expected: &ast.Relation{Name: "a_name"},
		},
		{
			name:  "basic inner join",
			input: "from tab1 INNER JOIN tab2 ON tab1.id = tab2.id",
			expected: &ast.InnerJoin{
				Left:  &ast.Relation{Name: "tab1"},
				Right: &ast.Relation{Name: "tab2"},
				On: &ast.EqualColumnCondition{
					Left:  &ast.Attribute{Qualifier: "tab1", Name: "id"},
					Right: &ast.Attribute{Qualifier: "tab2", Name: "id"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			a, err := from(lexer.NewFilterWhitespace(strings.NewReader(test.input)))
			assert.Nil(err)
			assert.Equal(test.expected, a)
		})
	}
}

func TestFieldEqualsField(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ast.EqualColumnCondition
		err      error
	}{
		{
			name:  "field1 = field2",
			input: "field1 = field2",
			expected: &ast.EqualColumnCondition{
				Left:  &ast.Attribute{Name: "field1"},
				Right: &ast.Attribute{Name: "field2"},
			},
		},
		{
			name:  "qual.field1 = field2",
			input: "qual.field1 = field2",
			expected: &ast.EqualColumnCondition{
				Left:  &ast.Attribute{Qualifier: "qual", Name: "field1"},
				Right: &ast.Attribute{Name: "field2"},
			},
		},
		{
			name:  "field1 = qual.field2",
			input: "field1 = qual.field2",
			expected: &ast.EqualColumnCondition{
				Left:  &ast.Attribute{Name: "field1"},
				Right: &ast.Attribute{Qualifier: "qual", Name: "field2"},
			},
		},
		{
			name:  "qual1.field1 = qual.field2",
			input: "qual1.field1 = qual.field2",
			expected: &ast.EqualColumnCondition{
				Left:  &ast.Attribute{Qualifier: "qual1", Name: "field1"},
				Right: &ast.Attribute{Qualifier: "qual", Name: "field2"},
			},
		},
		{
			name:  "qual1.field1",
			input: "qual1.field1",
			err:   fmt.Errorf(`expected = after "field1"`),
		},
		{
			name:  "qual1.field1 =",
			input: "qual1.field1 =",
			err:   fmt.Errorf(`expected field name, found ""`),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			ecc, err := fieldEqualsField(lexer.NewFilterWhitespace(strings.NewReader(test.input)))

			assert.Equal(test.err, err)
			assert.Equal(test.expected, ecc)
		})
	}
}

func TestField(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ast.Attribute
		err      error
	}{
		{
			name:     "simple field anem",
			input:    "a_name",
			expected: &ast.Attribute{Name: "a_name"},
		},
		{
			name:     "simple qualified attribute",
			input:    "qual.a_name",
			expected: &ast.Attribute{Qualifier: "qual", Name: "a_name"},
		},
		{
			name:  "qualified field name missing name",
			input: "qual.",
			err:   fmt.Errorf(`expected field name, found ""`),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			a, err := field(lexer.NewFilterWhitespace(strings.NewReader(test.input)))
			assert.Equal(test.err, err)
			assert.Equal(test.expected, a)
		})
	}
}

func TestCondition(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ast.Condition
		err      error
	}{
		{
			name:  "a = 'abc'",
			input: "a = 'abc'",
			expected: &ast.EqualCondition{
				LHS: &ast.Attribute{Name: "a"},
				RHS: &ast.Constant{Type: ast.StringType, Value: "abc", Raw: "'abc'"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			a, err := condition(lexer.NewFilterWhitespace(strings.NewReader(test.input)))

			assert.Equal(test.err, err)
			assert.Equal(test.expected, a)
		})
	}
}

func TestOrderByClause(t *testing.T) {
	tests := []struct {
		input    string
		expected *ast.OrderCriteria
		err      error
	}{
		{
			input: "a",
			expected: &ast.OrderCriteria{
				Attribute: &ast.Attribute{Name: "a"},
				SortOrder: ast.Asc,
			},
		},
		{
			input: "a ASC",
			expected: &ast.OrderCriteria{
				Attribute: &ast.Attribute{Name: "a"},
				SortOrder: ast.Asc,
			},
		},
		{
			input: "table_name.field_name DESC",
			expected: &ast.OrderCriteria{
				Attribute: &ast.Attribute{Qualifier: "table_name", Name: "field_name"},
				SortOrder: ast.Desc,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			assert := assert.New(t)

			obc, err := orderByClause(lexer.NewFilterWhitespace(strings.NewReader(test.input)))

			assert.Equal(test.err, err)
			assert.Equal(test.expected, obc)
		})
	}
}

func TestOrderBy(t *testing.T) {
	tests := []struct {
		input    string
		expected *ast.OrderBy
		err      error
	}{
		{
			input: "order by a",
			expected: &ast.OrderBy{
				Criteria: []*ast.OrderCriteria{
					{
						Attribute: &ast.Attribute{Name: "a"},
						SortOrder: ast.Asc,
					},
				},
			},
		},
		{
			input: "order by a ASC, b DESC",
			expected: &ast.OrderBy{
				Criteria: []*ast.OrderCriteria{
					{Attribute: &ast.Attribute{Name: "a"}, SortOrder: ast.Asc},
					{Attribute: &ast.Attribute{Name: "b"}, SortOrder: ast.Desc},
				},
			},
		},
		{
			input: "ORDER BY tn.fn DESC, bbb ASC, cde DESC",
			expected: &ast.OrderBy{
				Criteria: []*ast.OrderCriteria{
					{Attribute: &ast.Attribute{Qualifier: "tn", Name: "fn"}, SortOrder: ast.Desc},
					{Attribute: &ast.Attribute{Name: "bbb"}, SortOrder: ast.Asc},
					{Attribute: &ast.Attribute{Name: "cde"}, SortOrder: ast.Desc},
				},
			},
		},
		{
			input: "ORDER BY tn.fn KKKKKKK",
			expected: &ast.OrderBy{
				Criteria: []*ast.OrderCriteria{
					{Attribute: &ast.Attribute{Qualifier: "tn", Name: "fn"}, SortOrder: ast.Asc},
				},
			},
		},
		{
			input: "group by abcd",
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			assert := assert.New(t)

			ob, err := orderBy(lexer.NewFilterWhitespace(strings.NewReader(test.input)))

			assert.Equal(test.err, err)
			assert.Equal(test.expected, ob)
		})
	}
}

func TestWhere(t *testing.T) {
	tests := []struct {
		input    string
		expected ast.Condition
		err      error
	}{
		{
			input: "WHERE a = 'abc'",
			expected: &ast.EqualCondition{
				LHS: &ast.Attribute{Name: "a"},
				RHS: &ast.Constant{Type: ast.StringType, Value: "abc", Raw: "'abc'"},
			},
		},
		{
			input: "ORDER BY table_name.field_name",
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			assert := assert.New(t)

			ob, err := where(lexer.NewFilterWhitespace(strings.NewReader(test.input)))

			assert.Equal(test.err, err)
			assert.Equal(test.expected, ob)
		})
	}
}

func TestParseInnerJoin(t *testing.T) {
	tests := []struct {
		table    *ast.Relation
		input    string
		expected *ast.InnerJoin
		err      error
	}{
		{
			table: &ast.Relation{Name: "tab1"},
			input: "INNER JOIN tab2 ON tab1.id = tab2.id",
			expected: &ast.InnerJoin{
				Left:  &ast.Relation{Name: "tab1"},
				Right: &ast.Relation{Name: "tab2"},
				On: &ast.EqualColumnCondition{
					Left:  &ast.Attribute{Qualifier: "tab1", Name: "id"},
					Right: &ast.Attribute{Qualifier: "tab2", Name: "id"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			assert := assert.New(t)

			from, err := innerJoin(lexer.NewFilterWhitespace(strings.NewReader(test.input)), test.table)

			assert.Equal(test.err, err)
			assert.Equal(test.expected, from)
		})
	}
}
