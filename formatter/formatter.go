package formatter

import (
	"fmt"
	"io"
	"os"

	"github.com/jacobsimpson/csvsql/physical"
	"golang.org/x/crypto/ssh/terminal"
)

type Formatter interface {
	Print(io.Writer)
}

type tableFormatter struct {
	rowReader physical.RowReader
	width     int
}

func (f *tableFormatter) Print(w io.Writer) {
	for _, name := range f.rowReader.Columns() {
		fmt.Fprintf(w, "%10s", name)
	}
	fmt.Fprintf(w, "\n")
	for {
		row, err := f.rowReader.Read()
		if err != nil {
			fmt.Fprintf(w, "\n\n%+v\n\n", err)
			return
		}
		for _, cell := range row {
			fmt.Fprintf(w, "%10s", cell)
		}
		fmt.Fprintf(w, "\n")
	}
}

func NewTableFormatter(rowReader physical.RowReader) Formatter {
	width, _, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width = 80
	}

	return &tableFormatter{
		rowReader: rowReader,
		width:     width,
	}
}
