// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

/*
Chris Liu changes:
* Restructure the code into runnable go unit test cases
*/

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

	exp = "{}"
	if s := y.String(); s != exp {
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

func TestLen(t *testing.T) {
	var x IntSet
	var exp int
	exp = 0
	if i := x.Len(); i != exp {
		t.Errorf(`got %s, expect %s`, i, exp)
	}

	x.Add(0)
	exp++
	if i := x.Len(); i != exp {
		t.Errorf(`got %s, expect %s`, i, exp)
	}

	x.Add(9)
	exp++
	if i := x.Len(); i != exp {
		t.Errorf(`got %s, expect %s`, i, exp)
	}
}

func TestRemove(t *testing.T) {
	var x IntSet
	var exp string
	exp = "{}"
	if s := x.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
	x.Add(0)
	exp = "{0}"
	if s := x.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
	x.Remove(0)
	exp = "{}"
	if s := x.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
	x.Add(32)
	x.Add(64)
	exp = "{32 64}"
	if s := x.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
	x.Remove(0)
	exp = "{32 64}"
	if s := x.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
	x.Remove(64)
	exp = "{32}"
	if s := x.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
}

func TestClear(t *testing.T) {
	var x IntSet
	var exp int
	exp = 0
	if i := x.Len(); i != exp {
		t.Errorf(`got %s, expect %s`, i, exp)
	}
	x.Add(32)
	x.Add(320)
	exp = 2
	if i := x.Len(); i != exp {
		t.Errorf(`got %s, expect %s`, i, exp)
	}
	x.Clear()
	exp = 0
	if i := x.Len(); i != exp {
		t.Errorf(`got %s, expect %s`, i, exp)
	}
}

func TestCopy(t *testing.T) {
	var x IntSet
	var exp, expY string
	exp = "{}"
	if s := x.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
	x.Add(0)
	exp = "{0}"
	if s := x.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
	y := x.Copy()
	expY = "{0}"
	x.Remove(0)
	exp = "{}"
	if s := x.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
	if s := y.String(); s != expY {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
	y.Add(320)
	expY = "{0 320}"
	if s := y.String(); s != expY {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
	x.Add(1)
	exp = "{1}"
	if s := x.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
}

func TestAddAll(t *testing.T) {
	var x IntSet
	var exp string
	x.AddAll(0, 1, 2, 3)
	exp = "{0 1 2 3}"
	if s := x.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
}

func TestUnionWith(t *testing.T) {
	var x, y IntSet
	var exp string
	x.AddAll(0, 1, 2, 3)
	y.AddAll(2, 3, 4, 320)
	x.UnionWith(&y)
	exp = "{0 1 2 3 4 320}"
	if s := x.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
}

func TestIntersectWith(t *testing.T) {
	var x, y IntSet
	var exp string
	x.AddAll(0, 1, 2, 3, 144)
	y.AddAll(2, 3, 4, 144, 320)
	x.IntersectWith(&y)
	exp = "{2 3 144}"
	if s := x.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
}

func TestDifferenceWith(t *testing.T) {
	var x, y IntSet
	var exp string
	x.AddAll(0, 1, 2, 3)
	y.AddAll(2, 3, 4, 320)
	z := x.Copy()
	x.DifferenceWith(&y)
	exp = "{0 1}"
	if s := x.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
	y.DifferenceWith(z)
	exp = "{4 320}"
	if s := y.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
}

func TestSymmetricDifference(t *testing.T) {
	var x, y IntSet
	var exp string
	x.AddAll(0, 1, 2, 3)
	y.AddAll(2, 3, 4, 320)
	x.SymmetricDifference(&y)
	exp = "{0 1 4 320}"
	if s := x.String(); s != exp {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
}

func TestElems(t *testing.T) {
	var x IntSet
	var exp []int
	x.AddAll(0, 1, 2, 3)
	exp = []int{0, 1, 2, 3}
	if s := x.Elems(); !compareSlice(s, exp) {
		t.Errorf(`got %s, expect %s`, s, exp)
	}
	x.AddAll(0, 1, 2, 3, 320)
	exp = []int{0, 1, 2, 3, 320}
	if s := x.Elems(); !compareSlice(s, exp) {
		t.Errorf(`got %s, expect %s`, s, exp)
	}

}

func TestProperContains(t *testing.T) {
	var x, y IntSet
	x.AddAll(0, 1, 2, 3)
	if !x.ProperContains(&y) {
		t.Errorf(`%s should proper contains %s`, &x, &y)
	}
	y.AddAll(0, 1, 2)
	if !x.ProperContains(&y) {
		t.Errorf(`%s should proper contains %s`, &x, &y)
	}
	y.Add(3)
	if x.ProperContains(&y) {
		t.Errorf(`%s should not proper contains %s`, &x, &y)
	}
	y.Add(4)
	if !y.ProperContains(&x) {
		t.Errorf(`%s should proper contains %s`, &y, &x)
	}
	x.Add(4)
	if y.ProperContains(&x) {
		t.Errorf(`%s should not proper contains %s`, &y, &x)
	}
	x.Clear()
	if !y.ProperContains(&x) {
		t.Errorf(`%s should proper contains %s`, &y, &x)
	}
	y.Clear()
	if y.ProperContains(&x) {
		t.Errorf(`%s should not proper contains %s`, &y, &x)
	}
	y.Add(100)
	if !y.ProperContains(&x) {
		t.Errorf(`%s proper contains %s`, &y, &x)
	}
	x.AddAll(1, 100)
	if !x.ProperContains(&y) {
		t.Errorf(`%s proper contains %s`, &x, &y)
	}
}

func compareSlice(s1, s2 []int) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}
	return true
}

func BenchmarkAddThenRemove(b *testing.B) {
	var x IntSet
	for i := 0; i < b.N; i++ {
		x.Add(i)
		x.Remove(i)
	}
}
