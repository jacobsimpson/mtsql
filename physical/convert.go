package physical

import (
	"fmt"

	"github.com/jacobsimpson/mtsql/logical"
	md "github.com/jacobsimpson/mtsql/metadata"
)

func Convert(o logical.Operation, tables map[string]*md.Relation) (RowReader, error) {
	if o == nil {
		return nil, fmt.Errorf("unable to covert nil value")
	}

	if _, ok := o.(*logical.Difference); ok {
	}

	if _, ok := o.(*logical.Intersection); ok {
	}

	if _, ok := o.(*logical.Product); ok {
	}

	if _, ok := o.(*logical.Union); ok {
	}

	if s, ok := o.(*logical.Selection); ok {
		rr, err := Convert(s.Child, tables)
		if err != nil {
			return nil, err
		}
		return NewColumnFilter(rr, nil, nil)
	}

	if p, ok := o.(*logical.Projection); ok {
		rr, err := Convert(p.Child, tables)
		if err != nil {
			return nil, err
		}
		return NewProjection(rr, p.Provides())
	}

	if _, ok := o.(*logical.Distinct); ok {
	}

	if s, ok := o.(*logical.Sort); ok {
		rr, err := Convert(s.Child, tables)
		if err != nil {
			return nil, err
		}
		return NewSortScan(rr, []SortScanCriteria{})
	}

	if s, ok := o.(*logical.Source); ok {
		relation := tables[s.Name]
		return NewTableScan(relation.Name, relation.Source)
	}

	return nil, nil
}
