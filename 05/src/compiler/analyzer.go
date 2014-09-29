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

	fmt.Printf("productions ----------------\n")
	productions.Print()

	fmt.Printf("nonterminals ---------------\n")
	nonterminals.Print()

	fmt.Printf("nonterminals ---------------\n")
	terminals.Print()
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
	
	return nil
}

// Reads all the nonterminals in a buffer
func readNonterminals (buf bytes.Buffer) {

	s := strings.Replace(buf.String(), "->", "", 1)

	for strings.Contains(s, "<") {
		start := strings.Index(s, "<")
		end   := strings.Index(s, ">") + 1

		nonterminals.Add(s[start:end])
		s = s[end:]
	}
}

// Reads all the terminals in a buffer
func readTerminals (buf bytes.Buffer) {

	// remove terminal symbols and arrow
	re := regexp.MustCompile("(?:\\<[a-zA-Z0-9 ]*\\>|->|\\|)")
	s := re.ReplaceAllString(buf.String(), " ")

	strs := strings.Fields(s)

	for _, e := range strs {
		terminals.Add(e)
	}
}

func (a *Analyzer) readRHS() {

}

func (a *Analyzer) readLHS() {

}

func addUnique() {
	
}

