package main

import (
   "compiler"
   "os"
   "bufio"
   "fmt"
)

func main() {
   // open file, if an error is returned, stop all processes
   file, fileError := os.Open("sample.micro")
   if fileError != nil {
      panic(fileError)
   }

   // make sure we close the file to avoid errors later, use go's defer
   defer func() {
      if error := file.Close(); error != nil {
         panic(error)
      }
   } ()

   // create a new reader and initialze a compiler with it
   reader := bufio.NewReader(file)
   c := compiler.Compiler { Reader: *reader }

   // read through the file
   for tok := c.Scanner() ; tok != compiler.EofSym; tok = c.Scanner() {
      fmt.Printf("%v\n", tok)
   }
}

// TODO
// By inspection, what the list of tokens should be:
// L1: BeginSym 
// L2: ReadSym, LParen, Id, Comma, Id, Comma, Id, Comma, Id, RParen, SemiColon
// L3: Id, AssignOp, Id, PlusOp, LParen, Id, MinusOp, Id, RParen, MinusOp,
//     IntLiteral, SemiColon
// L4: Id, AssignOp, LParen, LParen, Id, MinusOp, LParen, IntLiteral, RParen,
//     PlusOp, LParen, PlusOp, LParen, Id, PlusOp, Id, RParen, RParen, RParen
//     MinusOp, LParen, IntLiteral, MinusOp, Id, RParen, SemiColon
// L5: WriteSym, LParen, Id, Comma, Id, PlusOp, Id, RParen, SemiColon
// L6: Nothing, comments aren't dealt with
// L7: EndSym
