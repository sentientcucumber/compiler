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
	Comment                    // 14
	EofSym                     // 15
	Whitespace
)

// state enumerations and definitions
type State uint8
const (
	StartState    State = iota // 0
	EndState                   // 1
	ScanAlpha                  // 2
	ScanNumeric                // 3
	ScanWhitespace             // 4
	ProcessAlpha               // 5
	ProcessNumeric             // 6
	ProcessPlusOp              // 7
	ProcessSemicolon           // 8
	ProcessLParen              // 9
	ProcessRParen              // 10
	ProcessComma               // 11
	ProcessWhitespace
	ProcessAssign
	ProcessMinusOp
	ProcessComment
	ScanColon
	ScanDash
	ScanEquals
)
