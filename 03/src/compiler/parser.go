// Author: Michael Hunsinger
// Date:   Aug 30 2014
// File:   parser.go
// Definiton for a parser.

package compiler

import (
	"fmt"
	"strings"
   "bufio"
   "bytes"
   "strconv"
)

type Parser struct {
	Scanner   Scanner
   Writer    bufio.Writer
	currState string
   maxSymbol int
   maxTemp   int
   symbolTable []string
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
   p.Finish()
}

// Program definition according to grammar
func (p *Parser) Program() {

	p.currState = strings.Replace(p.currState, "<program>", "<statement list>", 1)
	fmt.Println(p.currState)

   p.Start()
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

      var identifier, expr ExprRec

		p.Ident(identifier)
		p.Match(AssignOp)
		p.Expression(expr)
      p.Assign(identifier, expr)
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

   var identifier ExprRec

	p.Ident(identifier)
   p.ReadId(identifier)

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
   var expr ExprRec

	p.Expression(expr)
   p.WriteId(expr)

	if next := p.NextToken(); next == Comma {
		p.Match(Comma)
		p.ExprList()
	}

	fmt.Printf("\n")
}

// Expression definition according to grammar
func (p *Parser) Expression(result ExprRec) {

	// some issues here when expression is contained in parens
	if ahead := p.ReadNTokensAhead(2); ahead == PlusOp || ahead == MinusOp {
		p.currState = strings.Replace(p.currState, "<expression>", "<primary> <add op> <expression>", 1)
		fmt.Println(p.currState)
	} else {
		p.currState = strings.Replace(p.currState, "<expression>", "<primary>", 1)
		fmt.Println(p.currState)
	}

   var leftOperand, rightOperand ExprRec
   var op OpRec

	p.Primary(leftOperand)
	next := p.NextToken()

	if next == PlusOp || next == MinusOp {
		p.AddOp(op)
		p.Expression(rightOperand)
      result = p.GenInfix(leftOperand, op, rightOperand)
	} else {
      result = leftOperand
   }
}

// Primary definition according to grammar
func (p *Parser) Primary(result ExprRec) {

	next := p.NextToken()

	switch next {
	case LParen:
		p.currState = strings.Replace(p.currState, "<primary>", "( <expression> )", 1)
		fmt.Println(p.currState)

		p.Match(LParen)
		p.Expression(result)
		p.Match(RParen)
		break

	case Id:
		p.currState = strings.Replace(p.currState, "<primary>", "<ident>", 1)
		fmt.Println(p.currState)

		p.Ident(result)
		break

	case IntLiteral:
		p.currState = strings.Replace(p.currState, "<primary>", "IntLiteral", 1)
		fmt.Println(p.currState)

		p.Match(IntLiteral)
      p.ProcessLiteral(result)
		break

	default:
		panic(fmt.Errorf("Syntax error when reading %v\n", next))
		break
	}
}

// Identifier definition according to grammar
func (p *Parser) Ident(result ExprRec) {

	p.currState = strings.Replace(p.currState, "<ident>", "Id", 1)
	fmt.Println(p.currState)

	p.Match(Id)
   p.ProcessId(result)
}

// AddOp definition according to grammar
func (p *Parser) AddOp(op OpRec) {

	next := p.NextToken()

	switch next {
	case PlusOp:
		p.currState = strings.Replace(p.currState, "<add op>", "PlusOp", 1)
		fmt.Println(p.currState)

		p.Match(PlusOp)
      p.ProcessOp(op)
		break

	case MinusOp:
		p.currState = strings.Replace(p.currState, "<add op>", "MinusOp", 1)
		fmt.Println(p.currState)

		p.Match(MinusOp)
      p.ProcessOp(op)
		break

	default:
		panic(fmt.Errorf("Syntax error when reading %v\n", next))
		break
	}
}

// Initializes the maxSymbol and maxTemp variables
// These are used for the symbol table and temp variabl assignment
func (p *Parser) Start() {
   p.symbolTable = make([]string, 10)
   p.maxTemp   = 0
}

// Write the snippet of code to store the variable
func (p *Parser) Assign(target, src ExprRec) {
   p.Generate("STORE", p.Extract(src), target.Name)
}

// Write the snippet of code to read the variable
func (p *Parser) ReadId(in ExprRec) {
   p.Generate("READ", in.Name, "INTEGER")
}

// Write the snippet of code to write the variable
func (p *Parser) WriteId(out ExprRec) {
   p.Generate("WRITE", p.Extract(out), "INTEGER")
}

// TODO placeholder for gen infix function
func (p *Parser) GenInfix(e1 ExprRec, op OpRec, e2 ExprRec) ExprRec {
   er := ExprRec { Kind: TempExpr }
   er.Name = p.GetTemp()
   p.Generate(p.ExtractOp(op), p.Extract(e1), p.Extract(e2), er.Name)

   return er
}

// TODO placeholder for process id function
func (p *Parser) ProcessId(er ExprRec) {
   p.CheckId(p.Scanner.tokenBuffer.String())
   er.Kind = IdExpr
   er.Name = p.Scanner.tokenBuffer.String()
}

// TODO placeholder for process literal function
func (p *Parser) ProcessLiteral(er ExprRec) {
   er.Kind = LiteralExpr
   er.Val, _ = strconv.Atoi(p.Scanner.tokenBuffer.String())
}

// TODO placeholder for process op function
func (p *Parser) ProcessOp(o OpRec) {
   
}

// TODO placeholder for finish function
func (p *Parser) Finish() {
   p.Generate("HALT")
}

// Used to do the actual writing of code
func (p *Parser) Generate(strs ...string) {

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

// Extract various parts of the ExprRec
func (p *Parser) Extract(er ExprRec) (val string) {

   kind := er.Kind

   switch kind {
   case IdExpr:
   case TempExpr:
      val = er.Name
      break
   case LiteralExpr:
      val = strconv.Itoa(er.Val)
      break
   // default:
	// 	panic(fmt.Errorf("Trouble extracting the ExprRec %v\n", kind))
   }

   return 
}

// TODO placeholder for extract op function
func (p *Parser) ExtractOp(o OpRec) string {
   if o.Op == PlusOp {
      return "ADD"
   } else {
      return "SUB"
   }
}

// Checks to see if the symbol already exists
func (p *Parser) LookUp(s string) (found bool) {

   found = false
   for _, sym := range p.symbolTable {
      if sym == s {
         found = true
      }
   }

   return
}

// TODO placeholder for check id function
func (p *Parser) CheckId(s string) {

   if !p.LookUp(s) {
      p.Enter(s)
      p.Generate("DECLARE", s, "INTEGER")
   }
}

// TODO placeholder for enter function
func (p *Parser) Enter(s string) {

   if len(p.symbolTable) < p.maxSymbol {
      p.symbolTable[len(p.symbolTable) - 1] = s
   }
}

// TODO placeholder for get temp function
func (p *Parser) GetTemp() (tempName string) {

   p.maxTemp++

   var buf bytes.Buffer
   buf.WriteString("Temp&")
   buf.WriteString(strconv.Itoa(p.maxTemp))

   p.CheckId(buf.String())
   
   return
}
