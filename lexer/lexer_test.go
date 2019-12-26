package lexer_test

import (
	"strings"
	"testing"

	"github.com/jacobsimpson/mtsql/lexer"
	"github.com/stretchr/testify/assert"
)

func TestNewLexer(t *testing.T) {
	assert := assert.New(t)
	l := lexer.New(strings.NewReader("SELECT * FROM table_name"))
	assert.NotNil(l)
}

func TestLexWhitespace(t *testing.T) {
	assert := assert.New(t)
	l := lexer.New(strings.NewReader(`
 `))

	l.Next()
	token := l.Token()

	assert.Equal(lexer.WhitespaceType, token.Type)
	assert.Equal(`
 `, token.Raw)

	l.Next()
	token = l.Token()

	assert.Equal(lexer.EOFType, token.Type)
	assert.Equal("", token.Raw)
}

func TestLexLotsOfWhitespace(t *testing.T) {
	assert := assert.New(t)
	l := lexer.New(strings.NewReader("           "))

	l.Next()
	token := l.Token()

	assert.Equal(lexer.WhitespaceType, token.Type)
	assert.Equal("           ", token.Raw)

	l.Next()
	token = l.Token()

	assert.Equal(lexer.EOFType, token.Type)
	assert.Equal("", token.Raw)
}

func TestLexIdentifierUnread(t *testing.T) {
	assert := assert.New(t)
	l := lexer.New(strings.NewReader("PROFILE SELECT col1 FROM table_name"))

	l.Next()
	token := l.Token()

	assert.Equal(lexer.IdentifierType, token.Type)
	assert.Equal("PROFILE", token.Raw)

	l.UnreadToken()
	token = l.Token()

	assert.Equal(lexer.IdentifierType, token.Type)
	assert.Equal("PROFILE", token.Raw)

}

func TestLexIdentifier(t *testing.T) {
	assert := assert.New(t)
	l := lexer.New(strings.NewReader("SELECT"))

	l.Next()
	token := l.Token()

	assert.Equal(lexer.IdentifierType, token.Type)
	assert.Equal("SELECT", token.Raw)
}

func TestLexIdentifiersAndWhitespace(t *testing.T) {
	assert := assert.New(t)
	l := lexer.New(strings.NewReader("SELECT col   FROM table_name"))
	expected := []lexer.Token{
		lexer.Token{Type: lexer.IdentifierType, Raw: "SELECT"},
		lexer.Token{Type: lexer.WhitespaceType, Raw: " "},
		lexer.Token{Type: lexer.IdentifierType, Raw: "col"},
		lexer.Token{Type: lexer.WhitespaceType, Raw: "   "},
		lexer.Token{Type: lexer.IdentifierType, Raw: "FROM"},
		lexer.Token{Type: lexer.WhitespaceType, Raw: " "},
		lexer.Token{Type: lexer.IdentifierType, Raw: "table_name"},
	}

	for _, t := range expected {
		l.Next()
		token := l.Token()

		assert.Equal(t.Type, token.Type)
		assert.Equal(t.Raw, token.Raw)
	}

}

func TestLexQuery(t *testing.T) {
	assert := assert.New(t)
	l := lexer.New(strings.NewReader("SELECT col1,col   FROM table_name WHERE col3 = 3 AND col_12_a='abcd'"))
	expected := []lexer.Token{
		lexer.Token{Type: lexer.IdentifierType, Raw: "SELECT"},
		lexer.Token{Type: lexer.WhitespaceType, Raw: " "},
		lexer.Token{Type: lexer.IdentifierType, Raw: "col1"},
		lexer.Token{Type: lexer.CommaType, Raw: ","},
		lexer.Token{Type: lexer.IdentifierType, Raw: "col"},
		lexer.Token{Type: lexer.WhitespaceType, Raw: "   "},
		lexer.Token{Type: lexer.IdentifierType, Raw: "FROM"},
		lexer.Token{Type: lexer.WhitespaceType, Raw: " "},
		lexer.Token{Type: lexer.IdentifierType, Raw: "table_name"},
		lexer.Token{Type: lexer.WhitespaceType, Raw: " "},
		lexer.Token{Type: lexer.IdentifierType, Raw: "WHERE"},
		lexer.Token{Type: lexer.WhitespaceType, Raw: " "},
		lexer.Token{Type: lexer.IdentifierType, Raw: "col3"},
		lexer.Token{Type: lexer.WhitespaceType, Raw: " "},
		lexer.Token{Type: lexer.EqualType, Raw: "="},
		lexer.Token{Type: lexer.WhitespaceType, Raw: " "},
		lexer.Token{Type: lexer.IntegerType, Raw: "3"},
		lexer.Token{Type: lexer.WhitespaceType, Raw: " "},
		lexer.Token{Type: lexer.IdentifierType, Raw: "AND"},
		lexer.Token{Type: lexer.WhitespaceType, Raw: " "},
		lexer.Token{Type: lexer.IdentifierType, Raw: "col_12_a"},
		lexer.Token{Type: lexer.EqualType, Raw: "="},
		lexer.Token{Type: lexer.StringType, Raw: "'abcd'"},
		lexer.Token{Type: lexer.EOFType, Raw: ""},
	}

	for _, t := range expected {
		l.Next()
		token := l.Token()

		assert.Equal(t.Type, token.Type)
		assert.Equal(t.Raw, token.Raw)
	}
}

func TestLexAnotherQuery(t *testing.T) {
	assert := assert.New(t)
	l := lexer.NewFilterWhitespace(strings.NewReader("SELECT State, City FROM cities WHERE State = 123"))
	expected := []lexer.Token{
		lexer.Token{Type: lexer.IdentifierType, Raw: "SELECT"},
		lexer.Token{Type: lexer.IdentifierType, Raw: "State"},
		lexer.Token{Type: lexer.CommaType, Raw: ","},
		lexer.Token{Type: lexer.IdentifierType, Raw: "City"},
		lexer.Token{Type: lexer.IdentifierType, Raw: "FROM"},
		lexer.Token{Type: lexer.IdentifierType, Raw: "cities"},
		lexer.Token{Type: lexer.IdentifierType, Raw: "WHERE"},
		lexer.Token{Type: lexer.IdentifierType, Raw: "State"},
		lexer.Token{Type: lexer.EqualType, Raw: "="},
		lexer.Token{Type: lexer.IntegerType, Raw: "123"},
		lexer.Token{Type: lexer.EOFType, Raw: ""},
	}

	for _, t := range expected {
		l.Next()
		token := l.Token()

		assert.Equal(t.Type, token.Type)
		assert.Equal(t.Raw, token.Raw)
	}
}

func TestLexQualifiedIdentifier(t *testing.T) {
	assert := assert.New(t)
	l := lexer.NewFilterWhitespace(strings.NewReader("abc.def"))
	expected := []lexer.Token{
		lexer.Token{Type: lexer.IdentifierType, Raw: "abc"},
		lexer.Token{Type: lexer.PeriodType, Raw: "."},
		lexer.Token{Type: lexer.IdentifierType, Raw: "def"},
		lexer.Token{Type: lexer.EOFType, Raw: ""},
	}

	for _, t := range expected {
		l.Next()
		token := l.Token()

		assert.Equal(t.Type, token.Type)
		assert.Equal(t.Raw, token.Raw)
	}
}
