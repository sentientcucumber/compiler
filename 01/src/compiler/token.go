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
   Empty Token = ""
   BadToken Token = "BadToken"
)

