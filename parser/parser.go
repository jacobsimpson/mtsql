package parser

import (
	"fmt"
	"strings"

	"github.com/jacobsimpson/mtsql/ast"
	"github.com/jacobsimpson/mtsql/lexer"
)

func Parse(lex lexer.Lexer) (ast.Query, error) {
	return query(lex)
}

func query(lex lexer.Lexer) (ast.Query, error) {
	if !lex.Next() {
		return nil, nil
	}
	token := lex.Token()
	if token.Type == lexer.ErrorType {
		return nil, fmt.Errorf("could not tokenize input: %v", token.Raw)
	}
	if token.Type != lexer.IdentifierType && strings.ToUpper(token.Raw) != "SELECT" {
		return nil, fmt.Errorf("expected SELECT, found %q", token.Raw)
	}

	q := ast.SFW{}
	selList, err := selList(lex)
	if err != nil {
		return nil, err
	}
	q.SelList = selList

	if !lex.Next() {
		return nil, fmt.Errorf("expected FROM, found end of query")
	}
	token = lex.Token()
	if token.Type == lexer.ErrorType {
		return nil, fmt.Errorf("could not tokenize input: %v", token.Raw)
	}
	if token.Type != lexer.IdentifierType && strings.ToUpper(token.Raw) != "FROM" {
		return nil, fmt.Errorf("expected FROM, found %q", token.Raw)
	}
	from, err := from(lex)
	if err != nil {
		return nil, err
	}
	q.From = from

	if !lex.Next() {
		return &q, nil
	}
	token = lex.Token()
	if token.Type == lexer.ErrorType {
		return nil, fmt.Errorf("could not tokenize input: %v", token.Raw)
	}
	if token.Type == lexer.EOFType {
		return &q, nil
	}
	if token.Type != lexer.IdentifierType && strings.ToUpper(token.Raw) != "WHERE" {
		return nil, fmt.Errorf("expected WHERE, found %q", token.Raw)
	}
	condition, err := condition(lex)
	if err != nil {
		return nil, err
	}
	q.Condition = condition

	return &q, nil
}

func selList(lex lexer.Lexer) (*ast.SelList, error) {
	result := ast.SelList{}
	for {
		if !lex.Next() {
			return nil, fmt.Errorf("expected column, found nothing")
		}
		token := lex.Token()
		if token.Type == lexer.StarType {
			result.All = true
		} else if token.Type != lexer.IdentifierType {
			return nil, fmt.Errorf("expected column name, found %q", token.Raw)
		}
		if strings.ToUpper(token.Raw) == "FROM" {
			lex.UnreadToken()
			return &result, nil
		}
		result.Attributes = append(result.Attributes, &ast.Attribute{Name: token.Raw})

		if !lex.Next() {
			return nil, fmt.Errorf("unexpected end of query, no FROM clause specified")
		}
		token = lex.Token()
		if token.Type != lexer.CommaType {
			lex.UnreadToken()
			return &result, nil
		}
	}
}

func from(lex lexer.Lexer) (ast.From, error) {
	result := ast.Relation{}
	if !lex.Next() {
		return nil, fmt.Errorf("expected table, found nothing")
	}
	token := lex.Token()
	if token.Type != lexer.IdentifierType {
		return nil, fmt.Errorf("expected table name, found %q", token.Raw)
	}
	result.Name = token.Raw
	return &result, nil
}

func condition(lex lexer.Lexer) (ast.Condition, error) {
	result := ast.EqualCondition{}
	if !lex.Next() {
		return nil, fmt.Errorf("expected an attribute, found nothing")
	}
	token := lex.Token()
	if token.Type != lexer.IdentifierType {
		return nil, fmt.Errorf("expected attribute name, found %q", token.Raw)
	}
	result.LHS = &ast.Attribute{Name: token.Raw}

	if !lex.Next() {
		return nil, fmt.Errorf("expected =, found nothing")
	}
	token = lex.Token()
	if token.Type != lexer.EqualType {
		return nil, fmt.Errorf("expected =, found %q", token.Raw)
	}

	if !lex.Next() {
		return nil, fmt.Errorf("expected an attribute, found nothing")
	}
	token = lex.Token()
	if token.Type != lexer.StringType {
		return nil, fmt.Errorf("expected attribute name, found %q", token.Raw)
	}
	result.RHS = &ast.Attribute{Name: token.Raw}
	return &result, nil
}
