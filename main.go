package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jacobsimpson/csvsql/formatter"
	"github.com/jacobsimpson/csvsql/lexer"
	"github.com/jacobsimpson/csvsql/parser"
	"github.com/jacobsimpson/csvsql/physical"
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
	qp, err := physical.NewQueryPlan(q)
	if err != nil {
		return err
	}
	f := formatter.NewTableFormatter(qp)
	f.Print(os.Stdout)
	return nil
}
