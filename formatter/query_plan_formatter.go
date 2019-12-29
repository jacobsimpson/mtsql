package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/jacobsimpson/mtsql/physical"
)

type queryPlanFormatter struct {
	root physical.RowReader
}

func NewQueryPlanFormatter(rowReader physical.RowReader) Formatter {
	return &queryPlanFormatter{
		root: rowReader,
	}
}

func (f *queryPlanFormatter) Print(w io.Writer) {
	fmt.Fprintf(w, "\n")
	printPlanDescription(w, f.root, 0)
	fmt.Fprintf(w, "\n")
}

func printPlanDescription(w io.Writer, rowReader physical.RowReader, indentation int) {
	planDescription := rowReader.PlanDescription()
	if len(planDescription.Description) > 0 {
		fmt.Fprintf(w, "%so %s (%s)\n",
			strings.Repeat("| ", indentation),
			planDescription.Name,
			planDescription.Description)
	} else {
		fmt.Fprintf(w, "%so %s\n",
			strings.Repeat("| ", indentation),
			planDescription.Name)
	}
	children := rowReader.Children()

	// If there are multiple children, they need to be indented so there is
	// room for the line to continue down. The last child doesn't have to be
	// indented.
	indentationIncrement := 1
	for i, rr := range children {
		if i == len(children)-1 {
			indentationIncrement = 0
		}
		fmt.Fprintf(w, "%s|%s\n", strings.Repeat("| ", indentation), strings.Repeat("\\", indentationIncrement))
		printPlanDescription(w, rr, indentation+indentationIncrement)
	}
}
