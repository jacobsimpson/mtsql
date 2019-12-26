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
	RHS *Attribute
}
type LikeCondition struct {
	LHS *Attribute
	RHS string
}
