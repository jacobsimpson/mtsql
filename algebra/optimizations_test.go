package algebra

import (
	"fmt"
	"testing"

	md "github.com/jacobsimpson/mtsql/metadata"
	"github.com/stretchr/testify/assert"
)

func TestPushDown(t *testing.T) {
	tests := []struct {
		name     string
		input    Operation
		expected Operation
	}{
		{
			name: "pushdown past union",
			input: &Projection{
				Child: &Selection{
					Child: &Union{
						LHS: &Source{},
						RHS: &Source{},
					},
				},
			},
			expected: &Projection{
				Child: &Union{
					LHS: &Selection{
						Child: &Source{},
					},
					RHS: &Selection{
						Child: &Source{},
					},
				},
			},
		},
		{
			name: "pushdown two steps",
			input: &Selection{
				requires: []*md.Column{{Qualifier: "tab1", Name: "name1"}},
				Child: &Projection{
					Child: &Union{
						LHS: &Source{},
						RHS: &Source{},
					},
				},
			},
			expected: &Projection{
				Child: &Union{
					LHS: &Selection{
						Child: &Source{},
					},
					RHS: &Selection{
						Child: &Source{},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			result := PushDownSelection(test.input)

			assert.Equal(test.expected, result)
		})
	}
}

func TestContainsAll(t *testing.T) {
	assert := assert.New(t)

	assert.True(containsAll(
		[]*md.Column{},
		[]*md.Column{}))
	assert.False(containsAll(
		[]*md.Column{},
		[]*md.Column{&md.Column{Name: "abc"}}))
	assert.True(containsAll(
		[]*md.Column{&md.Column{Name: "abc"}},
		[]*md.Column{&md.Column{Name: "abc"}}))
	assert.False(containsAll(
		[]*md.Column{&md.Column{Qualifier: "qual", Name: "abc"}},
		[]*md.Column{&md.Column{Name: "abc"}}))
}

func TestCanPushDownSelection(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name      string
		operation Operation
		selection *Selection
		expected  bool
	}{
		{
			name:      "push down over projection",
			operation: &Projection{},
			selection: &Selection{},
			expected:  true,
		},
		{
			name: "push down over missing columns",
			operation: &Projection{
				columns: []*md.Column{
					&md.Column{Name: "abc"},
				},
			},
			selection: &Selection{
				requires: []*md.Column{
					&md.Column{Name: "def"},
				},
			},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("provides = %+v\n", test.operation.Provides())
			fmt.Printf("requires = %+v\n", test.selection.Requires())
			assert.Equal(
				test.expected,
				canPushDownSelection(test.operation, test.selection))
		})
	}
}
