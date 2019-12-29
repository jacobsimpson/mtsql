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
	fileName  string
	columns   []*metadata.Column
	firstRow  int64
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
	_, err := t.file.Seek(t.firstRow, 0)
	return err
}

func (t *tableScan) init() error {
	f, err := os.Open(t.fileName)
	if err != nil {
		return err
	}
	t.file = f
	reader := csv.NewReader(f)
	columns, err := reader.Read()
	if err != nil {
		return err
	}
	firstRow, err := f.Seek(0, 1)
	if err != nil {
		return err
	}
	t.firstRow = firstRow

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
		Description: fmt.Sprintf("%s, %s", t.tableName, t.fileName),
	}
}

func (t *tableScan) Children() []RowReader { return []RowReader{} }

func NewTableScan(tableName, fileName string) (RowReader, error) {
	ts := &tableScan{
		tableName: tableName,
		fileName:  fileName,
	}
	if err := ts.init(); err != nil {
		return nil, err
	}
	return ts, nil
}
