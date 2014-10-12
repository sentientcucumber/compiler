// Author: Michael Hunsinger
// Date:   Oct 11 2014
// File:   table.go
// Definition and procedures for a table

package compiler

import (
	"sort"
	"fmt"
)

// Definition of a Table
type Table [][]uint8

// Initializes a table with 0s, and sets up the column and row headers
func (t Table) initTable (g Grammar) {
	rows := len(g.nonterminals)
	cols := len(g.terminals)

	// Initialize the table
	t = make([][]uint8, rows)
	for i := range t {
		t[i] = make([]uint8, cols)
	}

	rowTitle := make([]string, rows)
	colTitle := make([]string, cols)

	var tKeys []string
	var nKeys []string

	for k := range g.terminals {
		tKeys = append(tKeys, k)
	}

	sort.Strings(tKeys)

	for i, k := range tKeys {
		rowTitle[i] = k
	}

	for k := range g.nonterminals {
		nKeys = append(nKeys, k)
	}

	sort.Strings(nKeys)

	for i, k := range nKeys {
		rowTitle[i] = k
	}

	i := 0
	fmt.Printf("\t\t")
	for _, v := range tKeys {
		colTitle[i] = v
		fmt.Printf("%s\t", colTitle[i])
		i++
	}

	fmt.Printf("\n")
	
	for i := range t {
		fmt.Printf("%s\t", rowTitle[i])

		if (len(rowTitle[i]) < 8) {
			fmt.Printf("\t")
		}

		for j := range t {
			fmt.Printf("%d\t", t[i][j])
		}
		fmt.Printf("\n")
	}
}
