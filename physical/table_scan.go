package physical

import (
	"encoding/csv"
	"fmt"
	"os"
)

type tableScan struct {
	file      *os.File
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

func (t *tableScan) Close() {
	t.file.Close()
}

func (t *tableScan) Reset() error {
	t.Close()
	return t.init()
}

func (t *tableScan) init() error {
	f, err := os.Open(t.tableName)
	if err != nil {
		return err
	}
	t.file = f
	reader := csv.NewReader(f)
	columns, err := reader.Read()

	t.reader = reader
	t.columns = columns
	return nil
}

func (t *tableScan) PlanDescription() *PlanDescription {
	return &PlanDescription{
		Name:        "TableScan",
		Description: fmt.Sprintf("%s", t.tableName),
	}
}

func (t *tableScan) Children() []RowReader { return []RowReader{} }

func NewTableScan(tableName string) (RowReader, error) {
	ts := &tableScan{
		tableName: tableName,
	}
	if err := ts.init(); err != nil {
		return nil, err
	}
	return ts, nil
}
