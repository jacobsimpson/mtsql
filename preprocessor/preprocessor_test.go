package preprocessor_test

import (
	"testing"

	"github.com/jacobsimpson/mtsql/ast"
	"github.com/jacobsimpson/mtsql/preprocessor"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	assert := assert.New(t)
	q := &ast.SFW{
		SelList: &ast.SelList{},
		From:    &ast.Relation{Name: "mock_table"},
		Condition: &ast.EqualCondition{
			LHS: &ast.Attribute{Name: "col1"},
			RHS: &ast.Constant{
				Type:  ast.StringType,
				Value: "raw_value",
				Raw:   "'raw_value'",
			},
		},
		OrderBy: &ast.OrderBy{},
	}

	err := preprocessor.Validate(q)

	assert.NotNil(err)
}
