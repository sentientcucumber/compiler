// Author: Michael Hunsinger
// Date:   Aug 24 2014
// File:   compiler.go
// The Token class, enumerations for allowable Tokens

package compiler

type Token string

const (
   BeginSym Token = "BeginSym"
   EndSym Token = "EndSym"
   ReadSym Token = "ReadSym"
   WriteSym Token = "WriteSym"
   Id Token = "Id"
   IntLiteral Token = "IntLiteral"
   LParen Token = "LParen"
   RParen Token = "RParen"
   SemiColon Token = "SemiColon"
   Comma Token = "Comma"
   AssignOp Token = "AssignOp"
   PlusOp Token = "PlusOp"
   MinusOp Token = "MinusOp"
   EofSym Token = "EofSym"
)

