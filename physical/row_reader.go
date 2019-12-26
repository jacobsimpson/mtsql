package physical

type RowReader interface {
	Columns() []string
	Read() ([]string, error)
	Close()
}
