// Author: Michael Hunsinger
// Date:   Oct 18 2014
// File:   parser.go
// Parser implementation for the compiler

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!! IMPORTANT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// There are more Parser functions in the semantic.go file, they are separated
// to keep things organized

package compiler

import (
	"fmt"
	"bytes"
	"strings"
	"bufio"
	"strconv"
	"regexp"
)

type Parser struct {
	Grammar   Grammar
	Scanner   Scanner
	Table     Table
	Reader    bytes.Reader
	Writer    bufio.Writer

	maxSymbol   int
	maxTemp     int
	lastSymbol  int
	SymbolTable []string
}

type Input struct {
	str       []string
	i         int
}

var (
	ss           []Symbol
	currentIndex int
	topIndex     int
	leftIndex    int
	rightIndex   int 
)

func (p *Parser) Compiler() {
	// Initialize parse stack
	start := findStartSymbol(p.Grammar)
	stack := new (Stack)
	stack.Push(start)

	// Initialize semantic stack
	rightIndex = 0; leftIndex = 0; currentIndex = 1; topIndex = 2
	ss = []Symbol {}
	ss = append(ss, start)

	tokenCode := 0
	p.Scanner.Scan(&tokenCode, bytes.NewBuffer(*new([]byte)))
	state := Input { strings.Fields(p.getInput()), 0 }

	for !stack.Empty() {
		x := stack.Peek().(Symbol)

		if _, ok := p.Grammar.nonterminals[x.name]; ok {
			if i := p.Table.FindProd(x, Symbol {name: tokenString(tokenCode)}); i > 0 {
				printInput(&state, false)
				fmt.Printf("PS:\t\t%s\n", printStack(*stack))
				fmt.Printf("SS:\t\t%s\n", printArray(ss))

				stack.Pop()
				stack.Push(EOPSymbol(leftIndex, rightIndex, currentIndex, topIndex))

				rhs := stripRhs(p.Table.Production[i])
				strs := strings.Fields(rhs)
				
				// Add symbols in reverse order for the parse stack
				for i := len(strs) - 1; i >= 0; i-- {
					if strs[i] != lambda.name {
						stack.Push(Symbol { name: strs[i] })
					}
				}

				// Add symbols in order for the semantic stack
				count := 0
				for i := 0; i < len(strs); i++ {
					if strs[i] != lambda.name && strs[i][0] != '#' {
						ss = insert(ss, 0, Symbol { name: strs[i] })
						// ss = append(ss, Symbol { name: strs[i] })
						count++
					}
				}
				
				// Print then update indices
				fmt.Printf("Indices:(%d, %d, %d, %d)\n", leftIndex, rightIndex, currentIndex, topIndex)
				leftIndex = currentIndex; rightIndex = topIndex; currentIndex = rightIndex; topIndex += count
			} else {
				panic(fmt.Errorf("Could not find a production for <%s, %s>",
					x.name, tokenString(tokenCode)))
			}
		} else if _, ok := p.Grammar.terminals[x.name]; ok {
			if x.name == tokenString(tokenCode) {
				printInput(&state, true)
				fmt.Printf("PS:\t\t%s\n", printStack(*stack))
				fmt.Printf("SS:\t\t%s\n", printArray(ss))
				fmt.Printf("Indices:(%d, %d, %d, %d)\n", leftIndex, rightIndex, currentIndex, topIndex)

				stack.Pop()
				p.Scanner.Scan(&tokenCode, bytes.NewBuffer(*new([]byte)))
				currentIndex++
			} else {
				panic(fmt.Errorf("Expected %s, scanned %s", tokenString(tokenCode), x.name))
			}
		} else if x.category == EOP {
			leftIndex, rightIndex, currentIndex, topIndex = unpackEOP(x)
			stack.Pop()
		} else {
			stack.Pop()
			fmt.Printf("%s\n", x.name)
		}
		fmt.Printf("\n")
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
	fmt.Printf("Input:\t%s\n", s.str[s.i : len(s.str)])

	if d {
		s.i++
	}
}

func EOPSymbol (l, r, c, t int) Symbol {
	lstr := strconv.Itoa(l)
	rstr := strconv.Itoa(r)
	cstr := strconv.Itoa(c)
	tstr := strconv.Itoa(t)

	return Symbol { "EOP(" + lstr + ", " + rstr + ", " + cstr + ", " + tstr + ")", EOP }
}

func unpackEOP (s Symbol) (l, r, c, t int) {
	temp := s.name
	re := regexp.MustCompile("[A-Z\\W]*")
	temp = re.ReplaceAllString(temp, " ")
	strs := strings.Fields(temp)

	l, _ = strconv.Atoi(strs[0])
	r, _ = strconv.Atoi(strs[1])
	c, _ = strconv.Atoi(strs[2])
	t, _ = strconv.Atoi(strs[3])

	return
}

func printArray (s []Symbol) string {
	b := bytes.NewBuffer(*new ([]byte))

	for _, v := range s {
		b.WriteString(v.name + " ")
	}

	return b.String()
}

func insert(slice []Symbol, index int, value Symbol) []Symbol {
	slice = append(slice[0:], append(make([]Symbol, 1), slice[:0]...)...)
	copy(slice[index + 1:], slice[index:])
	slice[index] = value

	return slice
}
