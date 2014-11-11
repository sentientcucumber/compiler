// Author: Michael Hunsinger
// Date:   Oct 18 2014
// File:   semantic.go
// Definitions and implementation for semantic routines

package compiler

import (
	"strconv"
	"bytes"
	"fmt"
)

// Write the snippet of code to store the variable
func (p *Parser) Assign(target, src SemanticRecord) {
	p.generate("STORE", p.extract(src), target.exprRec.Name)
}

// Write the snippet of code to read the variable
func (p *Parser) ReadId(in SemanticRecord) {
	p.generate("READ", in.exprRec.Name, "INTEGER")
}

// Write the snippet of code to write the variable
func (p *Parser) WriteExpr(out SemanticRecord) {
	p.generate("WRITE", p.extract(out), "INTEGER")
}

// Create an assembly language-esque infix
func (p *Parser) GenInfix(e1, op, e2 SemanticRecord, out *SemanticRecord) {
	s := new(SemanticRecord)
	s.exprRec.Kind = TempExpr
	s.exprRec.Name = p.getTemp()

	p.generate(p.extract(op), p.extract(e1), p.extract(e2), s.exprRec.Name)
	out = s
}

// TODO placeholder for process id function
func (p *Parser) ProcessId(s SemanticRecord) {
	s.exprRec.Kind = IdExpr
	s.exprRec.Name = ss[currentIndex - 1].exprRec.Name

	p.checkId(p.extract(s))
	s.exprRec.Kind = IdExpr
	s.exprRec.Name = p.extract(s)

	ss[leftIndex] = s
}

// TODO placeholder for process literal function
func (p *Parser) ProcessLiteral(s SemanticRecord) {
	s.exprRec.Kind = LiteralExpr
	s.exprRec.Name = ss[currentIndex - 1].exprRec.Name

	s.exprRec.Kind = LiteralExpr
	s.exprRec.Val, _ = strconv.Atoi(p.extract(s))

	ss[leftIndex] = s
}

// TODO placeholder for process op function
func (p *Parser) ProcessOp(s SemanticRecord) {
	s.opRec.Op = ss[currentIndex - 1].opRec.Op
	ss[leftIndex] = s
}

// Used to do the actual writing of code
func (p *Parser) generate(strs ...string) {
	var buf bytes.Buffer

	for i, s := range strs {
		buf.WriteString(s)
		if i == 0 && len(strs) > 1 {
			buf.WriteString(" ")
		} else if i < len(strs) - 1 {
			buf.WriteString(", ")
		} else {
			buf.WriteString("\n")
		}
	}

	p.Writer.Write(buf.Bytes())
	p.Writer.Flush()
}

// Extract various parts of the SemanticRecord
func (p *Parser) extract(s SemanticRecord) string {
	// if Kind is empty, then semantic record must be an op rec
	// otherwise exprRec
	if len(s.exprRec.Name) == 0 {
		return p.extractOp(s.opRec)
	} else {
		if s.exprRec.Kind == IdExpr ||
			s.exprRec.Kind == TempExpr {
			return s.exprRec.Name
		} else {
			return strconv.Itoa(s.exprRec.Val)
		}
	}
}

// Determine's the type of operation
func (p *Parser) extractOp(o OpRec) string {
	if o.Op == PlusOp {
		return "ADD"
	} else {
		return "SUB"
	}
}

// func (p *Parser) semanticCopy(src, dest SemanticRecord) {
func (p *Parser) semanticCopy(src, dest int) {
	ss[dest] = ss[src]
}

// Checks to see if the symbol already exists
func (p *Parser) lookUp(s string) bool {
	for _, sym := range p.SymbolTable {
		if sym == s {
			return true
		}
	}
	return false
}
// Checks to see if the Id being passed is in the symbol table
func (p *Parser) checkId(s string) {
	if !p.lookUp(s) {
		p.enter(s)
		p.generate("DECLARE", s, "INTEGER")
	}
}
// Enters the symbol into the symbol table, so long as there's room
func (p *Parser) enter(s string) {
	if p.lastSymbol < p.maxSymbol {
		p.SymbolTable[p.lastSymbol] = s
		p.lastSymbol++
	} else {
		panic(fmt.Errorf("Symbol table overflow"))
	}
}

// Returns the current Temp name
func (p *Parser) getTemp() string {
	var buf bytes.Buffer

	p.maxTemp++
	buf.WriteString("Temp&")
	buf.WriteString(strconv.Itoa(p.maxTemp))
	p.checkId(buf.String())

	return buf.String()
}

func (p *Parser) Start() {
	p.maxSymbol = 100
	p.lastSymbol = 0
	p.maxTemp = 0
	p.SymbolTable = make([]string, p.maxSymbol)
}

// generates the halt function
func (p *Parser) finish() {
	p.generate("HALT")
}
