// Author: Michael Hunsinger
// Date:   Oct 4 2014
// File:   analyzer.go
// Implementation of a grammar analyzer

package compiler

import (
	"bytes"
	"strings"
	"regexp"
)

type Analyzer struct {
	Reader  bytes.Reader
}

var (
	terminals     = map[string]bool {}
	nonterminals  = map[string]bool {}
	productions   = map[string]bool {}
	rhs           = map[string]bool {}
	lhs           = map[string]bool {}
)

// Reads whatever grammar is passed in.
func (a *Analyzer) ReadGrammar() Grammar {
	for err := a.readProduction(); err == nil; err = a.readProduction() {
		a.readProduction()
	}

	return Grammar {
		terminals:     terminals,
		nonterminals:  nonterminals,
		productions:   productions,
		rhs:           rhs,
		lhs:           lhs,
	}
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

	productions[buf.String()] = true

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

		nonterminals[s[start:end]] = true
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

		terminals[i] = true
	}
}

// Reads the RHS of each production, notice this passes them into a set so any
// repetitions will be ignored
func readRHS (buf bytes.Buffer) {

	r := strings.Split(buf.String(), "->")
	rhs[strings.TrimSpace(r[1])] = true
}

// Reads the LHS of each production, notice this passes them into a set, so any
// repetitions will be ignored
func readLHS(buf bytes.Buffer) {

	l := strings.Split(buf.String(), "->")
	lhs[strings.TrimSpace(l[0])] = true
}

// Pull the RHS from a string representation of a string
func stripRhs (s string) string {

	strs := strings.Split(s, "->")
	return strings.TrimSpace(strs[1])
}

// Pull the LHS from a string representation of a string
func stripLhs (s string) string {

	strs := strings.Split(s, "->")
	return strings.TrimSpace(strs[0])
}

// Pull the symbols from a string representation
func stripSymbols (s string) []string {

	strArr := make([]string, 0)
	nonStr := strings.Replace(s, "->", "", 1)
	
	for strings.Contains(nonStr, "<") {
		start := strings.Index(nonStr, "<")
		end   := strings.Index(nonStr, ">") + 1

		strArr = append(strArr, nonStr[start:end])
		nonStr = nonStr[end:]
	}

	re := regexp.MustCompile("(?:\\<[a-zA-Z0-9 ]*\\>|->|\\|)")
	tStr := re.ReplaceAllString(s, " ")

	strs := strings.Fields(tStr)
	
	for _, i := range strs {
		strArr = append(strArr, i)
	}

	return strArr
}

// Returns the next symbol, whether terminal or nonterminal
func firstTerm (s string) string {

	strs := strings.Fields(s)
	if len(strs) == 0 {
		return ""
	} else {
		return strs[0]
	}
}
