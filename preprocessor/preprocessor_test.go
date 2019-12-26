package preprocessor_test

import (
	"strings"
	"testing"

	"github.com/jacobsimpson/mtsql/lexer"
	"github.com/jacobsimpson/mtsql/parser"
	"github.com/jacobsimpson/mtsql/preprocessor"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	assert := assert.New(t)
	q, err := parser.Parse(lexer.NewFilterWhitespace(strings.NewReader("SELECT a, b, c FROM mock_table WHERE col1='raw_value'")))
	assert.Nil(err)

	err = preprocessor.Validate(q)

	assert.NotNil(err)
}
