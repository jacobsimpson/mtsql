package preprocessor

import (
	"testing"

	"github.com/jacobsimpson/mtsql/ast"
	md "github.com/jacobsimpson/mtsql/metadata"
	"github.com/stretchr/testify/assert"
)

func TestMapper(t *testing.T) {
	tests := []struct {
		name     string
		columns  []*md.Column
		input    *ast.Attribute
		expected []*md.Column
	}{
		{
			name:     "empty columns",
			columns:  []*md.Column{},
			input:    &ast.Attribute{Name: "abc"},
			expected: []*md.Column{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			m := newMapper(test.columns)
			result := m.findMatches(test.input)

			assert.Equal(test.expected, result)
		})
	}
}
