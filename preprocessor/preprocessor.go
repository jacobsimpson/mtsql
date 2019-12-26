package preprocessor

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/jacobsimpson/mtsql/ast"
)

func Validate(q ast.Query) error {
	var sfw *ast.SFW
	if p, ok := q.(*ast.Profile); ok {
		sfw = p.SFW
	} else if s, ok := q.(*ast.SFW); ok {
		sfw = s
	} else {
		return fmt.Errorf("expected a select query, but got something else")
	}

	rel, ok := sfw.From.(*ast.Relation)
	if !ok {
		return fmt.Errorf("expected a relation in the FROM clause, but got something else")
	}
	filename := rel.Name + ".csv"
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("table %q could not be located at %q", rel.Name, filename)
	}
	reader := csv.NewReader(f)
	columns, err := reader.Read()
	if err != nil {
		return fmt.Errorf("unable to read columns for table %q at %q", rel.Name, filename)
	}
	columnsMap := map[string]bool{}
	for _, c := range columns {
		columnsMap[c] = true
	}

	if sfw.Condition != nil {
		eq, ok := sfw.Condition.(*ast.EqualCondition)
		if !ok {
			return fmt.Errorf("only = conditions are currently supported")
		}
		if !columnsMap[eq.LHS.Name] {
			return fmt.Errorf("no column %q in table %q", eq.LHS.Name, rel.Name)
		}
	}

	if !sfw.SelList.All {
		for _, a := range sfw.SelList.Attributes {
			if !columnsMap[a.Name] {
				return fmt.Errorf("no column %q in table %q", a.Name, rel.Name)
			}
		}
	}

	return nil
}
