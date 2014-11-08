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
	EOP           SymbolCategory = "EOP"
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

// Operator type
type Operator string

// ExprKind type
type ExprKind string

// OpRec type
type OpRec struct {
	Op          int
}

// ExprRec type
type ExprRec struct {
	Kind        ExprKind
	Name        string
	Val         int
}

// Enumerations for kinds
const (
	IdExpr      ExprKind = "IdExpr"
	LiteralExpr ExprKind = "LiteralExpr"
	TempExpr    ExprKind = "TempExpr"
)

type SemanticRecord struct {
	exprRec     ExprRec
	opRec       OpRec
}
