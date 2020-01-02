package algebra

import (
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
					columns: []*md.Column{{Qualifier: "tab1", Name: "name1"}},
					Child: &Union{
						LHS: &Source{
							provides: []*md.Column{{Qualifier: "tab1", Name: "name1"}},
						},
						RHS: &Source{
							provides: []*md.Column{{Qualifier: "tab1", Name: "name1"}},
						},
					},
				},
			},
			expected: &Projection{
				columns: []*md.Column{{Qualifier: "tab1", Name: "name1"}},
				Child: &Union{
					LHS: &Selection{
						requires: []*md.Column{{Qualifier: "tab1", Name: "name1"}},
						Child: &Source{
							provides: []*md.Column{{Qualifier: "tab1", Name: "name1"}},
						},
					},
					RHS: &Selection{
						requires: []*md.Column{{Qualifier: "tab1", Name: "name1"}},
						Child: &Source{
							provides: []*md.Column{{Qualifier: "tab1", Name: "name1"}},
						},
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
		[]*md.Column{{Name: "abc"}}))
	assert.True(containsAll(
		[]*md.Column{{Name: "abc"}},
		[]*md.Column{{Name: "abc"}}))
	assert.False(containsAll(
		[]*md.Column{{Qualifier: "qual", Name: "abc"}},
		[]*md.Column{{Name: "abc"}}))
}

func TestCanPushDownSelection(t *testing.T) {
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
			name: "attempt push down over missing columns",
			operation: &Projection{
				columns: []*md.Column{
					{Name: "abc"},
				},
			},
			selection: &Selection{
				requires: []*md.Column{
					{Name: "def"},
				},
			},
			expected: false,
		},
		{
			name: "pushdown two steps",
			operation: &Projection{
				columns: []*md.Column{{Qualifier: "tab1", Name: "name1"}},
				Child: &Union{
					LHS: &Source{
						provides: []*md.Column{{Qualifier: "tab1", Name: "name1"}},
					},
					RHS: &Source{
						provides: []*md.Column{{Qualifier: "tab1", Name: "name1"}},
					},
				},
			},
			selection: &Selection{
				requires: []*md.Column{{Qualifier: "tab1", Name: "name1"}},
				Child: &Projection{
					columns: []*md.Column{{Qualifier: "tab1", Name: "name1"}},
					Child: &Union{
						LHS: &Source{
							provides: []*md.Column{{Qualifier: "tab1", Name: "name1"}},
						},
						RHS: &Source{
							provides: []*md.Column{{Qualifier: "tab1", Name: "name1"}},
						},
					},
				},
			},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(
				test.expected,
				canPushDownSelection(test.operation, test.selection))
		})
	}
}
