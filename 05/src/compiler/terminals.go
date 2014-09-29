// Author: Michael Hunsinger
// Date:   Sept 27 2014
// File:   terminal.go
// The various terminals

package compiler

type Terminal int
const (
	BeginSym    int = iota  // 1
	EndSym                  // 2
	ReadSym                 // 3
	WriteSym                // 4
	Id                      // 5
	IntLiteral              // 6
	LParen                  // 7
	RParen                  // 8
	SemiColon               // 9
	Comma                   // 10
	AssignOp                // 11
	PlusOp                  // 12
	MinusOp                 // 13
	Comment                 // 14
	EofSym                  // 15
)

