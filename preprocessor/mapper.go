package preprocessor

import (
	"fmt"

	"github.com/jacobsimpson/mtsql/ast"
	md "github.com/jacobsimpson/mtsql/metadata"
)

type mapper struct {
	names     map[string][]*md.Column
	qualified map[string][]*md.Column
	aliases   map[string][]*md.Column
}

func newMapper(columns []*md.Column) *mapper {
	result := &mapper{
		names:     map[string][]*md.Column{},
		qualified: map[string][]*md.Column{},
		aliases:   map[string][]*md.Column{},
	}
	for _, c := range columns {
		result.names[c.Name] = append(result.names[c.Name], c)
		if c.Qualifier != "" {
			result.qualified[c.QualifiedName()] = append(result.qualified[c.QualifiedName()], c)
		}
		if c.Alias != "" {
			result.aliases[c.Alias] = append(result.aliases[c.Alias], c)
		}
	}
	return result
}

func (m *mapper) findMatches(a *ast.Attribute) ([]*md.Column, error) {
	if a.Alias != "" {
		r := m.aliases[a.Alias]
		if r == nil {
			return nil, fmt.Errorf("no matching alias for %q", a.Alias)
		}
		if len(r) > 1 {
			return nil, fmt.Errorf("too many matching aliases %q", a.Alias)
		}
		return r, nil
	}
	if a.Qualifier != "" {
		qualifiedName := fmt.Sprintf("%s.%s", a.Qualifier, a.Name)
		r := m.qualified[qualifiedName]
		if r == nil {
			return nil, fmt.Errorf("no matching qualified name %q", qualifiedName)
		}
		if len(r) > 1 {
			return nil, fmt.Errorf("too many matching qualified names %q", qualifiedName)
		}
		return r, nil
	}
	if a.Name == "*" {
		r := []*md.Column{}
		for _, n := range m.names {
			r = append(r, n...)
		}
		return r, nil
	}
	r := m.names[a.Name]
	if r == nil {
		return nil, fmt.Errorf("no matching name %q", a.Name)
	}
	if len(r) > 1 {
		return nil, fmt.Errorf("too many matching names %q", a.Name)
	}
	return r, nil
}
