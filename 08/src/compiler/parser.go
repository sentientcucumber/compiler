// Author: Michael Hunsinger
// Date:   Oct 18 2014
// File:   parser.go
// Parser implementation for the compiler

package compiler

import (
	"fmt"
)

type Parser struct {
	Grammar   Grammar
}

func (p *Parser) Driver() {
	var _ = fmt.Printf

	start := findStartSymbol(p.Grammar)
	stack := new (Stack)
	stack.Push(start)

	for !stack.Empty() {
		top := stack.Peek().(string)

		if _, ok := p.Grammar.nonterminals[top]; ok {
			
		}
	}
}
