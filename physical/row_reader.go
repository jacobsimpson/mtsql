package physical

type PlanDescription struct {
	Name        string
	Description string
}

type RowReader interface {
	Columns() []string
	Read() ([]string, error)
	Close()

	PlanDescription() *PlanDescription
	Children() []RowReader
}
