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

func (p *Parser) ProcessId(s SemanticRecord) {
	s.exprRec.Kind = IdExpr
	s.exprRec.Name = ss[currentIndex - 1].exprRec.Name

	p.checkId(p.extract(s))
	s.exprRec.Kind = IdExpr
	s.exprRec.Name = p.extract(s)

	ss[leftIndex] = s
}

func (p *Parser) ProcessLiteral(s SemanticRecord) {
	s.exprRec.Kind = LiteralExpr
	s.exprRec.Name = ss[currentIndex - 1].exprRec.Name
	s.exprRec.Val, _ = strconv.Atoi(ss[currentIndex - 1].exprRec.Name)
	ss[leftIndex] = s
}

func (p *Parser) ProcessOp(s SemanticRecord) {
	if ss[currentIndex - 1].exprRec.Name == "+" {
		ss[currentIndex - 1].exprRec.Kind = NotExpr
		ss[currentIndex - 1].opRec.Op = PlusOp
		s.opRec.Op = PlusOp
		s.exprRec.Name = "+"
	} else {
		ss[currentIndex - 1].exprRec.Kind = NotExpr
		ss[currentIndex - 1].opRec.Op = MinusOp
		s.opRec.Op = MinusOp
		s.exprRec.Name = "-"
	}

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
	switch s.exprRec.Kind {
	case IdExpr:
		fallthrough
	case TempExpr:
		return s.exprRec.Name
	case LiteralExpr:
		return strconv.Itoa(s.exprRec.Val)
	default:
		return p.extractOp(s.opRec)
	}
}

// Determine's the type of operation
func (p *Parser) extractOp(o OpRec) string {
	switch o.Op {
	case PlusOp:
		return "ADD"
	case MinusOp:
		return "SUB"
	default:
		panic(fmt.Errorf("Incorrect operation, found %v, must be MinusOp or PlusOp", o.Op))
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
