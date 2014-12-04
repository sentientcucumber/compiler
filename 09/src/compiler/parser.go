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
	ss           []SemanticRecord
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
	blank := SemanticRecord {}
	ss = []SemanticRecord { blank }

	srStart := new(SemanticRecord)
	srStart.exprRec.Name = start.name

	ss = append(ss, *srStart)

	tokenCode := 0
	tokenText := bytes.NewBuffer(*new([]byte))
	p.Scanner.Scan(&tokenCode, tokenText)
	state := Input { strings.Fields(p.getInput()), 0 }

	fmt.Printf("Initial state before entering loop\n")
	printInput(&state, false)
	fmt.Printf("PS:\t\t%s\n", printStack(*stack))
	fmt.Printf("SS:\t\t%s\n", printArray(ss))
	fmt.Printf("Indices:(%d, %d, %d, %d)\n\n", leftIndex, rightIndex, currentIndex, topIndex)

	for !stack.Empty() {
		x := stack.Peek().(Symbol)

		if _, ok := p.Grammar.nonterminals[x.name]; ok {
			if i := p.Table.FindProd(x, Symbol {name: tokenString(tokenCode)}); i > 0 {
				stack.Pop()
				rhs := stripRhs(p.Table.Production[i])

				if rhs != lambda.name {
					stack.Push(EOPSymbol(leftIndex, rightIndex, currentIndex, topIndex))
					strs := strings.Fields(rhs)
					
					// Add symbols in reverse order for the parse stack
					for i := len(strs) - 1; i >= 0; i-- {
						// if strs[i] != lambda.name {
							stack.Push(Symbol { name: strs[i] })
						// }
					}

					// Add symbols in order for the semantic stack
					count := 0
					for i := 0; i < len(strs); i++ {
						if strs[i] != lambda.name && strs[i][0] != '#' {
							sr := new(SemanticRecord)
							sr.exprRec.Name = strs[i]

							ss = append(ss, *sr)
							count++
						}
					}

					leftIndex = currentIndex
					rightIndex = topIndex
					currentIndex = rightIndex
					topIndex = topIndex + count

					fmt.Printf("NONTERMINAL T(%s, %s) = %d\n", x.name, tokenString(tokenCode), i)
					printInput(&state, false)
					fmt.Printf("PS:\t\t%s\n", printStack(*stack))
					fmt.Printf("SS:\t\t%s\n", printArray(ss))
					fmt.Printf("Indices:(%d, %d, %d, %d)\n\n", leftIndex, rightIndex, currentIndex, topIndex)
				}
			} else {
				panic(fmt.Errorf("Could not find a production for <%s, %s>",
					x.name, tokenString(tokenCode)))
			}
		} else if _, ok := p.Grammar.terminals[x.name]; ok {
			if x.name == tokenString(tokenCode) {
				ss[currentIndex].exprRec.Name = tokenText.String()
				stack.Pop()
				p.Scanner.Scan(&tokenCode, tokenText)
				currentIndex++

				fmt.Printf("TERMINAL X = %s\n", x.name)
				if (x.name == "EofSym") {
					printInput(&state, false)
				} else {
					printInput(&state, true)
				}
				fmt.Printf("PS:\t\t%s\n", printStack(*stack))
				fmt.Printf("SS:\t\t%s\n", printArray(ss))
				fmt.Printf("Indices:(%d, %d, %d, %d)\n\n", leftIndex, rightIndex, currentIndex, topIndex)
			} else {
				panic(fmt.Errorf("Expected %s, scanned %s", tokenString(tokenCode), x.name))
			}
		} else if x.category == EOP {
			leftIndex, rightIndex, currentIndex, topIndex = unpackEOP(x)
			currentIndex++
			ss = ss[:len(ss) - 1]
			stack.Pop()

			fmt.Printf("X = EOP %s\n", x.name)
			printInput(&state, false)
			fmt.Printf("PS:\t\t%s\n", printStack(*stack))
			fmt.Printf("SS:\t\t%s\n", printArray(ss))
			fmt.Printf("Indices:(%d, %d, %d, %d)\n\n", leftIndex, rightIndex, currentIndex, topIndex)
		} else {
			stack.Pop()
			p.processSemanticRoutine(x.name, leftIndex, rightIndex)

			fmt.Printf("X = %s\n", x.name)
			printInput(&state, false)
			fmt.Printf("PS:\t\t%s\n", printStack(*stack))
			fmt.Printf("SS:\t\t%s\n", printArray(ss))
			fmt.Printf("Indices:(%d, %d, %d, %d)\n\n", leftIndex, rightIndex, currentIndex, topIndex)			
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

func printStack(s Stack) string {
	b := bytes.NewBuffer(*new ([]byte))

	for !s.Empty() {
		v := s.Pop().(Symbol)
		if v.category == EOP {
			b.WriteString("EOP(" + v.name + ") ")
		} else {
			b.WriteString(v.name + " ")
		}
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
	if d {
		s.i++
	}

	fmt.Printf("Input:\t%s\n", s.str[s.i : len(s.str)])
}

func EOPSymbol (l, r, c, t int) Symbol {
	lstr := strconv.Itoa(l)
	rstr := strconv.Itoa(r)
	cstr := strconv.Itoa(c)
	tstr := strconv.Itoa(t)

	return Symbol { lstr + " " + rstr + " " + cstr + " " + tstr, EOP }
}

func unpackEOP (s Symbol) (l, r, c, t int) {
	temp := s.name
	strs := strings.Fields(temp)

	l, _ = strconv.Atoi(strs[0])
	r, _ = strconv.Atoi(strs[1])
	c, _ = strconv.Atoi(strs[2])
	t, _ = strconv.Atoi(strs[3])

	return
}

func printArray (s []SemanticRecord) string {
	b := bytes.NewBuffer(*new ([]byte))

	// for _, v := range s {
	for i := len(s) - 1; i >= 0; i-- {
		b.WriteString("\"" + s[i].exprRec.Name + "\" ")
	}

	return b.String()
}

func insert(slice []SemanticRecord, index int, value SemanticRecord) []SemanticRecord {
	slice = append(slice[0:], append(make([]SemanticRecord, 1), slice[:0]...)...)
	copy(slice[index + 1:], slice[index:])
	slice[index] = value

	return slice
}

func (p *Parser) processSemanticRoutine(s string, l, r int) {
	switch s {
	case "#Start":
		p.Start()
	case "#Assign($1,$3)":
		p.Assign(ss[r], ss[r + 2])
	case "#ReadId($1)":
		p.ReadId(ss[r])
	case "#WriteExpr($1)":
		p.WriteExpr(ss[r])
	case "#Copy($1,$2)":
		p.semanticCopy(r, r + 1)
	case "#Copy($2,$$)":
		p.semanticCopy(r + 1, l)
	case "#Copy($1,$$)":
		p.semanticCopy(r, l)
	case "#GenInfix($$,$1,$2,$$)":
		p.GenInfix(ss[l], ss[r], ss[r + 1], &ss[l])
	case "#ProcessLiteral($$)":
		p.ProcessLiteral(ss[l])
	case "#ProcessOp($$)":
		p.ProcessOp(ss[l])
	case "#ProcessId($$)":
		p.ProcessId(ss[l])
	case "#Finish":
		p.finish()
	default:
		panic(fmt.Errorf("Unknown semantic routine %s", s))
	}
}
