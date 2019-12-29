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
	if s, err := sfw(lex); err != nil {
		return nil, err
	} else if s != nil {
		return s, nil
	}

	if p, err := profile(lex); err != nil {
		return nil, err
	} else if p != nil {
		return p, nil
	}

	return nil, fmt.Errorf("expected SELECT or PROFILE")
}

func profile(lex lexer.Lexer) (*ast.Profile, error) {
	if ok, err := ifKeywords(lex, "PROFILE"); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
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
	if ok, err := ifKeywords(lex, "SELECT"); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
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

	where, err := where(lex)
	if err != nil {
		return nil, err
	}
	q.Where = where

	orderby, err := orderBy(lex)
	if err != nil {
		return nil, err
	}
	q.OrderBy = orderby

	if !lex.Next() {
		return &q, nil
	}
	token := lex.Token()
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
	if token.Type == lexer.PeriodType {
		attribute.Qualifier = attribute.Name

		if !lex.Next() {
			return nil, fmt.Errorf("partially specified column '%s.'", attribute.Qualifier)
		}
		token = lex.Token()
		if token.Type != lexer.IdentifierType {
			return nil, fmt.Errorf("expected identifier as part of qualified column name '%s.', found %q", attribute.Qualifier, token.Raw)
		}
		attribute.Name = token.Raw
		if !lex.Next() {
			return attribute, nil
		}
		token = lex.Token()
	}

	if token.Type == lexer.CommaType || token.Type == lexer.EOFType ||
		(token.Type == lexer.IdentifierType && strings.ToUpper(token.Raw) == "FROM") {
		lex.UnreadToken()
		return attribute, nil
	}
	if token.Type != lexer.IdentifierType {
		return nil, fmt.Errorf("expected alias for column name '%s.', found %q", attribute.Name, token.Raw)
	}
	attribute.Alias = token.Raw
	return attribute, nil
}

func from(lex lexer.Lexer) (ast.From, error) {
	if ok, err := ifKeywords(lex, "FROM"); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("expected FROM clause")
	}

	tableName, err := tableName(lex)
	if err != nil {
		return nil, err
	}

	if join, err := innerJoin(lex, tableName); err != nil {
		return nil, err
	} else if join != nil {
		return join, nil
	}
	return tableName, nil
}

func tableName(lex lexer.Lexer) (*ast.Relation, error) {
	if !lex.Next() {
		return nil, fmt.Errorf("expected table, found nothing")
	}
	token := lex.Token()
	if token.Type != lexer.IdentifierType {
		return nil, fmt.Errorf("expected table name, found %q", token.Raw)
	}
	return &ast.Relation{Name: token.Raw}, nil
}

func innerJoin(lex lexer.Lexer, left *ast.Relation) (*ast.InnerJoin, error) {
	if ok, err := ifKeywords(lex, "INNER", "JOIN"); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	right, err := tableName(lex)
	if err != nil {
		return nil, err
	}

	if ok, err := ifKeywords(lex, "ON"); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("INNER JOIN requires ON")
	}

	on, err := fieldEqualsField(lex)
	if err != nil {
		return nil, err
	}

	return &ast.InnerJoin{
		Left:  left,
		Right: right,
		On:    on,
	}, nil
}

func field(lex lexer.Lexer) (*ast.Attribute, error) {
	if !lex.Next() {
		return nil, fmt.Errorf("expected field name, found nothing")
	}
	token := lex.Token()
	if token.Type != lexer.IdentifierType {
		return nil, fmt.Errorf("expected field name, found %q", token.Raw)
	}

	result := &ast.Attribute{Name: token.Raw}

	// Check if this is a qualified field name.
	if !lex.Next() {
		return result, nil
	}
	token = lex.Token()
	if token.Type != lexer.PeriodType {
		lex.UnreadToken()
		return result, nil
	}

	if !lex.Next() {
		return nil, fmt.Errorf("expected an attribute, found nothing")
	}
	token = lex.Token()
	if token.Type != lexer.IdentifierType {
		return nil, fmt.Errorf("expected field name, found %q", token.Raw)
	}

	result.Qualifier = result.Name
	result.Name = token.Raw
	return result, nil
}

func ifToken(lex lexer.Lexer, t lexer.Type) (bool, error) {
	if !lex.Next() {
		return false, nil
	}
	switch lex.Token().Type {
	case t:
		return true, nil
	case lexer.ErrorType:
		return false, fmt.Errorf("unable to get next token")
	default:
		lex.UnreadToken()
		return false, nil
	}
}

func fieldEqualsField(lex lexer.Lexer) (*ast.EqualColumnCondition, error) {
	left, err := field(lex)
	if err != nil {
		return nil, err
	}

	if ok, err := ifToken(lex, lexer.EqualType); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("expected = after %q", left.Name)
	}

	right, err := field(lex)
	if err != nil {
		return nil, err
	}

	return &ast.EqualColumnCondition{
		Left:  left,
		Right: right,
	}, nil
}

func ifKeywords(lex lexer.Lexer, keyword string, keywords ...string) (bool, error) {
	if !lex.Next() {
		return false, nil
	}
	token := lex.Token()
	if token.Type == lexer.ErrorType {
		return false, fmt.Errorf("could not tokenize input: %v", token.Raw)
	}
	if token.Type == lexer.EOFType {
		return false, nil
	}
	if token.Type != lexer.IdentifierType || strings.ToUpper(token.Raw) != keyword {
		lex.UnreadToken()
		return false, nil
	}

	for _, keyword := range keywords {
		if !lex.Next() {
			return false, fmt.Errorf("expected keyword %q, found nothing", keyword)
		}
		token := lex.Token()
		if token.Type == lexer.ErrorType {
			return false, fmt.Errorf("could not tokenize input: %v", token.Raw)
		}
		if token.Type == lexer.EOFType {
			return false, fmt.Errorf("expected keyword %q, found nothing", keyword)
		}
		if token.Type != lexer.IdentifierType || strings.ToUpper(token.Raw) != keyword {
			return false, fmt.Errorf("expected keyword %q, found nothing", keyword)
		}
	}
	return true, nil
}

func where(lex lexer.Lexer) (ast.Condition, error) {
	if ok, err := ifKeywords(lex, "WHERE"); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return condition(lex)
}

func condition(lex lexer.Lexer) (ast.Condition, error) {
	field, err := field(lex)
	if err != nil {
		return nil, err
	}
	result := ast.EqualCondition{LHS: field}
	//if !lex.Next() {
	//	return nil, fmt.Errorf("expected an attribute, found nothing")
	//}
	//token = lex.Token()
	//if token.Type != lexer.IdentifierType {
	//	return nil, fmt.Errorf("expected attribute name, found %q", token.Raw)
	//}
	//result.LHS = &ast.Attribute{Name: token.Raw}

	if ok, err := ifToken(lex, lexer.EqualType); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("expected =")
	}
	//if !lex.Next() {
	//	return nil, fmt.Errorf("expected =, found nothing")
	//}
	//token = lex.Token()
	//if token.Type != lexer.EqualType {
	//	return nil, fmt.Errorf("expected =, found %q", token.Raw)
	//}

	if !lex.Next() {
		return nil, fmt.Errorf("expected an attribute, found nothing")
	}
	token := lex.Token()
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
	if ok, err := ifKeywords(lex, "ORDER", "BY"); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	result := ast.OrderBy{}
	oc, err := orderByClause(lex)
	if err != nil {
		return nil, err
	}
	result.Criteria = append(result.Criteria, oc)

	for {
		if ok, err := ifToken(lex, lexer.CommaType); err != nil {
			return nil, err
		} else if !ok {
			break
		}

		oc, err := orderByClause(lex)
		if err != nil {
			return nil, err
		}
		result.Criteria = append(result.Criteria, oc)
	}

	return &result, nil
}

func orderByClause(lex lexer.Lexer) (*ast.OrderCriteria, error) {
	field, err := field(lex)
	if err != nil {
		return nil, err
	}

	oc := &ast.OrderCriteria{
		Attribute: field,
		SortOrder: ast.Asc,
	}
	if lex.Next() {
		token := lex.Token()
		if token.Type == lexer.IdentifierType && strings.ToUpper(token.Raw) == "ASC" {
			oc.SortOrder = ast.Asc
		} else if token.Type == lexer.IdentifierType && strings.ToUpper(token.Raw) == "DESC" {
			oc.SortOrder = ast.Desc
		} else {
			lex.UnreadToken()
		}
	}
	return oc, nil
}
