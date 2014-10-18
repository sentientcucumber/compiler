// Author: Michael Hunsinger
// Date:   Oct 10 2014
// File:   types.go
// Various definitions for types used in the compiler

package compiler

import (
	"sort"
	"fmt"
)

// Definition of a Symbol
type Symbol struct {
	name          string
	category      SymbolCategory
}

// Enumerations used for the SymbolCategory
type SymbolCategory string
const (
	TERMINAL      SymbolCategory = "TERMINAL"
	NONTERMINAL   SymbolCategory = "NONTERMINAL"
	LAMBDA        SymbolCategory = "LAMBDA"
)

// Definition of a MarkedVocabulary
type MarkedVocabulary map[string]bool

// Definition of a TermSet
type TermSet []Symbol

// Definition of a TermSet
type Set map[string][]Symbol

// Determines if a set contains a Symbol, takes a key k and v Symbol to find
// If found, it will return true and the index it was found, otherwise it will
// return false and return -1
func (s Set) contains(k string, v Symbol) (bool, int) {
	for i, e := range s[k] {
		if e == v {
			return true, i
		}
	}

	return false, -1
}

// Checks to see if a set contains lambda, if true, return the index it was 
// found at, otherwise return false and -1
func (s Set) containsLambda(k string) (bool, int) {
	l := Symbol { name: "λ" }

	for i, e := range s[k] {
		if e.name == l.name {
			return true, i
		}
	}

	return false, -1
}

// Checks to see if a set contains a value, removes it and returns a copy of
// the Set if found, otherwise it returns the same Set
func (s Set) remove(k string, v Symbol) Set {
	n := s

	if b, i := s.contains(k, v); b {
		copy(n[k][i:], n[k][i+1:])
		n[k] = n[k][:len(n[k]) - 1]

		return n
	} else {
		return n
	}
}

// Checks to see if a set contains lambda, removes it and returns a copy of
// the Set if found, otherwise it returns the same Set
func (s Set) removeLambda(k string) Set {
	n := s
	b, i := s.containsLambda(k)
	
	for b  {
		copy(n[k][i:], n[k][i+1:])
		n[k] = n[k][:len(n[k]) - 1]

		b, i = s.containsLambda(k)
	}
	
	return n
}

// Adds a unique value to the Set's array of Symbols.
func (s Set) add(k string, a...Symbol) {

	for _, v := range a {
		if b, _ := s.contains(k, v); !b {
			s[k] = append(s[k], v)
		}
	}
}

// Prints a set in the prettiest, most consistent way possible
func (s Set) print() {
	var (
		keys []string
	)

	for k := range s {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%s -> ", k)
		for _, v := range s[k] {
			fmt.Printf("%v ", v.name)
		}
		fmt.Printf("\n")
	}
}

// Definition of a Grammar
type Grammar struct {
	terminals     map[string]bool
	nonterminals  map[string]bool
	productions   map[string]bool
	rhs           map[string]bool
	lhs           map[string]bool
}

