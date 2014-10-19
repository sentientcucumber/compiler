// Author: Michael Hunsinger
// Date:   Oct 11 2014
// File:   generator.go
// Definition and implementation for the table struct

package compiler

import (
	"sort"
	"fmt"
	"strings"
)

type Table struct {
	production  map[int]string
	rowTitle    []string
	colTitle    []string
	rowCount    int
	colCount    int
	array       [][]int
}

// Sets up production, rowTitle, colTitle, rowCount, colCount
func (t *Table) init(g Grammar) {
	t.rowCount = len(g.nonterminals)

	if _, c := g.terminals[lambda.name]; c {
		t.colCount = len(g.terminals) - 1
	} else {
		t.colCount = len(g.terminals)
	}

	t.rowTitle = make([]string, t.rowCount)
	t.colTitle = make([]string, t.colCount)

	i := 0
	for k := range g.nonterminals {
		t.rowTitle[i] = k
		i++
	}

	i = 0
	for k := range g.terminals {
		if k != lambda.name {
			t.colTitle[i] = k
			i++
		}
	}

	sort.Strings(t.colTitle)
	sort.Strings(t.rowTitle)

	t.production = g.staticProd
}

// Prints out the table
func (t Table) print() {
	
	fmt.Printf("\t\t")
	for _, v := range t.colTitle {
		fmt.Printf("%s\t", v)
	}

	fmt.Printf("\n")

	for i := range t.array {
		fmt.Printf("%s\t", t.rowTitle[i])
		if len(t.rowTitle[i]) < 8 {
			fmt.Printf("\t")
		}

		for j := range t.array[1] {
			fmt.Printf("%d\t", t.array[i][j])
		}
		fmt.Printf("\n")
	}
}

// Performs a lookup based on a terminal and nonterminal symbol and returns the
// production number
func (t *Table) lookup(n, x Symbol, g *Generator) int {
	var p int
	l := false
	c := 0

	for i, v := range t.production {
		lhs := stripLhs(v)
		rhs := stripRhs(v)
		strs := strings.Fields(rhs)

		// fmt.Printf("%d %s", i, v)
		// If the first symbol on RHS is a terminal, see that it matches and
		// return the production. Otherwise, increment production counter
		// if there's only one, it must be this production for all terminals
		if lhs == n.name {
			if strs[0] == x.name {
				return i
			} else if strs[0] == lambda.name {
				p = i
				l = true
			} else if !isTerminal(strs[0], v) {
				for _, j := range g.computeFirst(strs[0]) {
					if j.name == x.name {
						// p = i
						return i
					}
				}
			} else if !l {
				c++
				p = i
			}
		}
	}

	if c == 1 {
		return p
	}

	if l {
		return p
	}

	return 0
}
