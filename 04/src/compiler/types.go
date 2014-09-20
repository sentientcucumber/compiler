// Author: Michael Hunsinger
// Date:   Sept 18 2014
// File:   types.go
// Lists the various types used in compilers

package compiler

// action enumeration and definitions
type Action uint8
const (
	ActionError  Action = iota // 0
	MoveAppend                 // 1
	MoveNoAppend               // 2
	HaltAppend                 // 3
	HaltNoAppend               // 4
	HaltReuse                  // 5
)

// token enumeration and definitions
type Token uint8
const (
	TokenError Token = iota    // 0
	BeginSym                   // 1
	EndSym                     // 2
	ReadSym                    // 3
	WriteSym                   // 4
	Id                         // 5
	IntLiteral                 // 6
	LParen                     // 7
	RParen                     // 8
	SemiColon                  // 9
	Comma                      // 10
	AssignOp                   // 11
	PlusOp                     // 12
	MinusOp                    // 13
	ExpOp                      // 14
	EqualityOp                 // 15
	EofSym                     // 16
)

// state enumerations and definitions
type State uint8
const (
	StartState State = iota    // 0
	EndState                   // 1
	ReadState                  // 2
	ProcessState               // 3
)
