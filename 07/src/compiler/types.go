// Author: Michael Hunsinger
// Date:   Oct 10 2014
// File:   types.go
// Various definitions for types used in the compiler

package compiler

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

// Definition of a Grammar
type Grammar struct {
	terminals     map[string]bool
	nonterminals  map[string]bool
	productions   map[string]bool
	staticProd    map[int]string
	rhs           map[string]bool
	lhs           map[string]bool
}
