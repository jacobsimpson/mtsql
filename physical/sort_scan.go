package physical

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type sortScan struct {
	rowReader     RowReader
	columnIndexes []int
	rows          [][]string
	next          int
}

func (t *sortScan) Columns() []string {
	return t.rowReader.Columns()
}

func (t *sortScan) Read() ([]string, error) {
	if t.next >= len(t.rows) {
		return nil, io.EOF
	}
	row := t.rows[t.next]
	t.next++
	return row, nil
}

func (t *sortScan) Close() {}

func (t *sortScan) PlanDescription() *PlanDescription {
	return &PlanDescription{
		Name:        "SortScan",
		Description: fmt.Sprintf("%v", strings.Join(t.Columns(), ", ")),
	}
}

func (t *sortScan) Children() []RowReader { return []RowReader{t.rowReader} }

// columnSorter joins has a slice of rows to be sorted.
type columnSorter struct {
	rows      [][]string
	columns   []int
	sortOrder []SortOrder
}

// Len is part of sort.Interface.
func (s *columnSorter) Len() int {
	return len(s.rows)
}

// Swap is part of sort.Interface.
func (s *columnSorter) Swap(i, j int) {
	s.rows[i], s.rows[j] = s.rows[j], s.rows[i]
}

// Less is part of sort.Interface.
func (s *columnSorter) Less(i, j int) bool {
	for n, c := range s.columns {
		if s.rows[i][c] < s.rows[j][c] {
			if s.sortOrder[n] == Asc {
				return true
			}
			return false
		}
		if s.rows[i][c] > s.rows[j][c] {
			if s.sortOrder[n] == Asc {
				return false
			}
			return true
		}
	}
	return false
}

func NewSortScan(rowReader RowReader, columns []SortScanCriteria) (RowReader, error) {
	columnMap := map[string]int{}
	for i, c := range rowReader.Columns() {
		columnMap[c] = i
	}
	cols := []int{}
	sortOrder := []SortOrder{}
	for _, c := range columns {
		cols = append(cols, columnMap[c.Column])
		sortOrder = append(sortOrder, c.SortOrder)
	}

	rows := [][]string{}
	for {
		row, err := rowReader.Read()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		if row == nil {
			break
		}
		rows = append(rows, row)
	}
	sort.Sort(&columnSorter{rows: rows, columns: cols, sortOrder: sortOrder})
	return &sortScan{
		rowReader:     rowReader,
		rows:          rows,
		columnIndexes: cols,
	}, nil
}
