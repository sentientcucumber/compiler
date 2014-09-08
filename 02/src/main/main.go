// Author: Michael Hunsinger
// Date:   Aug 24 2014
// File:   main.go
// Reads the specified file and prints out a list of tokens
// Read the README.pdf for more information on compiling and running the file

package main

import (
   "compiler"
   "os"
   "bytes"
   "io/ioutil"
)

func main() {
   // create a new reader and initialze a parser with it
   buf, _ := ioutil.ReadFile(os.Args[1])
   reader := bytes.NewReader(buf)

   p := new (compiler.Parser)
   p.Scanner = compiler.Scanner { Reader: *reader }
   p.SystemGoal()
}
