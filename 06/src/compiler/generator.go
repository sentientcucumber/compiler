// Author: Michael Hunsinger
// Date:   Oct 4 2014
// File:   generator.go
// Implementation of a predict generator for LL(1) grammars

package compiler

import (
	"fmt"
	"strings"
	"sort"
	"regexp"
)

type Generator struct {
	grammar           Grammar
	FirstSet          map[string][]string
	FollowSet         map[string][]string
	derivesLambda     MarkedVocabulary
}

// Generates a predict set
func (g *Generator) Predict() {
	g.MarkLambda(g.grammar)

	g.FirstSet = make(map[string][]string, 0)
	g.FollowSet = make(map[string][]string, 0)
	
	g.FillFirstSet()
	g.FillFollowSet()

	// fmt.Printf("FirstSet ----------------------------------\n")
	// printSet(g.FirstSet)
	// fmt.Printf("FollowSet ----------------------------------\n")
	// printSet(g.FollowSet)

	for p := range g.grammar.productions {
		predictSet := make([]string, 0)
		rhs := stripRhs(p)
		lhs := stripLhs(p)

		// Skip over where rhs is empty
		strs := strings.Fields(rhs)
		term := false
		fmt.Printf("First ( %s )", rhs)

		for i := 0; i < len(strs) && !term; i++ {

			// If the first symbol is a terminal, add it on, it's the
			// predict set, otherwise, find the first set for the nonterminal
			if (regexp.MustCompile("[[:punct:]]\\s").MatchString(strs[i]) ||
				!regexp.MustCompile("<[a-zA-Z\\s]*>").MatchString(strs[i])) &&
				lhs == "lambda" {
				predictSet = append(predictSet, strs[i])
				term = true

			} else {
				predictSet = append(predictSet, g.FirstSet[strs[i]]...)

				// This should be safe in this nonterminal branch, as
				// terminals will never result in lambda
				if b, _ := contains(g.FirstSet[strs[i]], "lambda"); b {
					predictSet = remove(predictSet, "lambda")

					// This bit ensures the set contains unique values
					for _, v := range g.FollowSet[lhs] {
						if c, _ := contains(predictSet, v); !c && v != "lambda" {
							predictSet = append(predictSet, v)
						}
					}

					fmt.Printf(" ∪ Follow ( %s ) - λ", lhs)
				}

				term = containsTerminal(predictSet)
			}
		}
		
		fmt.Printf(" = %s\n", predictSet)
	}
}

// Mark which parts of a vocabulary (terminals and nonterminals) from a grammar
// can produce lambda. If reading an LL(1) grammar, the grammar should be
// formatted that the LHS produces nothing instead of nil or a lambda unicode
// (e.g. "<lhs> -> ")
func (g *Generator) MarkLambda (gmr Grammar) MarkedVocabulary {
	g.grammar = gmr
	changes := true
	g.derivesLambda = pullVocabulary(gmr)
	
	for k, _ := range g.derivesLambda {
		g.derivesLambda[k] = false
	}

	for changes {
		changes = false

		for p := range gmr.productions {
			rhsDerivesLambda := true
			rhs := stripRhs(p)
			
			for _, s := range stripSymbols(rhs) {
				rhsDerivesLambda = rhsDerivesLambda && g.derivesLambda[s];
			}

			lhs := stripLhs(p)
			if rhsDerivesLambda && !g.derivesLambda[lhs] {
				changes = true
				g.derivesLambda[lhs] = true
			}
		}
	}

	return g.derivesLambda
}

// Determines the first terminal or lambda for a given set of symbols,
// terminals and nonterminals
func (g *Generator) ComputeFirst (s string) (result TermSet) {
	strs := strings.Fields(s)

	if k := len(strs); k == 0 {
		result = append(result, "lambda")
	} else {

		if b, _ := contains(g.FirstSet[strs[0]], "lambda"); !b {
			tmp := remove(g.FirstSet[strs[0]], "lambda")
			result = append(result, tmp...) // g.FirstSet[strs[0]]...)
		} else {
			i := 0
			// tmp := remove(g.FirstSet[strs[i]], "lambda")
			// result = append(result, tmp...)

			for b, _ := contains(g.FirstSet[strs[0]], "lambda");
			!b && i < k - 1;
			b, _ = contains(g.FirstSet[strs[i]], "lambda") {
				tmp := remove(g.FirstSet[strs[i]], "lambda")
				result = append(result, tmp...) //g.FirstSet[strs[i]]...)
			}

			// never reaches this point... 
			if b, _ := contains(g.FirstSet[strs[k - 1]], "lambda"); b && i == k - 1 {
				result = append(result, "lambda")
			}
		}
	}
	
	return
}


// Fill the FirstSet
func (g *Generator) FillFirstSet() {
	for A := range g.grammar.nonterminals {
		if g.derivesLambda[A] {
			g.FirstSet[A] = []string { "lambda" }
		} else {
			g.FirstSet[A] = make([]string, 0)
		}
	}

	for a := range g.grammar.terminals {
		g.FirstSet[a] = []string { a }

		for A := range g.grammar.nonterminals {
			for p := range g.grammar.productions {
				rhs := stripRhs(p)
				lhs := stripLhs(p)

				// Extra bit to make sure this is a "set"
				if _, s := firstTerm(rhs); s == a && lhs == A {
					if b, _ := contains(g.FirstSet[A], a); !b {
						g.FirstSet[A] = append(g.FirstSet[A], a);
					}
				}
			}
		}
	}

	// TODO this is poor programming... 
	for i := 0; i < 2; i++ {
		for p := range g.grammar.productions {
			lhs := stripLhs(p)
			rhs := stripRhs(p)
			first := g.ComputeFirst(rhs)

			// Extra bit to make sure this is a "set"
			for i, _ := range first {
				if b, _ := contains(g.FirstSet[lhs], first[i]); !b {
					g.FirstSet[lhs] = append(g.FirstSet[lhs], first[i])
				}
			}
		}
	}
}

// Fill the FollowSet
func (g *Generator) FillFollowSet() {
	for A := range g.grammar.nonterminals {
		g.FollowSet[A] = make([]string, 0)
	}

	// TODO this is also poor programming...
	g.FollowSet["<S>"] = []string { "lambda" }

	for i := 0; i < 2; i++ {
		for p := range g.grammar.productions {
			rhs := stripRhs(p)
			lhs := stripLhs(p)
			a := stripNonTerminals(rhs)
			
			for _, B := range a {
				next := nextSymbol(rhs, B)
				g.FollowSet[B] = append(g.FollowSet[B], g.FirstSet[next]...)

				if b, _ := contains(g.ComputeFirst(next), "lambda");
				b || len(g.FirstSet[next]) == 0 {

					// This bit added to maintain the set contains unique values
					for _, v := range g.FollowSet[lhs] {
						if c, _ := contains(g.FollowSet[B], v); !c {
							g.FollowSet[B] = append(g.FollowSet[B], v)
						}
					}
				}
			}
		}
	}
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

	n := a
	if b, i := contains(n, s); b {
		copy(n[i:], n[i+1:])
		n = n[:len(n) - 1]
	}

	return n
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

	return "lambda"
}

// Determines if a set has a terminal symbol
func containsTerminal(set []string) (found bool) {

	for _, v := range set {
		if !regexp.MustCompile("<[a-zA-Z\\s]*>").MatchString(v) {
			found = true
			return
		}
	}

	return
}
