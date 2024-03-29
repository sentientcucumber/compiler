// Author: Michael Hunsinger
// Date:   Oct 11 2014
// File:   generator.go
// Implementation of a predict generator for LL(1) grammars

package compiler

import (
	"fmt"
	"strings"
	"regexp"
)

type Generator struct {
	Grammar           Grammar
	FirstSet          Set
	FollowSet         Set
	PredictSet        Set
	predictNonTerm    Set
	derivesLambda     MarkedVocabulary
}

// Used throughout the program, this should be const, but can't do so
var lambda = Symbol { "λ", "LAMBDA" }

// Generates a predict set
func (g *Generator) Predict() {
	// Consume the Grammar
	g.markLambda(g.Grammar)

	// Initialize sets
	g.FirstSet   = make(map[string][]Symbol, 0)
	g.FollowSet  = make(map[string][]Symbol, 0)
	g.PredictSet = make(map[string][]Symbol, 0)
	g.predictNonTerm = make(map[string][]Symbol, 0)

	g.FillFirstSet()
	g.FillFollowSet()

	for p := range g.Grammar.productions {
		rhs := stripRhs(p)
		lhs := stripLhs(p)

		// Skip over where rhs is empty
		strs := strings.Fields(rhs)
		term := false
		fmt.Printf("First ( '%s' )", rhs)

		for i := 0; i < len(strs) && !term; i++ {

			// If the first symbol is a terminal, add it on, it's the
			// predict set, otherwise, find the first set for the nonterminal
			if isTermial(strs[i], lhs) {
				g.PredictSet.add(strs[i], Symbol { strs[i], "TERMINAL"})
				g.predictNonTerm.add(lhs, Symbol {strs[i], "TERMINAL"})
				term = true
			} else {
				for _, v := range g.FirstSet[strs[i]] {
					g.PredictSet.add(strs[i], v)
					g.predictNonTerm.add(lhs, v)
				}

				// This should be safe in this nonterminal branch, as
				// terminals will never result in lambda
				if b, _ := g.FirstSet.containsLambda(strs[i]); b {
					g.PredictSet.removeLambda(strs[i])
					g.predictNonTerm.removeLambda(lhs)

					for _, v := range g.FollowSet[lhs] {

						// Used to keep the various lambdas in line
						if v.name != lambda.name {
							temp := []string { lambda.name, lhs }
							g.PredictSet.add(strings.Join(temp, " "), v)
							g.predictNonTerm.add(lhs, v)
						}
					}
					
					fmt.Printf(" ∪ Follow ( %s ) - λ", lhs)
				}

				term = true
			}

			fmt.Printf(" = ")

			// This looks up the correct name since we stored lambda based on
			// their non-terminal name (e.g. "λ <expressiontail>")
			if strs[i] == lambda.name {
				temp := []string { lambda.name, lhs }
				strs[i] = strings.Join(temp, " ")
			}

			for _, v := range g.PredictSet[strs[i]] {
				fmt.Printf("%s ", v.name)
			}

			fmt.Printf("\n")
		}
	}
	fmt.Printf("\nPredict set\n")
	g.predictNonTerm.print()
}

// Mark which parts of a vocabulary (terminals and nonterminals) from a Grammar
// can produce lambda. If reading an LL(1) Grammar, the Grammar should be
// formatted that the LHS produces nothing instead of nil or a lambda unicode
// (e.g. "<lhs> -> ")
func (g *Generator) markLambda (gmr Grammar) MarkedVocabulary {
	g.Grammar = gmr
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
		result = append(result, lambda)
	} else {

		if b, _ := g.FirstSet.containsLambda(strs[0]); !b {
			temp := g.FirstSet.removeLambda(strs[0])
			result = append(result, temp[strs[0]]...) // g.FirstSet[strs[0]]...)

		} else {
			i := 0

			b, _ := g.FirstSet.containsLambda(strs[0])

			for !b && i < k - 1 {
				temp := g.FirstSet.removeLambda(strs[i])
				result = append(result, temp[strs[i]]...)
				b, _ = g.FirstSet.containsLambda(strs[0])
			}

			if b, _ := g.FirstSet.containsLambda(strs[k - 1]); b && i == k - 1 {
				result = append(result, lambda)
			}
		}
	}
	
	return
}


// Fill the FirstSet
func (g *Generator) FillFirstSet() {
	for A := range g.Grammar.nonterminals {
		if g.derivesLambda[A] {
			g.FirstSet[A] = []Symbol { lambda }
		} else {
			g.FirstSet[A] = make([]Symbol, 0)
		}
	}

	for a := range g.Grammar.terminals {
		g.FirstSet[a] = []Symbol { Symbol { a, "TERMINAL" } }

		for A := range g.Grammar.nonterminals {
			for p := range g.Grammar.productions {
				rhs := stripRhs(p)
				lhs := stripLhs(p)

				// Added bit of logic to ensure SymbolCategory is correct
				if _, s := firstTerm(rhs); s == a && lhs == A {
					if a == lambda.name {
						g.FirstSet.add(A, lambda)
					} else {
						g.FirstSet.add(A, Symbol { a, "TERMINAL" })
					}
				}
			}
		}
	}

	for i := 0; i < 2; i++ {
		for p := range g.Grammar.productions {
			lhs := stripLhs(p)
			rhs := stripRhs(p)
			first := g.ComputeFirst(rhs)

			for _, v := range first {
				g.FirstSet.add(lhs, v)
			}
		}
	}
}

// Fill the FollowSet
func (g *Generator) FillFollowSet() {
	for A := range g.Grammar.nonterminals {
		g.FollowSet[A] = make([]Symbol, 0)
	}

	start := findStartSymbol(g.Grammar)
	g.FollowSet[start.name] = []Symbol {lambda}

	for i := 0; i < 2; i++ {
		for p := range g.Grammar.productions {
			rhs := stripRhs(p)
			lhs := stripLhs(p)
			a   := stripNonTerminals(rhs)
			
			for _, B := range a {
				next := nextSymbol(rhs, B)
				g.FollowSet.add(B, g.FirstSet[next.name]...)

				if b, _ := g.FirstSet.containsLambda(next.name); b {
					g.FollowSet.add(B, g.FollowSet[lhs]...)
				}
			}
		}
	}
}

// Checks to see if a string exists in an array of strings
func contains(a []Symbol, v Symbol) (found bool, ind int) {
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
func remove(a []Symbol, s Symbol) []Symbol {

	n := a
	if b, i := contains(n, s); b {
		copy(n[i:], n[i+1:])
		n = n[:len(n) - 1]
	}

	return n
}

// Pull the vocabulary from a Grammar
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

// Grabs the next symbol
func nextSymbol(s, v string) Symbol {
	a := make([]Symbol, 0)

	for _, v := range strings.Fields(s) {
		a = append(a, Symbol { name: v })
	}

	if b, i := contains(a, Symbol { name: v }); b {
		if i + 1 < len(a) {
			return a[i + 1]
		}
	}

	return lambda
}

// Determines if the string is terminal or not
func isTermial(s string, l string) bool {
	if (regexp.MustCompile("[[:punct:]]\\s").MatchString(s) ||
		!regexp.MustCompile("<[a-zA-Z\\s]*>").MatchString(s)) &&
		l == "λ" {
		return true
	}

	return false
}

// Determine's the start symbol in a Grammar, must be defined in the Grammar
// passed in (e.g. <Start> -> <nonterminal> $)
func findStartSymbol(g Grammar) Symbol {

	for p := range g.productions {
		if strings.Index(p, "$") > 0 {
			start := stripLhs(p)
			return Symbol { start, "NONTERMINAL"}
		}
	}

	panic(fmt.Errorf("No start symbol defined in the Grammar"))
}
