package physical

type PlanDescription struct {
	Name        string
	Description string
}

type RowReader interface {
	Columns() []string
	Read() ([]string, error)
	Reset() error
	Close()

	PlanDescription() *PlanDescription
	Children() []RowReader
}
