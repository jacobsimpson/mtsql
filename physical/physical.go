package physical

import (
	"fmt"

	"github.com/jacobsimpson/mtsql/ast"
)

func NewQueryPlan(q ast.Query) (RowReader, error) {
	var sfw *ast.SFW
	if p, ok := q.(*ast.Profile); ok {
		sfw = p.SFW
	} else if s, ok := q.(*ast.SFW); ok {
		sfw = s
	} else {
		return nil, fmt.Errorf("expected a select query, but got something else")
	}

	rel, ok := sfw.From.(*ast.Relation)
	if !ok {
		return nil, fmt.Errorf("expected a relation in the FROM clause, but got something else")
	}
	rowReader, err := NewTableScan(rel.Name + ".csv")
	if err != nil {
		return nil, err
	}

	if sfw.Condition != nil {
		eq, ok := sfw.Condition.(*ast.EqualCondition)
		if !ok {
			return nil, fmt.Errorf("only = conditions are currently supported")
		}
		rowReader, err = NewFilter(rowReader, eq.LHS.Name, eq.RHS)
		if err != nil {
			return nil, err
		}
	}

	if !sfw.SelList.All {
		columns := []string{}
		for _, a := range sfw.SelList.Attributes {
			columns = append(columns, a.Name)
		}
		rowReader, err = NewProjection(rowReader, columns)
		if err != nil {
			return nil, err
		}
	}

	return rowReader, nil
}
