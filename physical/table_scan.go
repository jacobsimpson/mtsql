package physical

import (
	"encoding/csv"
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
