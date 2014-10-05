// Author: Michael Hunsinger
// Date:   Oct 4 2014
// File:   generator.go
// Implementation of a predict generator for LL(1) grammars

package compiler

import (
	"strings"
)

// Mark which parts of a vocabulary (terminals and nonterminals) from a grammar
// can produce lambda. If reading an LL(1) grammar, the grammar should be
// formatted that the LHS produces nothing instead of nil or a lambda unicode
// (e.g. "<lhs> -> ")
func MarkLambda (g Grammar) MarkedVocabulary {

	derivesLambda := pullVocabulary(g)
	changes := true
	
	for k, _ := range derivesLambda.vocabulary {
		derivesLambda.vocabulary[k] = false
	}

	for changes {
		changes = false

		for p := range g.productions {
			rhsDerivesLambda := true
			rhs := stripRhs(p)
			
			for _, s := range stripSymbols(rhs) {
				rhsDerivesLambda = rhsDerivesLambda && derivesLambda.vocabulary[s];
			}

			lhs := stripLhs(p)
			if rhsDerivesLambda && !derivesLambda.vocabulary[lhs] {
				changes = true
				derivesLambda.vocabulary[lhs] = true
			}
		}
	}

	return derivesLambda
}

// Determines the first terminal or lambda for a given set of symbols,
// terminals and nonterminals
func computeFirst (s string) (result TermSet) {
	strs := strings.Fields(s)

	if k := len(strs); k == 0 {
		result = TermSet {}
	} else {
		result.symbols = append(result.symbols, firstSet(strs[0]))
		i := 0

		for i < k && firstSet(strs[i]) == "" {
			i++
			result.symbols = append(result.symbols, firstSet(strs[i]))
		}

		if i == k && firstSet(strs[k]) == "" {
			result.symbols = append(result.symbols, "lambda")
		}
	}

	return
}

// Use in conjunction with ComputeFirst to fill the FirstSet
func fillFirstSet() {
	for A := range nonterminals {
		if derivesLambda(A) {
			firstSet(A) = ""
		} else {
			firstSet(A) = "" 	// not sure what this should be 
		}
	}

	for a := range terminals {
		firstSet(a) := a

		for A := range nonterminals {
			if first(rhs) == a { // first is the first symbol of a production
				firstSet(A) := firstSet(A) + a;
			}
		}
	}

	for p := range productions {
		firstSet(stripLhs(p)) := firstSet(stripLhs) + computeFirst(stripRhs(p))
	}							// exit when no changes
}

// Returns the first set of a given string???
func firstSet (s string) string {
	if len(s) == 0 {
		return ""
	} else {
		return s
	}
}

// Pull the vocabulary from a grammar
func pullVocabulary (g Grammar) MarkedVocabulary {
	v := make(map[string]bool, 0)

	for k, _ := range g.nonterminals {
		v[k] = g.nonterminals[k]
	}

	for k, _ := range g.terminals {
		v[k] = g.terminals[k]
	}

	return MarkedVocabulary { v }
}
