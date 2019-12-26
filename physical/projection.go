package physical

type projection struct {
	rowReader RowReader
	columns   []int
}

func (t *projection) Columns() []string {
	c := t.rowReader.Columns()
	r := []string{}
	for _, col := range t.columns {
		r = append(r, c[col])
	}
	return r
}

func (t *projection) Read() ([]string, error) {
	row, err := t.rowReader.Read()
	if err != nil {
		return nil, err
	}
	r := []string{}
	for _, col := range t.columns {
		r = append(r, row[col])
	}
	return r, nil
}

func (t *projection) Close() {}

func NewProjection(rowReader RowReader, columns []string) (RowReader, error) {
	columnMap := map[string]int{}
	for i, c := range rowReader.Columns() {
		columnMap[c] = i
	}
	cols := []int{}
	for _, c := range columns {
		cols = append(cols, columnMap[c])
	}
	return &projection{
		rowReader: rowReader,
		columns:   cols,
	}, nil
}
