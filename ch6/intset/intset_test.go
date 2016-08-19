// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package intset

import (
	"testing"
)

func TestString(t *testing.T) {
	var x, y IntSet
	var exp string
	x.Add(1)
	x.Add(144)
	x.Add(9)
	exp = "{1 9 144}"
	if s := x.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}

	y.Add(9)
	y.Add(42)
	exp = "{9 42}"
	if s := y.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}

	x.UnionWith(&y)
	exp = "{1 9 42 144}"
	if s := x.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
}

func TestHas(t *testing.T) {
	var x IntSet
	x.Add(9)
	if !x.Has(9) {
		t.Errorf("should has 9")
	}
	if x.Has(123) {
		t.Errorf("should not has 123")
	}
}
