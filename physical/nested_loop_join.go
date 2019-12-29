package physical

import (
	"io"

	"github.com/jacobsimpson/mtsql/metadata"
)

type nestedLoopJoin struct {
	left    RowReader
	right   RowReader
	leftRow []string
}

func (t *nestedLoopJoin) Columns() []*metadata.Column {
	result := []*metadata.Column{}
	for _, c := range t.left.Columns() {
		result = append(result, c)
	}
	for _, c := range t.right.Columns() {
		result = append(result, c)
	}
	return result
}

func (t *nestedLoopJoin) Read() ([]string, error) {
	if t.leftRow == nil {
		row, err := t.left.Read()
		if err != nil {
			return nil, err
		}
		t.leftRow = row
	}
	for {
		rightRow, err := t.right.Read()
		if err == io.EOF {
			row, err := t.left.Read()
			if err != nil {
				return nil, err
			}
			t.leftRow = row
			if err := t.right.Reset(); err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		} else {
			return append(t.leftRow, rightRow...), nil
		}
	}
}

func (t *nestedLoopJoin) Close() {}
func (t *nestedLoopJoin) Reset() error {
	if err := t.left.Reset(); err != nil {
		return err
	}
	return t.right.Reset()
}

func (t *nestedLoopJoin) PlanDescription() *PlanDescription {
	return &PlanDescription{
		Name:        "NestedLoopJoin",
		Description: "",
	}
}

func (t *nestedLoopJoin) Children() []RowReader { return []RowReader{} }

func NewNestedLoopJoin(left, right RowReader) (RowReader, error) {
	return &nestedLoopJoin{
		left:  left,
		right: right,
	}, nil
}
