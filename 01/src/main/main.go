package main

import (
   "compiler"
   "os"
   "bufio"
   "fmt"
)

func main() {
   // open file, if an error is returned, stop all processes
   file, fileError := os.Open("eof.micro")
   if fileError != nil {
      panic(fileError)
   }

   // create a new reader
   reader := bufio.NewReader(file)

   c := compiler.Compiler {Reader: *reader}
   fmt.Printf("%v\n", c.Scanner())
}

// open up file for reading, if its not there, 
// file, error := os.Open("sample.micro")
// if error != nil {
//    panic(error)
// }

// make sure we close the file to avoid errors later, use go's defer
// defer func() {
//    if error := file.Close(); error != nil {
//       panic(error)
//    }
// } ()

// make a new read buffer
// reader := bufio.NewReader(file)
// for {
//    _, _, err := reader.ReadRune()
//    if err != nil && err != io.EOF { panic(err) }
// }

// TODO

// Useful for reading the file
// file, error := ioutil.ReadFile("input.micro")
// if error != nil { panic(error) }
// error = ioutil.WriteFile("output.micro", file, 0644)
// if error != nil { panic(error) }

// create a scanner function that doesn't take any arguments and returns a token
// need all the helper functions as outlined in the assignment

// by inspection, what the list of tokens should be:
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
