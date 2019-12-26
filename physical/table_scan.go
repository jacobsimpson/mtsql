package physical

import (
	"encoding/csv"
	"fmt"
	"os"
)

type tableScan struct {
	reader    *csv.Reader
	tableName string
	columns   []string
}

func (t *tableScan) Columns() []string {
	return t.columns
}

func (t *tableScan) Read() ([]string, error) {
	return t.reader.Read()
}

func (t *tableScan) Close() {}

func (t *tableScan) PlanDescription() *PlanDescription {
	return &PlanDescription{
		Name:        "TableScan",
		Description: fmt.Sprintf("%s", t.tableName),
	}
}

func (t *tableScan) Children() []RowReader { return []RowReader{} }

func NewTableScan(tableName string) (RowReader, error) {
	f, err := os.Open(tableName)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(f)
	columns, err := reader.Read()
	return &tableScan{
		tableName: tableName,
		reader:    reader,
		columns:   columns,
	}, nil
}
