package preprocessor

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/jacobsimpson/mtsql/ast"
	"github.com/jacobsimpson/mtsql/metadata"
)

func Validate(q ast.Query, tables map[string]*metadata.Relation) error {
	var sfw *ast.SFW
	if p, ok := q.(*ast.Profile); ok {
		sfw = p.SFW
	} else if s, ok := q.(*ast.SFW); ok {
		sfw = s
	} else {
		return fmt.Errorf("expected a select query, but got something else")
	}

	tableMetadata, err := validateFrom(sfw.From, tables)
	if err != nil {
		return err
	}

	columnsMap := tableMetadata.ColumnsMap()
	if sfw.Condition != nil {
		eq, ok := sfw.Condition.(*ast.EqualCondition)
		if !ok {
			return fmt.Errorf("only = conditions are currently supported")
		}
		if columnsMap[eq.LHS.Name] == nil {
			return fmt.Errorf("no column %q in table %q", eq.LHS.Name, tableMetadata.Name)
		}
	}

	for _, a := range sfw.SelList.Attributes {
		if columnsMap[a.Name] == nil && a.Name != "*" {
			return fmt.Errorf("no column %q in table %q", a.Name, tableMetadata.Name)
		}
	}

	// Expand any '*' references to the appropriate column list.
	r := []*ast.Attribute{}
	for _, a := range sfw.SelList.Attributes {
		if a.Name == "*" {
			for _, c := range tableMetadata.Columns {
				r = append(r, &ast.Attribute{Name: c.Name})
			}
			continue
		}
		r = append(r, a)
	}
	sfw.SelList.Attributes = r

	return nil
}

func validateFrom(from ast.From, tables map[string]*metadata.Relation) (*metadata.Relation, error) {
	rel, ok := from.(*ast.Relation)
	if !ok {
		return nil, fmt.Errorf("expected a relation in the FROM clause, but got something else")
	}

	if t := tables[rel.Name]; t != nil {
		return t, nil
	}

	md := &metadata.Relation{
		Name:   rel.Name,
		Type:   metadata.CsvType,
		Source: rel.Name + ".csv",
	}
	f, err := os.Open(md.Source)
	if err != nil {
		return nil, fmt.Errorf("table %q could not be located at %q", rel.Name, md.Source)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	columnNames, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("unable to read columns for table %q at %q", rel.Name, md.Source)
	}

	for _, cn := range columnNames {
		md.Columns = append(md.Columns, &metadata.Column{Name: cn, Type: metadata.AnyType})
	}

	tables[rel.Name] = md
	return md, nil
}
