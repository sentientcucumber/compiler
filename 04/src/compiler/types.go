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
type Token int
const (
	UnknownToken    int = iota // 0
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
	EofSym                     // 14
)
// TokenArray [1 7 5 8 12 5 10 5 9 14 0 0 0 0 0 0 0 0 0 0]

// state enumerations and definitions
type State uint8
const (
	StartState    State = iota // 0
	ScanAlpha                  // 1
	ScanNumeric                // 2
	ScanWhitespace             // 3
	ProcessAlpha               // 4
	ProcessNumeric             // 5
	ProcessPlusOp              // 6
	ProcessSemicolon           // 7
	ProcessLParen              // 8
	ProcessRParen              // 9
	ProcessComma               // 10
	EndState                   // 11
)
