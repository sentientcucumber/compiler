// Author: Michael Hunsinger
// Date:   Oct 18 2014
// File:   parser.go
// Parser implementation for the compiler

package compiler

import (
	"fmt"
	"bytes"
	"strings"
)

type Parser struct {
	Grammar   Grammar
	Scanner   Scanner
	Table     Table
	Reader    bytes.Reader
}

type Input struct {
	str       []string
	i         int
}

func (p *Parser) Driver() {

	tokenCode := 0
	start := findStartSymbol(p.Grammar)
	stack := new (Stack)
	stack.Push(start)
	p.Scanner.Scan(&tokenCode, bytes.NewBuffer(*new([]byte)))

	state := Input { strings.Fields(p.getInput()), 0 }

	printHeader()

	for !stack.Empty() {
		x := stack.Peek().(Symbol)

		if _, ok := p.Grammar.nonterminals[x.name]; ok {
			if i := p.Table.FindProd(x, Symbol {name: tokenString(tokenCode)}); i > 0 {
				fmt.Printf("Predict %d\t", i)
				printInput(&state, false)
				fmt.Printf("%s\n", printStack(*stack))

				stack.Pop()
				rhs := stripRhs(p.Table.Production[i])
				strs := strings.Fields(rhs)
				
				for i := len(strs) - 1; i >= 0; i-- {
					if strs[i] != lambda.name {
						stack.Push(Symbol {name: strs[i]})
					}
				}
			}


		} else if _, ok := p.Grammar.terminals[x.name]; ok {

			if x.name == tokenString(tokenCode) {
				fmt.Printf("Match!\t\t")
				printInput(&state, true)
				fmt.Printf("%s\n", printStack(*stack))

				stack.Pop()
				p.Scanner.Scan(&tokenCode, bytes.NewBuffer(*new([]byte)))
			} else {
				panic(fmt.Errorf("Expected %s, scanned %s", tokenString(tokenCode), x.name))
			}
		}

	}
}

func tokenString (t int) string {
	tokens := map[int]string{
		1:  "BeginSym",
		2:  "EndSym",
		3:  "ReadSym",
		4:  "WriteSym",
		5:  "Id",
		6:  "IntLiteral",
		7:  "LParen",
		8:  "RParen",
		9:  "SemiColon",
		10: "Comma",
		11: "AssignOp",
		12: "PlusOp",
		13: "MinusOp",
		14: "Comment",
		15: "EofSym",
	}

	return tokens[t]
}

func printHeader() {
	fmt.Printf("Parser Action\tRemaining Input\t\t\t\tParse Stack\n")
}

func printStack(s Stack) string {
	b := bytes.NewBuffer(*new ([]byte))

	for !s.Empty() {
		v := s.Pop().(Symbol).name
		b.WriteString(v + " ")
	}
	
	return b.String()
}

func (p *Parser) getInput() string {
	p.Reader.Seek(0,0)
	buf := bytes.NewBuffer(*new ([]byte))
	b, err := p.Reader.ReadByte(); 

	for err == nil {
		buf.WriteByte(b)
		b, err = p.Reader.ReadByte()
	}

	return buf.String()
}

func printInput(s *Input, d bool) {
	fmt.Printf("%s\t\t", s.str[s.i : len(s.str)])

	if d {
		s.i++
	}
}
