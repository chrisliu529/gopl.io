// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 165.

// Package intset provides a set of integers based on a bit vector.

/*
Chris Liu changes:
* Add functions:
    Len(), Remove(), Clear(), Copy(),
    AddAll(),
    IntersectWith(), DifferenceWith(), SymmetricDifference(),
    Elems(),
    Contains()
*/

package intset

import (
	"bytes"
	"fmt"
)

//!+intset

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint
}

const (
	wordBits = 32<<(^uint(0)>>63)
)

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/wordBits, uint(x%wordBits)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/wordBits, uint(x%wordBits)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

func (s *IntSet) AddAll(vals ...int) {
	for _, val := range vals {
		s.Add(val)
	}
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// IntersectWith sets s to the intersection of s and t.
func (s *IntSet) IntersectWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &= tword
		} else {
			return
		}
	}
}

// DifferenceWith sets s to the difference of s and t.
func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &^= tword
		} else {
			return
		}
	}
}

// SymmetricDifference sets s to the symmetric difference of s and t.
func (s *IntSet) SymmetricDifference(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] ^= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

//!-intset

//!+string

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < wordBits; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", wordBits*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

//!-string

func (s *IntSet) Len() int { //return the number of elements
	n := 0
	for _, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < wordBits; j++ {
			if word&(1<<uint(j)) != 0 {
				n++
			}
		}
	}
	return n
}

func (s *IntSet) Remove(x int) { // return x from the set
	word, bit := x/wordBits, uint(x%wordBits)
	if word >= len(s.words) {
		return
	}
	s.words[word] &^= 1 << bit
}

func (s *IntSet) Clear() { // remove all elements from the set
	s.words = nil
}

func (s *IntSet) Copy() *IntSet { // return a copy of the set
	words2 := make([]uint, len(s.words))
	copy(words2, s.words)
	return &IntSet{words: words2}
}

func (s *IntSet) Elems() []int {
	elems := []int{}
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < wordBits; j++ {
			if word&(1<<uint(j)) != 0 {
				elems = append(elems, wordBits*i+j)
			}
		}
	}
	return elems
}

func (s *IntSet) ProperContains(t *IntSet) bool {
	if s.Len() <= t.Len() {
		return false
	}
	for _, e := range t.Elems() {
		if !s.Has(e) {
			return false
		}
	}
	return true
}
