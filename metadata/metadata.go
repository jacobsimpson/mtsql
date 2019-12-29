package metadata

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
	Name string
	Type ColumnType
}
