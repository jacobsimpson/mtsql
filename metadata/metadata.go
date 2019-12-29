package metadata

import "fmt"

type RelationType string

const (
	CsvType RelationType = "csv"
)

type Relation struct {
	Name    string
	Type    RelationType
	Source  string
	Columns []*Column
}

func (r *Relation) ColumnsMap() map[string]*Column {
	columnsMap := map[string]*Column{}
	for _, c := range r.Columns {
		columnsMap[c.Name] = c
	}
	return columnsMap
}

type ColumnType string

const (
	StringType ColumnType = "string"
)

type Column struct {
	Qualifier string
	Name      string
	Alias     string
	Type      ColumnType
}

func (c *Column) QualifiedName() string { return fmt.Sprintf("%s.%s", c.Qualifier, c.Name) }

func (c *Column) String() string {
	return fmt.Sprintf("{Qualifier: %q, Name: %q, Alias: %q, Type: %s}",
		c.Qualifier,
		c.Name,
		c.Alias,
		c.Type)
}
