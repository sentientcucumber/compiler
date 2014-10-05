// Author: Michael Hunsinger
// Date:   Oct 4 2014
// File:   types.go

package compiler

type MarkedVocabulary struct {
	vocabulary    map[string]bool
}

type Grammar struct {
	terminals     map[string]bool
	nonterminals  map[string]bool
	productions   map[string]bool
	rhs           map[string]bool
	lhs           map[string]bool
}

type TermSet struct {
	symbols       []string
}
