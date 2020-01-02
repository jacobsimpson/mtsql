package algebra

import (
	"fmt"

	md "github.com/jacobsimpson/mtsql/metadata"
)

func Optimize(o Operation) Operation {
	return PushDownSelection(o)
}

func PushDownSelection(o Operation) Operation {
	if s, ok := o.(*Selection); ok {
		return helpPushDownSelection(s.Child, s)
	}
	result := []Operation{}
	for _, c := range o.Children() {
		result = append(result, PushDownSelection(c))
	}
	return o.Clone(result...)
}

func helpPushDownSelection(o Operation, s *Selection) Operation {
	result := []Operation{}
	didPushDown := false
	for _, c := range o.Children() {
		if canPushDownSelection(o, s) {
			result = append(result, helpPushDownSelection(c, s))
			didPushDown = true
		} else {
			result = append(result, c)
		}
	}
	if didPushDown {
		return o.Clone(result...)
	}
	return s.Clone(o)
}

func canPushDownSelection(o Operation, s *Selection) bool {
	fmt.Println("1")
	if !containsAll(o.Provides(), s.Requires()) {
		fmt.Println("2")
		return false
	}
	fmt.Println("3")
	if _, ok := o.(*Union); ok {
		fmt.Println("4")
		return ok
	}
	fmt.Println("5")
	if _, ok := o.(*Product); ok {
		fmt.Println("6")
		return ok
	}
	fmt.Println("7")
	if _, ok := o.(*Projection); ok {
		fmt.Println("8")
		return ok
	}
	fmt.Println("9")
	return false
}

func containsAll(provides []*md.Column, requires []*md.Column) bool {
	for _, r := range requires {
		found := false
		for _, p := range provides {
			if r.QualifiedName() == p.QualifiedName() {
				found = true
				break
			}
		}
		if found == false {
			return false
		}
	}
	return true
}
