package preprocessor

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/jacobsimpson/mtsql/algebra"
	"github.com/jacobsimpson/mtsql/ast"
	md "github.com/jacobsimpson/mtsql/metadata"
)

func Convert(q ast.Query, tables map[string]*md.Relation) (algebra.Operation, error) {
	var sfw *ast.SFW
	if p, ok := q.(*ast.Profile); ok {
		sfw = p.SFW
	} else if s, ok := q.(*ast.SFW); ok {
		sfw = s
	} else {
		return nil, fmt.Errorf("expected a select query, but got something else")
	}

	result, err := convertFrom(sfw.From, tables)
	if err != nil {
		return nil, err
	}

	if sfw.Where != nil {
		result = &algebra.Selection{
			Child: result,
		}
	}

	if sfw.SelList != nil {
		mapper := newMapper(result.Provides())

		columns := []*md.Column{}
		for _, a := range sfw.SelList.Attributes {
			matches := mapper.findMatches(a)
			if len(matches) == 0 {
				return nil, fmt.Errorf("no columns for %s", a)
			}
			if len(matches) > 1 {
				return nil, fmt.Errorf("multiple column matches %+v", matches)
			}
			fmt.Printf("the matching column = %+v\n", matches[0])
			columns = append(columns, matches[0])
		}

		result = algebra.NewProjection(
			result,
			columns)
	}

	return result, nil
}

func convertFrom(from ast.From, tables map[string]*md.Relation) (algebra.Operation, error) {
	if rel, ok := from.(*ast.Relation); ok {
		return convertRelation(rel, tables)
	}
	if ij, ok := from.(*ast.InnerJoin); ok {
		left, err := convertRelation(ij.Left, tables)
		if err != nil {
			return nil, err
		}
		right, err := convertRelation(ij.Right, tables)
		if err != nil {
			return nil, err
		}
		selection := &algebra.Selection{
			Child: &algebra.Product{
				LHS: left,
				RHS: right,
			},
			//requires: []*md.Column{
			//	&md.Column{},
			//},
		}
		return selection, nil
	}
	return nil, fmt.Errorf("unable to convert from relationship")
}

func convertRelation(relation *ast.Relation, tables map[string]*md.Relation) (*algebra.Source, error) {
	t := tables[relation.Name]
	if t == nil {
		t := &md.Relation{
			Name:   relation.Name,
			Type:   md.CsvType,
			Source: relation.Name + ".csv",
		}
		columns, err := loadColumns(t.Name, t.Source)
		if err != nil {
			return nil, err
		}
		t.Columns = columns

		tables[t.Name] = t
	}
	return algebra.NewSource(t.Name, t.Columns), nil
}

func loadColumns(tableName, file string) ([]*md.Column, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("table %q could not be located at %q", tableName, file)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	columnNames, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("unable to read columns for table %q at %q", tableName, file)
	}

	var columns []*md.Column
	for _, cn := range columnNames {
		columns = append(columns, &md.Column{
			Qualifier: tableName,
			Name:      cn,
			Type:      md.StringType,
		})
	}
	return columns, nil
}
