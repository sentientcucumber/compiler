// Author: Michael Hunsinger
// Date:   Sept 13 2014
// File:   types.go
// Lists the various types used in compilers

package compiler

// Constant values
const (
   // tokens
	BeginSym    Token = "BeginSym"
	EndSym      Token = "EndSym"
	ReadSym     Token = "ReadSym"
	WriteSym    Token = "WriteSym"
	Id          Token = "Id"
	IntLiteral  Token = "IntLiteral"
	LParen      Token = "LParen"
	RParen      Token = "RParen"
	SemiColon   Token = "SemiColon"
	Comma       Token = "Comma"
	AssignOp    Token = "AssignOp"
	PlusOp      Token = "PlusOp"
	MinusOp     Token = "MinusOp"
	ExpOp       Token = "ExpOp"
	EqualityOp  Token = "EqualityOp"
	EofSym      Token = "EofSym"

   // operators
   OpPlusOp    Operator = "PlusOp"
   OpMinusOp   Operator = "MinusOp"

   // expression kind
   IdExpr      ExprKind = "IdExpr"
   LiteralExpr ExprKind = "LiteralExpr"
   TempExpr    ExprKind = "TempExpr"
)

// Token type
type Token string

// Operator type
type Operator string

// ExprKind type
type ExprKind string

// OpRec type
type OpRec struct {
   Op Operator
}

// ExprRec type
type ExprRec struct {
   Kind ExprKind
   Name string
   Val  int
}
