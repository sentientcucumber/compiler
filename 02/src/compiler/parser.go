// Author: Michael Hunsinger
// Date:   Aug 30 2014
// File:   parser.go
// Definiton for a parser.

package compiler

import (
	"fmt"
	"strings"
)

type Parser struct {
	Scanner   Scanner
	currState string
}

// Function matches a the legal, expected token (according to the grammar) to
// the token read in the file, throw an error if not
func (p *Parser) Match(legalTok Token) {
	currTok := p.Scanner.Scan()

	if currTok != legalTok {
		panic(fmt.Errorf("Syntax error when reading %v, doesn't match %v",
			currTok, legalTok))
	}
}

// Returns the next token, but does not advance the Scanner
func (p *Parser) NextToken() Token {
	off, _ := p.Scanner.Reader.Seek(0, 1)
	next := p.Scanner.Scan()
	newoff, _ := p.Scanner.Reader.Seek(0, 1)
	p.Scanner.Reader.Seek(off-newoff, 1)

	return next
}

// Returns the next n token ahead, but does not advance the Scanner
// Primarily used for printing parsing behavior, doesn't affect actual parsing
func (p *Parser) ReadNTokensAhead(n int) Token {
	var next Token

	off, _ := p.Scanner.Reader.Seek(0, 1)
	for i := 0; i < n; i++ {
		next = p.Scanner.Scan()
	}

	newoff, _ := p.Scanner.Reader.Seek(0, 1)
	p.Scanner.Reader.Seek(off-newoff, 1)

	return next
}

// Returns whether or not the token exists and if so, how many
// Primarily used for printing parsing behavior, doesn't affect actual parsing
func (p *Parser) HasToken(token Token) (found bool, count int) {
	off, _ := p.Scanner.Reader.Seek(0, 1)
	count = 0

	for t := p.Scanner.Scan(); t != EofSym; t = p.Scanner.Scan() {

		if t == token {
			count++
			found = true
		}
	}

	newoff, _ := p.Scanner.Reader.Seek(0, 1)
	p.Scanner.Reader.Seek(off-newoff, 1)

	return
}

// SystemGoal definition according to grammar
func (p *Parser) SystemGoal() {
	p.currState = "<system goal> --> BeginSym <program> EndSym EofSym"
	fmt.Println(p.currState)

	p.Program()
}

// Program definition according to grammar
func (p *Parser) Program() {
	p.currState = strings.Replace(p.currState, "<program>", "<statement list>", 1)
	fmt.Println(p.currState)

	p.Match(BeginSym)
	p.StatementList()
	p.Match(EndSym)
}

// StatementList definition according to grammar
func (p *Parser) StatementList() {

	if t, count := p.HasToken(SemiColon); t && count > 1 {
		p.currState = strings.Replace(p.currState, "<statement list>", "<statement> <statement list>", 1)
		fmt.Println(p.currState)
	} else {
		p.currState = strings.Replace(p.currState, "<statement list>", "<statement>", 1)
		fmt.Println(p.currState)
	}

	p.Statement()

	if next := p.NextToken(); next == Id || next == ReadSym || next == WriteSym {
		p.StatementList()
	}
}

// Statement definition according to grammar
func (p *Parser) Statement() {

	next := p.NextToken()

	switch next {
	case Id:
		p.currState = strings.Replace(p.currState, "<statement>", "<ident> := <expression> ;", 1)
		fmt.Println(p.currState)

		p.Ident()
		p.Match(AssignOp)
		p.Expression()
		p.Match(SemiColon)
		break

	case ReadSym:
		p.currState = strings.Replace(p.currState, "<statement>", "ReadSym ( <id list> ) ;", 1)
		fmt.Println(p.currState)

		p.Match(ReadSym)
		p.Match(LParen)
		p.IdList()
		p.Match(RParen)
		p.Match(SemiColon)
		break

	case WriteSym:
		p.currState = strings.Replace(p.currState, "<statement>", "WriteSym ( <expr list> ) ;", 1)
		fmt.Println(p.currState)

		p.Match(WriteSym)
		p.Match(LParen)
		p.ExprList()
		p.Match(RParen)
		p.Match(SemiColon)
		break

	default:
		panic(fmt.Errorf("Syntax error when reading %v\n", next))
		break
	}
}

// IdList definition according to grammar
func (p *Parser) IdList() {
	if ahead := p.ReadNTokensAhead(2); ahead == Comma {
		p.currState = strings.Replace(p.currState, "<id list>", "<ident>, <id list>", 1)
		fmt.Println(p.currState)
	} else {
		p.currState = strings.Replace(p.currState, "<id list>", "<ident>", 1)
		fmt.Println(p.currState)
	}

	p.Ident()

	if next := p.NextToken(); next == Comma {
		p.Match(Comma)
		p.IdList()
	}
}

// ExpressionList definition according to grammar
func (p *Parser) ExprList() {
	if ahead := p.ReadNTokensAhead(2); ahead == Comma {
		p.currState = strings.Replace(p.currState, "<expr list>", "<expression>, <expr list>", 1)
		fmt.Println(p.currState)
	} else {
		p.currState = strings.Replace(p.currState, "<expr list>", "<expression>", 1)
		fmt.Println(p.currState)
	}

	p.Expression()

	if next := p.NextToken(); next == Comma {
		p.Match(Comma)
		p.ExprList()
	}

	fmt.Printf("\n")
}

// Expression definition according to grammar
func (p *Parser) Expression() {

	// some issues here when expression is contained in parens
	if ahead := p.ReadNTokensAhead(2); ahead == PlusOp || ahead == MinusOp {
		p.currState = strings.Replace(p.currState, "<expression>", "<primary> <add op> <expression>", 1)
		fmt.Println(p.currState)
	} else {
		p.currState = strings.Replace(p.currState, "<expression>", "<primary>", 1)
		fmt.Println(p.currState)
	}

	p.Primary()
	next := p.NextToken()

	if next == PlusOp || next == MinusOp {
		p.AddOp()
		p.Expression()
	}
}

// Primary definition according to grammar
func (p *Parser) Primary() {
	next := p.NextToken()

	switch next {
	case LParen:
		p.currState = strings.Replace(p.currState, "<primary>", "( <expression> )", 1)
		fmt.Println(p.currState)

		p.Match(LParen)
		p.Expression()
		p.Match(RParen)
		break

	case Id:
		p.currState = strings.Replace(p.currState, "<primary>", "<ident>", 1)
		fmt.Println(p.currState)

		p.Ident()
		break

	case IntLiteral:
		p.currState = strings.Replace(p.currState, "<primary>", "IntLiteral", 1)
		fmt.Println(p.currState)

		p.Match(IntLiteral)
		break

	default:
		panic(fmt.Errorf("Syntax error when reading %v\n", next))
		break
	}
}

// Identifier definition according to grammar
func (p *Parser) Ident() {
	p.currState = strings.Replace(p.currState, "<ident>", "Id", 1)
	fmt.Println(p.currState)

	p.Match(Id)
}

// AddOp definition according to grammar
func (p *Parser) AddOp() {
	next := p.NextToken()

	switch next {
	case PlusOp:
		p.currState = strings.Replace(p.currState, "<add op>", "PlusOp", 1)
		fmt.Println(p.currState)

		p.Match(PlusOp)
		break

	case MinusOp:
		p.currState = strings.Replace(p.currState, "<add op>", "MinusOp", 1)
		fmt.Println(p.currState)

		p.Match(MinusOp)
		break

	default:
		panic(fmt.Errorf("Syntax error when reading %v\n", next))
		break
	}
}
