package physical

import "github.com/jacobsimpson/mtsql/metadata"

type PlanDescription struct {
	Name        string
	Description string
}

type RowReader interface {
	Columns() []*metadata.Column
	Read() ([]string, error)
	Reset() error
	Close()

	PlanDescription() *PlanDescription
	Children() []RowReader
}
