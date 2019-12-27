package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jacobsimpson/mtsql/ast"
	"github.com/jacobsimpson/mtsql/formatter"
	"github.com/jacobsimpson/mtsql/lexer"
	"github.com/jacobsimpson/mtsql/metadata"
	"github.com/jacobsimpson/mtsql/parser"
	"github.com/jacobsimpson/mtsql/physical"
	"github.com/jacobsimpson/mtsql/preprocessor"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute query: %+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "%s <SQL query>\n", filepath.Base(os.Args[0]))
		return nil
	}

	query := os.Args[1]
	q, err := parser.Parse(lexer.NewFilterWhitespace(strings.NewReader(query)))
	if err != nil {
		return err
	}

	if err := preprocessor.Validate(q, map[string]*metadata.Relation{}); err != nil {
		return err
	}

	qp, err := physical.NewQueryPlan(q)
	if err != nil {
		return err
	}
	if _, ok := q.(*ast.Profile); ok {
		fmt.Println("Showing the query plan...")
		f := formatter.NewQueryPlanFormatter(qp)
		f.Print(os.Stdout)
	} else {
		f := formatter.NewTableFormatter(qp)
		f.Print(os.Stdout)
	}
	return nil
}
