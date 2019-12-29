package physical

import (
	"fmt"

	"github.com/jacobsimpson/mtsql/ast"
	"github.com/jacobsimpson/mtsql/metadata"
)

type SortOrder string

const (
	Asc  SortOrder = "Asc"
	Desc SortOrder = "Desc"
)

type SortScanCriteria struct {
	Column    *metadata.Column
	SortOrder SortOrder
}

func NewQueryPlan(q ast.Query) (RowReader, error) {
	var sfw *ast.SFW
	if p, ok := q.(*ast.Profile); ok {
		sfw = p.SFW
	} else if s, ok := q.(*ast.SFW); ok {
		sfw = s
	} else {
		return nil, fmt.Errorf("expected a select query, but got something else")
	}

	var rowReader RowReader
	if rel, ok := sfw.From.(*ast.Relation); ok {
		rr, err := NewTableScan(rel.Name, rel.Name+".csv")
		if err != nil {
			return nil, err
		}
		rowReader = rr
	} else if ij, ok := sfw.From.(*ast.InnerJoin); ok {
		left, err := NewTableScan(ij.Left.Name, ij.Left.Name+".csv")
		if err != nil {
			return nil, err
		}
		right, err := NewTableScan(ij.Right.Name, ij.Right.Name+".csv")
		if err != nil {
			return nil, err
		}
		rowReader, err = NewNestedLoopJoin(left, right)
		if err != nil {
			return nil, err
		}
		// , ij.On
	} else {
		return nil, fmt.Errorf("expected a relation in the FROM clause, but got something else")
	}

	if sfw.OrderBy != nil {
		columns := []SortScanCriteria{}
		for _, c := range sfw.OrderBy.Criteria {
			sc := SortScanCriteria{
				Column: &metadata.Column{
					Qualifier: c.Attribute.Qualifier,
					Name:      c.Attribute.Name,
				},
				SortOrder: Asc,
			}
			if c.SortOrder == ast.Desc {
				sc.SortOrder = Desc
			}
			columns = append(columns, sc)
		}
		rr, err := NewSortScan(rowReader, columns)
		if err != nil {
			return nil, err
		}
		rowReader = rr
	}

	if sfw.Where != nil {
		eq, ok := sfw.Where.(*ast.EqualCondition)
		if !ok {
			return nil, fmt.Errorf("only = conditions are currently supported")
		}
		rr, err := NewFilter(rowReader,
			&metadata.Column{Qualifier: eq.LHS.Qualifier, Name: eq.LHS.Name},
			eq.RHS)
		if err != nil {
			return nil, err
		}
		rowReader = rr
	}

	columns := []*metadata.Column{}
	for _, a := range sfw.SelList.Attributes {
		columns = append(columns, &metadata.Column{
			Qualifier: a.Qualifier,
			Name:      a.Name,
		})
	}
	rr, err := NewProjection(rowReader, columns)
	if err != nil {
		return nil, err
	}
	rowReader = rr

	return rowReader, nil
}
