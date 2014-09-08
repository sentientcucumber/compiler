// Author: Michael Hunsinger
// Date:   Aug 30 2014
// File:   parser.go
// Definiton for a parser.

package compiler

import (
   "fmt"
)

type Parser struct {
   Scanner Scanner
   currState string
}

// Function matches a the legal, expected token (according to the grammar) to
// the token read in the file, throw an error if not
func (p *Parser) Match (legalTok Token) {
   currTok := p.Scanner.Scan()

   if currTok != legalTok {
      panic(fmt.Errorf("Syntax error when reading %v, doesn't match %v",
         currTok, legalTok))
   }

   fmt.Printf("Legal token: %v, passed token: %v\n", legalTok, currTok)
}

// Returns the next token, but does not advance the Scanner
func (p *Parser) NextToken () Token {
   off, _ := p.Scanner.Reader.Seek(0, 1)
   next := p.Scanner.Scan()
   newoff, _ := p.Scanner.Reader.Seek(0, 1)
   p.Scanner.Reader.Seek(off - newoff, 1)

   return next
}

// SystemGoal definition according to grammar
func (p *Parser) SystemGoal() {
   fmt.Println("Parsing system goal")
   p.Program()
}

// Program definition according to grammar
func (p *Parser) Program() {
   fmt.Println("Parsing program")
   p.Match(BeginSym)
   p.StatementList()
   p.Match(EndSym)
}

// StatementList definition according to grammar
func (p *Parser) StatementList() {
   fmt.Println("Parsing statement list")
   p.Statement()

   if next := p.NextToken();
   next == Id || next == ReadSym || next == WriteSym {
      p.StatementList()
   }
}

// Statement definition according to grammar
func (p *Parser) Statement() {
   fmt.Println("Parsing statement")
   next := p.NextToken()

   switch next {
   case Id:
      p.Ident()
      p.Match(AssignOp)
      p.Expression()
      p.Match(SemiColon)
      break

   case ReadSym:
      p.Match(ReadSym)
      p.Match(LParen)
      p.IdList()
      p.Match(RParen)
      p.Match(SemiColon)
      break

   case WriteSym:
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
   fmt.Println("Parsing IdList")
   p.Ident()
   
   if next := p.NextToken(); next == Comma {
      p.Match(Comma)
      p.IdList()
   }
}

// ExpressionList definition according to grammar
func (p *Parser) ExprList() {
   fmt.Println("Parsing ExprList")
   p.Expression()
   
   if next := p.NextToken(); next == Comma {
      p.Match(Comma)
      p.ExprList()
   }

   fmt.Printf("\n")
}

// Expression definition according to grammar
func (p *Parser) Expression() {
   fmt.Println("Parsing Expression")
   p.Primary()
   next := p.NextToken()

   if next == PlusOp || next == MinusOp {
      p.AddOp()
      p.Expression()
   }
}

// Primary definition according to grammar
func (p *Parser) Primary() {
   fmt.Println("Parsing Primary")
   next := p.NextToken()

   switch next {
   case LParen:
      p.Match(LParen)
      p.Expression()
      p.Match(RParen)
      break

   case Id:
      p.Ident()
      break

   case IntLiteral:
      p.Match(IntLiteral)
      break

   default:
      panic(fmt.Errorf("Syntax error when reading %v\n", next))
      break
   }
}

// Identifier definition according to grammar
func (p *Parser) Ident() {
   fmt.Println("Parsing Ident")
   p.Match(Id)
}

// AddOp definition according to grammar
func (p *Parser) AddOp() {
   fmt.Println("Parsing AddOp")
   next := p.NextToken()

   switch next {
   case PlusOp:
      p.Match(PlusOp)
      break

   case MinusOp:
      p.Match(MinusOp)
      break

   default:
      panic(fmt.Errorf("Syntax error when reading %v\n", next))
      break
   }
}
