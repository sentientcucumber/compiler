// Author: Michael Hunsinger
// Date:   Sept 27 2014
// File:   analyzer.go
// Implementation of a grammar analyzer

package compiler

import (
	"fmt"
	"bytes"
	"strings"
	"regexp"
)

type Analyzer struct {
	Reader        bytes.Reader
}

var (
	terminals       = Set { set: map[string]bool {} }
	nonterminals    = Set { set: map[string]bool {} }
	productions     = Set { set: map[string]bool {} }
	rhs             = Set { set: map[string]bool {} }
	lhs             = Set { set: map[string]bool {} }
)
	
// Reads whatever grammar is passed in.
func (a *Analyzer) ReadGrammar() {

	for err := a.readProduction(); err == nil; err = a.readProduction() {
		a.readProduction()
	}

	fmt.Printf("Productions ----------------------------\n")
	productions.Print()

	fmt.Printf("\nSymbols ------------------------------\n")
	terminals.Print()

	fmt.Printf("\nNon-Terminals ------------------------\n")
	nonterminals.Print()

	fmt.Printf("\nTerminals ----------------------------\n")
	terminals.Print()

	fmt.Printf("\nRHS ----------------------------------\n")
	rhs.Print()

	fmt.Printf("\nLHS ----------------------------------\n")
	lhs.Print()
}

// Read each production of the grammar. This assumes that each production is on
// a different line. It will then grab all the terminal, nonterminals, RHS and
// LHS for each production.
func (a *Analyzer) readProduction() error {

	buf := bytes.NewBuffer(*new ([]byte))

	for b, err := a.Reader.ReadByte(); b != '\n'; b, err = a.Reader.ReadByte() {
		buf.WriteByte(b)

		if err != nil {
			return err
		}
	}

	productions.Add(buf.String())

	readNonterminals(*buf)
	readTerminals(*buf)
	readRHS(*buf)
	readLHS(*buf)
	
	return nil
}

// Reads all the nonterminals in a buffer, notice this passes them into a set
// so any repetitions will be ignored
func readNonterminals (buf bytes.Buffer) {

	s := strings.Replace(buf.String(), "->", "", 1)

	for strings.Contains(s, "<") {
		start := strings.Index(s, "<")
		end   := strings.Index(s, ">") + 1

		nonterminals.Add(s[start:end])
		s = s[end:]
	}
}

// Reads all the terminals in a buffer, notice this passes them into a set so
// any repetitions will be ignored
func readTerminals (buf bytes.Buffer) {

	// remove terminal symbols, arrow, and pipe
	re := regexp.MustCompile("(?:\\<[a-zA-Z0-9 ]*\\>|->|\\|)")
	s := re.ReplaceAllString(buf.String(), " ")

	strs := strings.Fields(s)
	
	for _, i := range strs {

		terminals.Add(i)
	}
}

// Reads the RHS of each production, notice this passes them into a set so any
// repetitions will be ignored
func readRHS (buf bytes.Buffer) {

	r := strings.Split(buf.String(), "->")
	rhs.Add(strings.TrimSpace(r[1]))
}

// Reads the LHS of each production, notice this passes them into a set, so any
// repetitions will be ignored
func readLHS(buf bytes.Buffer) {

	l := strings.Split(buf.String(), "->")
	lhs.Add(strings.TrimSpace(l[0]))
}
