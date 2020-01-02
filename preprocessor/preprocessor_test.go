package preprocessor

import (
	"fmt"
	"testing"

	"github.com/jacobsimpson/mtsql/algebra"
	"github.com/jacobsimpson/mtsql/ast"
	md "github.com/jacobsimpson/mtsql/metadata"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name     string
		query    ast.Query
		tables   map[string]*md.Relation
		expected algebra.Operation
		err      error
	}{
		{
			name:  "empty query",
			query: &ast.SFW{},
			err:   fmt.Errorf("unable to convert from relationship"),
		},
		{
			name: "only from clause",
			query: &ast.SFW{
				From: &ast.Relation{Name: "this"},
			},
			tables: map[string]*md.Relation{
				"this": &md.Relation{Name: "this"},
			},
			expected: &algebra.Source{Name: "this"},
		},
		{
			name: "only select clause",
			query: &ast.SFW{
				SelList: &ast.SelList{
					Attributes: []*ast.Attribute{
						{Name: "name"},
					},
				},
				From: &ast.Relation{Name: "this"},
			},
			tables: map[string]*md.Relation{
				"this": &md.Relation{
					Name:   "this",
					Type:   md.CsvType,
					Source: "this",
					Columns: []*md.Column{
						{Name: "name", Type: md.StringType},
					},
				},
			},
			expected: algebra.NewProjection(
				algebra.NewSource("this", []*md.Column{
					{Name: "name", Type: md.StringType},
				}),
				[]*md.Column{
					{Name: "name", Type: md.StringType},
				},
			),
		},
		{
			name: "inner join",
			query: &ast.SFW{
				SelList: &ast.SelList{
					Attributes: []*ast.Attribute{
						{Name: "name"},
					},
				},
				From: &ast.InnerJoin{
					Left:  &ast.Relation{Name: "this"},
					Right: &ast.Relation{Name: "that"},
					On: &ast.EqualColumnCondition{
						Left:  &ast.Attribute{Qualifier: "this", Name: "id"},
						Right: &ast.Attribute{Qualifier: "that", Name: "id"},
					},
				},
			},
			tables: map[string]*md.Relation{
				"this": &md.Relation{
					Name:   "this",
					Type:   md.CsvType,
					Source: "this",
					Columns: []*md.Column{
						{Name: "id", Type: md.StringType},
						{Name: "name", Type: md.StringType},
					},
				},
				"that": &md.Relation{
					Name:   "that",
					Type:   md.CsvType,
					Source: "that",
					Columns: []*md.Column{
						{Name: "id", Type: md.StringType},
					},
				},
			},
			expected: algebra.NewProjection(
				algebra.NewSelection(
					&algebra.Product{
						LHS: algebra.NewSource("this", []*md.Column{
							{Name: "id", Type: md.StringType},
							{Name: "name", Type: md.StringType},
						}),
						RHS: algebra.NewSource("that", []*md.Column{
							{Name: "id", Type: md.StringType},
						}),
					},
					nil,
				),
				[]*md.Column{
					{Name: "name", Type: md.StringType},
				},
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			op, err := Convert(test.query, test.tables)

			assert.Equal(test.err, err)
			assert.Equal(test.expected, op)
		})
	}
}
