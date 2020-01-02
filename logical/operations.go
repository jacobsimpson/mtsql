package logical

import (
	"fmt"

	md "github.com/jacobsimpson/mtsql/metadata"
)

type Operation interface {
	Children() []Operation
	Clone(...Operation) Operation
	Provides() []*md.Column
	Requires() []*md.Column
	String() string
}

type Difference struct {
	LHS Operation
	RHS Operation
}

type Intersection struct {
	LHS Operation
	RHS Operation
}

type Product struct {
	LHS Operation
	RHS Operation
}

type Union struct {
	LHS Operation
	RHS Operation
}

type Selection struct {
	Child    Operation
	requires []*md.Column
}

func NewSelection(child Operation, requires []*md.Column) *Selection {
	return &Selection{
		Child:    child,
		requires: requires,
	}
}

type Projection struct {
	Child   Operation
	columns []*md.Column
}

func NewProjection(child Operation, columns []*md.Column) *Projection {
	return &Projection{
		Child:   child,
		columns: columns,
	}
}

type Distinct struct {
	Child Operation
}

type Sort struct {
	Child Operation
}

type Source struct {
	Name     string
	provides []*md.Column
}

func (o *Union) Children() []Operation {
	return []Operation{o.LHS, o.RHS}
}

func (o *Union) Clone(children ...Operation) Operation {
	if len(children) != 2 {
		panic("wrong number of children")
	}
	return &Union{
		LHS: children[0],
		RHS: children[1],
	}
}

func (o *Union) String() string {
	return fmt.Sprintf("Union{LHS: %s, RHS: %s}", o.LHS, o.RHS)
}

func (o *Union) Provides() []*md.Column { return o.LHS.Provides() }
func (o *Union) Requires() []*md.Column { return []*md.Column{} }

func (o *Intersection) Children() []Operation {
	return []Operation{o.LHS, o.RHS}
}

func (o *Intersection) Clone(children ...Operation) Operation {
	if len(children) != 2 {
		panic("wrong number of children")
	}
	return &Intersection{
		LHS: children[0],
		RHS: children[1],
	}
}

func (o *Intersection) String() string {
	return fmt.Sprintf("Intersection{LHS: %s, RHS: %s}", o.LHS, o.RHS)
}

func (o *Intersection) Provides() []*md.Column { return o.LHS.Provides() }
func (o *Intersection) Requires() []*md.Column { return []*md.Column{} }

func (o *Difference) Children() []Operation {
	return []Operation{o.LHS, o.RHS}
}

func (o *Difference) Clone(children ...Operation) Operation {
	if len(children) != 2 {
		panic("wrong number of children")
	}
	return &Difference{
		LHS: children[0],
		RHS: children[1],
	}
}

func (o *Difference) String() string {
	return fmt.Sprintf("Difference{LHS: %s, RHS: %s}", o.LHS, o.RHS)
}

func (o *Difference) Provides() []*md.Column { return o.LHS.Provides() }
func (o *Difference) Requires() []*md.Column { return []*md.Column{} }

func (o *Selection) Children() []Operation {
	return []Operation{o.Child}
}

func (o *Selection) Clone(children ...Operation) Operation {
	if len(children) != 1 {
		panic("wrong number of children")
	}
	return &Selection{
		Child:    children[0],
		requires: o.requires,
	}
}

func (o *Selection) String() string {
	return fmt.Sprintf("Selection{Child: %s}", o.Child)
}

func (o *Selection) Provides() []*md.Column { return o.Child.Provides() }
func (o *Selection) Requires() []*md.Column { return o.requires }

func (o *Projection) Children() []Operation {
	return []Operation{o.Child}
}

func (o *Projection) Clone(children ...Operation) Operation {
	if len(children) != 1 {
		panic("wrong number of children")
	}
	return &Projection{
		Child:   children[0],
		columns: o.columns,
	}
}

func (o *Projection) String() string {
	return fmt.Sprintf("Projection{Child: %s}", o.Child)
}

func (o *Projection) Provides() []*md.Column { return o.columns }
func (o *Projection) Requires() []*md.Column { return o.columns }

func (o *Product) Children() []Operation {
	return []Operation{o.LHS, o.RHS}
}

func (o *Product) Clone(children ...Operation) Operation {
	if len(children) != 2 {
		panic("wrong number of children")
	}
	return &Product{
		LHS: children[0],
		RHS: children[1],
	}
}

func (o *Product) String() string {
	return fmt.Sprintf("Product{LHS: %s, RHS: %s}", o.LHS, o.RHS)
}

func (o *Product) Provides() []*md.Column { return append(o.LHS.Provides(), o.RHS.Provides()...) }
func (o *Product) Requires() []*md.Column { return []*md.Column{} }

func (o *Distinct) Children() []Operation {
	return []Operation{o.Child}
}

func (o *Distinct) Clone(children ...Operation) Operation {
	if len(children) != 1 {
		panic("wrong number of children")
	}
	return &Distinct{Child: children[0]}
}

func (o *Distinct) String() string {
	return fmt.Sprintf("Distinct{Child: %s}", o.Child)
}

func (o *Distinct) Provides() []*md.Column { return o.Child.Provides() }
func (o *Distinct) Requires() []*md.Column { return []*md.Column{} }

func (o *Sort) Children() []Operation {
	return []Operation{o.Child}
}

func (o *Sort) Clone(children ...Operation) Operation {
	if len(children) != 1 {
		panic("wrong number of children")
	}
	return &Sort{Child: children[0]}
}

func (o *Sort) String() string {
	return fmt.Sprintf("Sort{Child: %s}", o.Child)
}

func (o *Sort) Provides() []*md.Column { return o.Child.Provides() }
func (o *Sort) Requires() []*md.Column { return []*md.Column{} }

func NewSource(name string, provides []*md.Column) *Source {
	return &Source{
		Name:     name,
		provides: provides,
	}
}

func (o *Source) Children() []Operation {
	return []Operation{}
}

func (o *Source) Clone(children ...Operation) Operation {
	if len(children) != 0 {
		panic("wrong number of children")
	}
	return &Source{
		provides: o.provides,
	}
}

func (o *Source) String() string {
	return fmt.Sprintf("Source{Name: %q, provides: %s}", o.Name, o.provides)
}

func (o *Source) Provides() []*md.Column { return o.provides }
func (o *Source) Requires() []*md.Column { return []*md.Column{} }
