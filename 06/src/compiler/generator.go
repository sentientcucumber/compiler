// Author: Michael Hunsinger
// Date:   Oct 4 2014
// File:   generator.go
// Implementation of a predict generator for LL(1) grammars

package compiler

import (
	"fmt"
	"strings"
)

var (
	g = Grammar {
		terminals:     terminals,
		nonterminals:  nonterminals,
		productions:   productions,
		rhs:           rhs,
		lhs:           lhs,
	}

	FirstSet = make(map[string][]string, 0)
	derivesLambda = pullVocabulary(g)
)


// Mark which parts of a vocabulary (terminals and nonterminals) from a grammar
// can produce lambda. If reading an LL(1) grammar, the grammar should be
// formatted that the LHS produces nothing instead of nil or a lambda unicode
// (e.g. "<lhs> -> ")
func MarkLambda (g Grammar) MarkedVocabulary {
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
	fmt.Printf("compute first s: %v\n", s)
	strs := strings.Fields(s)

	if k := len(strs); k == 0 {
		result.symbols = append(result.symbols, "")
	} else {
		fmt.Printf("compute first FirstSet[strs[0]]: '%v'\n", FirstSet[strs[0]])
		t := remove(FirstSet[strs[0]], "") // Remove lambda from FirstSet

		result.symbols = t
		i := 0
		
		for b, _ := contains(FirstSet[strs[i]], ""); i < k && b; {
			i++
			t = remove(FirstSet[strs[i]], "")

			fmt.Printf("result.symbols before %v\n", result.symbols)
			result.symbols = append(result.symbols, t...)
			fmt.Printf("result.symbols after %v\n", result.symbols)
		}

		if b, _ := contains(FirstSet[strs[k - 1]], ""); i == k - 1 && b {
			result.symbols = append(result.symbols, "")
		}
	}
	
	return
}


// Use in conjunction with ComputeFirst to fill the FirstSet
func FillFirstSet() {
	derivesLambda := MarkedVocabulary { nonterminals }

	for A := range g.nonterminals {
		if derivesLambda.vocabulary[A] {
			FirstSet[A] = append(FirstSet[A], "")
		} else {
			// add an empty set, so add nothing?
		}
	}

	for a := range g.terminals {
		FirstSet[a] = append(FirstSet[a], a)

		for A := range g.nonterminals {
			for p := range g.productions {
				rhs := stripRhs(p)

				if firstTerm(rhs) == a { // first is the first symbol of a production
					FirstSet[A] = append(FirstSet[A], a);
				}
			}
		}
	}

	for p := range g.productions {
		lhs := stripLhs(p)
		rhs := stripRhs(p)
		first := computeFirst(rhs).symbols

		FirstSet[lhs] = append(FirstSet[lhs], first...)
	}
}

func contains(a []string, v string) (found bool, ind int) {
	found = false

	for i, e := range a {
		if e == v {
			found = true
			ind = i
			break
		}
	}

	return
}

func remove(a []string, s string) []string {

	if b, i := contains(a, s); b {
		copy(a[i:], a[i+1:])
		a = a[:len(a) - 1]
	}

	return a
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
