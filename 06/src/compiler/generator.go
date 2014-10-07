// Author: Michael Hunsinger
// Date:   Oct 4 2014
// File:   generator.go
// Implementation of a predict generator for LL(1) grammars

package compiler

import (
	"fmt"
	"strings"
	"sort"
)

var (
	g = Grammar {
		terminals:     terminals,
		nonterminals:  nonterminals,
		productions:   productions,
		rhs:           rhs,
		lhs:           lhs,
	}

	FirstSet        = make(map[string][]string, 0)
	FollowSet       = make(map[string][]string, 0)
	derivesLambda   MarkedVocabulary
)

// Generates a predict set
func Predict() {
	MarkLambda(g)
	FillFirstSet()
	FillFollowSet()

	for p := range g.productions {

		PredictSet := make([]string, 0)

		rhs := stripRhs(p)
		lhs := stripLhs(p)

		if rhs != "" {		
			PredictSet = append(PredictSet, FirstSet[rhs]...)
			fmt.Printf("First ( %s )", rhs)

			if b, _ := contains(FirstSet[rhs], ""); b {
				t := remove(FollowSet[lhs], "")
				PredictSet = append(PredictSet, t...)
				fmt.Printf("∪ Follow ( %s ) - λ", )
			}
			
			fmt.Printf(" = %s\n", PredictSet)
		}
	}
}

// Mark which parts of a vocabulary (terminals and nonterminals) from a grammar
// can produce lambda. If reading an LL(1) grammar, the grammar should be
// formatted that the LHS produces nothing instead of nil or a lambda unicode
// (e.g. "<lhs> -> ")
func MarkLambda (g Grammar) MarkedVocabulary {
	changes := true
	derivesLambda = pullVocabulary(g)
	
	for k, _ := range derivesLambda {
		derivesLambda[k] = false
	}

	for changes {
		changes = false

		for p := range g.productions {
			rhsDerivesLambda := true
			rhs := stripRhs(p)
			
			for _, s := range stripSymbols(rhs) {
				rhsDerivesLambda = rhsDerivesLambda && derivesLambda[s];
			}

			lhs := stripLhs(p)
			if rhsDerivesLambda && !derivesLambda[lhs] {
				changes = true
				derivesLambda[lhs] = true
			}
		}
	}

	return derivesLambda
}

// Determines the first terminal or lambda for a given set of symbols,
// terminals and nonterminals
func ComputeFirst (s string) (result TermSet) {
	strs := strings.Fields(s)

	if k := len(strs); k == 0 {
		result = append(result, "")
	} else {
		t := remove(FirstSet[strs[0]], "") // Remove lambda from FirstSet

		result = t
		i := 0
		
		for b, _ := contains(FirstSet[strs[i]], ""); i < k && b; {
			i++
			t = remove(FirstSet[strs[i]], "")

			result = append(result, t...)
		}

		if b, _ := contains(FirstSet[strs[k - 1]], ""); i == k - 1 && b {
			result = append(result, "")
		}
	}
	
	return
}


// Fill the FirstSet
func FillFirstSet() {
	MarkLambda(g)

	for A := range g.nonterminals {
		if derivesLambda[A] {
			FirstSet[A] = []string { "" }
		} else {
			FirstSet[A] = make([]string, 0)
		}
	}

	for a := range g.terminals {
		FirstSet[a] = []string { a }

		for A := range g.nonterminals {
			for p := range g.productions {
				rhs := stripRhs(p)
				lhs := stripLhs(p)

				// Extra bit to make sure this is a "set"
				if _, s := firstTerm(rhs); s == a && lhs == A {
					if b, _ := contains(FirstSet[A], a); !b {
						FirstSet[A] = append(FirstSet[A], a);
					}
				}
			}
		}
	}

	// TODO this is poor programming... 
	for i := 0; i < 2; i++ {
		for p := range g.productions {
			lhs := stripLhs(p)
			rhs := stripRhs(p)
			first := ComputeFirst(rhs)

			// Extra bit to make sure this is a "set"
			for i, _ := range first {
				if b, _ := contains(FirstSet[lhs], first[i]); !b {
					FirstSet[lhs] = append(FirstSet[lhs], first[i])
				}
			}
		}
	}
}


// Fill the FollowSet
func FillFollowSet() {
	MarkLambda(g)
	FillFirstSet()

	for A := range g.nonterminals {
		FollowSet[A] = make([]string, 0)
	}

	// TODO this is also poor programming...
	FollowSet["<systemgoal>"] = []string { "" }

	for i := 0; i < 3; i++ {
	for p := range g.productions {
		rhs := stripRhs(p)
		lhs := stripLhs(p)
		a := stripNonTerminals(rhs)

		for _, B := range a {
			next := nextSymbol(rhs, B)
			t := remove(ComputeFirst(next), "")

			// fmt.Printf("rhs %s, t %v, next %s\n", rhs, t, next)
			// for _, s := range t {
			// 	if x, _ := contains(FollowSet[B], s); !x {
					FollowSet[B] = append(FollowSet[B], t...)
			// 	}
			// }
			
			first := ComputeFirst(next)

			// for _, s := range first {
				if b, _ := contains(first, ""); b {
					FollowSet[B] = append(FollowSet[B], FollowSet[lhs]...)
				}
			// }
		}
	}
}

	printSet(FollowSet)
}

// Checks to see if a string exists in an array of strings
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

// Removes a string from an array of strings
func remove(a []string, s string) []string {

	if b, i := contains(a, s); b {
		copy(a[i:], a[i+1:])
		a = a[:len(a) - 1]
	}

	return a
}

// Pull the vocabulary from a grammar
func pullVocabulary (g Grammar) (v MarkedVocabulary) {
	v = make(map[string]bool, 0)

	for k, _ := range g.nonterminals {
		v[k] = g.nonterminals[k]
	}

	for k, _ := range g.terminals {
		v[k] = g.terminals[k]
	}

	return
}

// Helper function to print a map
func printMap(m map[string]bool) {
	var keys []string

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("key '%s', value '%t'\n", k, m[k])
	}
}

// Helper function to print a map
func printSet(m map[string][]string) {
	var keys []string

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("key '%s', value '%v'\n", k, m[k])
	}
}

// Grabs the next symbol
func nextSymbol(s, v string) string {

	strs := strings.Fields(s)

	if b, i := contains(strs, v); b {
		if i + 1 < len(strs) {
			return strs[i + 1]
		}
	}

	return ""
}
