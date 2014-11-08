// Author: Michael Hunsinger
// Date:   Sept 18 2014
// File:   tokens.go
// Lists the various enumerations used in compiler

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
	UnknownToken int = iota // 0
	BeginSym                // 1
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

// state enumerations and definitions
type State uint8
const (
	StartState       State = iota // 0
	EndState                      // 1
	ErrorState                    // 2
	ScanAlpha                     // 3
	ScanNumeric                   // 4
	ScanWhitespace                // 5
	ScanColon                     // 6
	ScanDash                      // 7
	ScanEquals                    // 8
	ScanComment                   // 9
	ProcessAlpha                  // 10
	ProcessNumeric                // 11
	ProcessPlusOp                 // 12
	ProcessSemicolon              // 13
	ProcessLParen                 // 14
	ProcessRParen                 // 15
	ProcessComma                  // 16
	ProcessAssign                 // 17
	ProcessMinusOp                // 18
	ProcessComment                // 19
)
