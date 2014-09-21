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
	Whitespace                 // 15
	Comment                    // 16
)
// 1 5 11 5 12 6 12 5 9 9 5 5 2 14

// state enumerations and definitions
type State uint8
const (
	StartState    State = iota // 0
	EndState                   // 1
	ScanAlpha                  // 2
	ScanNumeric                // 3
	ScanWhitespace             // 4
	ScanColon                  // 5
	ScanDash                   // 6
	ScanEquals                 // 7
	ScanComment                // 8
	ProcessAlpha               // 9
	ProcessNumeric             // 10
	ProcessPlusOp              // 11
	ProcessSemicolon           // 12
	ProcessLParen              // 13
	ProcessRParen              // 14
	ProcessComma               // 15
	ProcessAssign              // 16
	ProcessMinusOp             // 17
	ProcessComment             // 18
)
