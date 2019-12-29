package formatter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jacobsimpson/mtsql/physical"
	"golang.org/x/crypto/ssh/terminal"
)

type Formatter interface {
	Print(io.Writer)
}

type tableFormatter struct {
	rowReader physical.RowReader
	width     int
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

func (f *tableFormatter) Print(w io.Writer) {
	for _, name := range f.rowReader.Columns() {
		fmt.Fprintf(w, "%10s", name.QualifiedName())
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

type queryPlanFormatter struct {
	root physical.RowReader
}

func (f *queryPlanFormatter) Print(w io.Writer) {
	fmt.Fprintf(w, "\n")
	printPlanDescription(w, f.root, 0)
	fmt.Fprintf(w, "\n")
}

func printPlanDescription(w io.Writer, rowReader physical.RowReader, indentation int) {
	planDescription := rowReader.PlanDescription()
	fmt.Fprintf(w, "%so %s (%s)\n",
		strings.Repeat("| ", indentation),
		planDescription.Name,
		planDescription.Description)
	children := rowReader.Children()
	indentationIncrement := 0
	if len(children) > 1 {
		indentationIncrement = 1
	}
	for _, rr := range children {
		fmt.Fprintf(w, "%s|%s\n", strings.Repeat("| ", indentation), strings.Repeat("\\", indentationIncrement))
		printPlanDescription(w, rr, indentation+indentationIncrement)
	}
}

func NewQueryPlanFormatter(rowReader physical.RowReader) Formatter {
	return &queryPlanFormatter{
		root: rowReader,
	}
}
