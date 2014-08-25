// Author: Michael Hunsinger
// Date:   Aug 24 2014
// File:   main.go
// Reads the specified file and prints out a list of tokens
// Read the README.pdf for more information on compiling and running the file

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

   // read through the file by calling Scanner
   for tok := c.Scanner() ; tok != compiler.EofSym; tok = c.Scanner() {
      fmt.Printf("%v\n", tok)
   }
}
