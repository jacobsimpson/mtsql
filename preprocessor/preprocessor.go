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

	columnsMap := map[string]*metadata.Column{}
	tablesMap := map[string]*metadata.Relation{}
	for _, tmd := range tableMetadata {
		tablesMap[tmd.Name] = tmd
		for k, v := range tmd.ColumnsMap() {
			columnsMap[k] = v
		}
	}
	if sfw.Where != nil {
		eq, ok := sfw.Where.(*ast.EqualCondition)
		if !ok {
			return fmt.Errorf("only = conditions are currently supported")
		}
		if columnsMap[eq.LHS.Name] == nil {
			return fmt.Errorf("no column %q in query", eq.LHS.Name)
		}
	}

	for _, a := range sfw.SelList.Attributes {
		if columnsMap[a.Name] == nil && a.Name != "*" {
			return fmt.Errorf("no column %q in query", a.Name)
		}
	}

	// Expand any '*' references to the appropriate column list.
	r := []*ast.Attribute{}
	for _, a := range sfw.SelList.Attributes {
		if a.Name == "*" {
			if a.Qualifier == "" {
				for _, tmd := range tableMetadata {
					for _, c := range tmd.Columns {
						r = append(r, &ast.Attribute{Name: c.Name})
					}
				}
				continue
			} else {
				tmd := tablesMap[a.Qualifier]
				if tmd == nil {
					return fmt.Errorf("table %q isn't in the query", a.Qualifier)
				}
				for _, c := range tmd.Columns {
					r = append(r, &ast.Attribute{Name: c.Name})
				}
			}
		}
		r = append(r, a)
	}
	sfw.SelList.Attributes = r

	return nil
}

func validateFrom(from ast.From, tables map[string]*metadata.Relation) ([]*metadata.Relation, error) {
	result := []*metadata.Relation{}
	for _, rel := range from.Tables() {
		if t := tables[rel.Name]; t != nil {
			result = append(result, t)
			continue
		}

		md := &metadata.Relation{
			Name:   rel.Name,
			Type:   metadata.CsvType,
			Source: rel.Name + ".csv",
		}
		columns, err := loadColumns(md.Name, md.Source)
		if err != nil {
			return nil, err
		}
		md.Columns = columns

		tables[rel.Name] = md
		result = append(result, md)
	}
	return result, nil
}

func loadColumns(tableName, file string) ([]*metadata.Column, error) {
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

	var columns []*metadata.Column
	for _, cn := range columnNames {
		columns = append(columns, &metadata.Column{
			Qualifier: tableName,
			Name:      cn,
			Type:      metadata.StringType,
		})
	}
	return columns, nil
}
