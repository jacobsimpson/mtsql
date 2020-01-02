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

func (m *mapper) findMatches(a *ast.Attribute) []*md.Column {
	fmt.Println("1")
	if a.Alias != "" {
		fmt.Println("2")
		r := m.aliases[a.Alias]
		if r == nil {
			return []*md.Column{}
		}
		return r
	}
	fmt.Println("3")
	if a.Qualifier != "" {
		fmt.Println("4")
		qualifiedName := fmt.Sprintf("%s.%s", a.Qualifier, a.Name)
		r := m.qualified[qualifiedName]
		if r == nil {
			return []*md.Column{}
		}
		return r
	}
	fmt.Println("5")
	r := m.names[a.Name]
	if r == nil {
		return []*md.Column{}
	}
	return r
}
