package physical

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/jacobsimpson/mtsql/metadata"
)

type tableScan struct {
	file      *os.File
	reader    *csv.Reader
	tableName string
	columns   []*metadata.Column
}

func (t *tableScan) Columns() []*metadata.Column {
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
	for _, c := range columns {
		t.columns = append(t.columns, &metadata.Column{
			Qualifier: t.tableName,
			Name:      c,
		})
	}
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
