package parser

import (
	"fmt"
	"strconv"
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
	if token.Type == lexer.IdentifierType {
		if strings.ToUpper(token.Raw) == "SELECT" {
			lex.UnreadToken()
			return sfw(lex)
		} else if strings.ToUpper(token.Raw) == "PROFILE" {
			lex.UnreadToken()
			return profile(lex)
		}
	}

	return nil, fmt.Errorf("expected SELECT, found %q", token.Raw)
}

func profile(lex lexer.Lexer) (*ast.Profile, error) {
	if !lex.Next() {
		return nil, nil
	}
	token := lex.Token()
	if token.Type != lexer.IdentifierType || strings.ToUpper(token.Raw) != "PROFILE" {
		return nil, fmt.Errorf("expected PROFILE, found %q", token.Raw)
	}
	sfw, err := sfw(lex)
	if err != nil {
		return nil, err
	}
	return &ast.Profile{
		SFW: sfw,
	}, nil
}

func sfw(lex lexer.Lexer) (*ast.SFW, error) {
	if !lex.Next() {
		return nil, nil
	}
	token := lex.Token()
	if token.Type != lexer.IdentifierType || strings.ToUpper(token.Raw) != "SELECT" {
		return nil, fmt.Errorf("expected SELECT, found %q", token.Raw)
	}

	q := ast.SFW{}
	selList, err := selList(lex)
	if err != nil {
		return nil, err
	}
	q.SelList = selList

	from, err := from(lex)
	if err != nil {
		return nil, err
	}
	q.From = from

	condition, err := condition(lex)
	if err != nil {
		return nil, err
	}
	q.Condition = condition

	orderby, err := orderBy(lex)
	if err != nil {
		return nil, err
	}
	q.OrderBy = orderby

	if !lex.Next() {
		return &q, nil
	}
	token = lex.Token()
	if token.Type == lexer.ErrorType {
		return nil, fmt.Errorf("could not tokenize input: %v", token.Raw)
	}
	if token.Type != lexer.EOFType {
		return nil, fmt.Errorf("extra stuff left over: %v", token.Raw)
	}
	return &q, nil
}

func selList(lex lexer.Lexer) (*ast.SelList, error) {
	result := ast.SelList{}
	for {
		attribute, err := attribute(lex)
		if err != nil {
			return nil, err
		}
		if attribute == nil {
			return &result, nil
		}
		result.Attributes = append(result.Attributes, attribute)

		if !lex.Next() {
			return nil, fmt.Errorf("unexpected end of query, no FROM clause specified")
		}
		token := lex.Token()
		if token.Type != lexer.CommaType {
			lex.UnreadToken()
			return &result, nil
		}
	}
}

func attribute(lex lexer.Lexer) (*ast.Attribute, error) {
	if !lex.Next() {
		return nil, fmt.Errorf("expected column, found nothing")
	}
	token := lex.Token()
	if token.Type == lexer.StarType {
		return &ast.Attribute{Name: token.Raw}, nil
	} else if token.Type != lexer.IdentifierType {
		return nil, fmt.Errorf("expected column name, found %q", token.Raw)
	}
	if strings.ToUpper(token.Raw) == "FROM" {
		lex.UnreadToken()
		return nil, nil
	}

	attribute := &ast.Attribute{Name: token.Raw}
	if !lex.Next() {
		return attribute, nil
	}
	token = lex.Token()
	if token.Type != lexer.PeriodType {
		lex.UnreadToken()
		return attribute, nil
	}
	attribute.Qualifier = attribute.Name

	if !lex.Next() {
		return nil, fmt.Errorf("partially specified column '%s.'", attribute.Qualifier)
	}
	token = lex.Token()
	if token.Type != lexer.IdentifierType {
		return nil, fmt.Errorf("expected identifier as part of qualified column name '%s.', found %q", attribute.Qualifier, token.Raw)
	}
	attribute.Name = token.Raw
	return attribute, nil
}

func from(lex lexer.Lexer) (ast.From, error) {
	if !lex.Next() {
		return nil, fmt.Errorf("expected FROM, found end of query")
	}
	token := lex.Token()
	if token.Type == lexer.ErrorType {
		return nil, fmt.Errorf("could not tokenize input: %v", token.Raw)
	}
	if token.Type != lexer.IdentifierType && strings.ToUpper(token.Raw) != "FROM" {
		return nil, fmt.Errorf("expected FROM, found %q", token.Raw)
	}
	result := ast.Relation{}
	if !lex.Next() {
		return nil, fmt.Errorf("expected table, found nothing")
	}
	token = lex.Token()
	if token.Type != lexer.IdentifierType {
		return nil, fmt.Errorf("expected table name, found %q", token.Raw)
	}
	result.Name = token.Raw
	return &result, nil
}

func condition(lex lexer.Lexer) (ast.Condition, error) {
	if !lex.Next() {
		return nil, nil
	}
	token := lex.Token()
	if token.Type == lexer.ErrorType {
		return nil, fmt.Errorf("could not tokenize input: %v", token.Raw)
	}
	if token.Type == lexer.EOFType {
		return nil, nil
	}
	if token.Type != lexer.IdentifierType && strings.ToUpper(token.Raw) != "WHERE" {
		lex.UnreadToken()
		return nil, nil
	}
	result := ast.EqualCondition{}
	if !lex.Next() {
		return nil, fmt.Errorf("expected an attribute, found nothing")
	}
	token = lex.Token()
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
	switch token.Type {
	case lexer.StringType:
		result.RHS = &ast.Constant{
			Type:  ast.StringType,
			Value: token.Raw[1 : len(token.Raw)-1],
			Raw:   token.Raw,
		}
	case lexer.IntegerType:
		i, err := strconv.Atoi(token.Raw)
		if err != nil {
			return nil, fmt.Errorf("unable to convert constant %q to integer", token.Raw)
		}
		result.RHS = &ast.Constant{
			Type:  ast.IntegerType,
			Value: i,
			Raw:   token.Raw,
		}
	default:
		return nil, fmt.Errorf("unexpected token type %s for %q", token.Type, token.Raw)
	}
	return &result, nil
}

func orderBy(lex lexer.Lexer) (*ast.OrderBy, error) {
	if !lex.Next() {
		return nil, nil
	}
	token := lex.Token()
	if token.Type == lexer.ErrorType {
		return nil, fmt.Errorf("could not tokenize input: %v", token.Raw)
	}
	if token.Type == lexer.EOFType {
		return nil, nil
	}
	if token.Type != lexer.IdentifierType && strings.ToUpper(token.Raw) != "ORDER" {
		lex.UnreadToken()
		return nil, nil
	}
	if !lex.Next() {
		return nil, fmt.Errorf("expected BY, found nothing")
	}
	token = lex.Token()
	if token.Type == lexer.ErrorType {
		return nil, fmt.Errorf("could not tokenize input: %v", token.Raw)
	}
	if token.Type == lexer.EOFType {
		return nil, fmt.Errorf("expected BY, found nothing")
	}
	if token.Type != lexer.IdentifierType && strings.ToUpper(token.Raw) != "BY" {
		return nil, fmt.Errorf("expected BY, found %q", token.Raw)

	}

	result := ast.OrderBy{}
	if !lex.Next() {
		return nil, fmt.Errorf("expected an attribute, found nothing")
	}
	token = lex.Token()
	if token.Type != lexer.IdentifierType {
		return nil, fmt.Errorf("expected attribute name, found %q", token.Raw)
	}
	oc := ast.OrderCriteria{
		Attribute: &ast.Attribute{Name: token.Raw},
		SortOrder: ast.Asc,
	}
	if lex.Next() {
		token = lex.Token()
		if token.Type == lexer.IdentifierType && strings.ToUpper(token.Raw) == "ASC" {
			oc.SortOrder = ast.Asc
		} else if token.Type == lexer.IdentifierType && strings.ToUpper(token.Raw) == "DESC" {
			oc.SortOrder = ast.Desc
		} else {
			lex.UnreadToken()
		}
	}
	result.Criteria = append(result.Criteria, &oc)

	for {
		if !lex.Next() {
			break
		}
		token = lex.Token()
		if token.Type == lexer.EOFType {
			break
		}
		if token.Type != lexer.CommaType {
			return nil, fmt.Errorf("expected ',', found %q", token.Raw)
		}

		if !lex.Next() {
			return nil, fmt.Errorf("expected an attribute, found nothing")
		}
		token = lex.Token()
		if token.Type != lexer.IdentifierType {
			return nil, fmt.Errorf("expected attribute name, found %q", token.Raw)
		}
		oc := ast.OrderCriteria{
			Attribute: &ast.Attribute{Name: token.Raw},
			SortOrder: ast.Asc,
		}
		if lex.Next() {
			token = lex.Token()
			if token.Type == lexer.IdentifierType && strings.ToUpper(token.Raw) == "ASC" {
				oc.SortOrder = ast.Asc
			} else if token.Type == lexer.IdentifierType && strings.ToUpper(token.Raw) == "DESC" {
				oc.SortOrder = ast.Desc
			} else {
				lex.UnreadToken()
			}
		}
		result.Criteria = append(result.Criteria, &oc)
	}

	return &result, nil
}
