package ast

type Query interface{}

type SFW struct {
	SelList   *SelList
	From      From
	Condition Condition
}

type SelList struct {
	All        bool
	Attributes []*Attribute
}

type Attribute struct {
	Name string
}

type From interface{}
type Relation struct {
	Name string
}

type Condition interface{}
type AndCondition struct {
	LHS Condition
	RHS Condition
}
type InCondition struct{}
type EqualCondition struct {
	LHS *Attribute
	RHS *Constant
}
type LikeCondition struct {
	LHS *Attribute
	RHS string
}
type Type string

const (
	StringType  Type = "string"
	IntegerType Type = "integer"
)

type Constant struct {
	Type  Type
	Value interface{}
	Raw   string
}
