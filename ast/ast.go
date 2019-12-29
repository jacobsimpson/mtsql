package ast

type Query interface{}

type Profile struct {
	SFW *SFW
}

type SFW struct {
	SelList *SelList
	From    From
	Where   Condition
	OrderBy *OrderBy
}

type SelList struct {
	Attributes []*Attribute
}

type Attribute struct {
	Qualifier string
	Name      string
	Alias     string
}

type From interface {
	Tables() []*Relation
}
type Relation struct {
	Name string
}

func (r *Relation) Tables() []*Relation { return []*Relation{r} }

type InnerJoin struct {
	Left  *Relation
	Right *Relation
	On    *EqualColumnCondition
}

func (r *InnerJoin) Tables() []*Relation { return []*Relation{r.Left, r.Right} }

type Condition interface{}
type AndCondition struct {
	LHS Condition
	RHS Condition
}
type InCondition struct{}
type EqualColumnCondition struct {
	Left  *Attribute
	Right *Attribute
}
type EqualCondition struct {
	LHS *Attribute
	RHS *Constant
}
type LikeCondition struct {
	LHS *Attribute
	RHS string
}

type SortOrder string

const (
	Asc  SortOrder = "ASC"
	Desc SortOrder = "DESC"
)

type OrderCriteria struct {
	Attribute *Attribute
	SortOrder SortOrder
}

type OrderBy struct {
	Criteria []*OrderCriteria
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
