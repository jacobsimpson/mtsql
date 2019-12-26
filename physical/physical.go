package physical

import (
	"fmt"

	"github.com/jacobsimpson/csvsql/ast"
)

func NewQueryPlan(q ast.Query) (RowReader, error) {
	swf, ok := q.(*ast.SFW)
	if !ok {
		return nil, fmt.Errorf("expected a select query, but got something else")
	}

	rel, ok := swf.From.(*ast.Relation)
	if !ok {
		return nil, fmt.Errorf("expected a select query, but got something else")
	}
	rowReader, err := NewTableScan(rel.Name + ".csv")
	if err != nil {
		return nil, err
	}

	//assert.NotNil(swf.SelList)
	//assert.Equal(len(swf.SelList.Attributes), 1)
	//assert.Equal(swf.SelList.Attributes[0].Name, "col1")

	//assert.NotNil(swf.Condition)
	//eq, ok := swf.Condition.(*ast.EqualCondition)
	//assert.True(ok)
	//assert.Equal(eq.LHS.Name, "col1")
	//assert.Equal(eq.RHS.Name, "'abcd'")
	return rowReader, nil
}
