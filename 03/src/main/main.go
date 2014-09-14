// Author: Michael Hunsinger
// Date:   Aug 24 2014
// File:   main.go
// Reads the specified file and prints out a list of tokens
// Read the README.pdf for more information on compiling and running the file

package main

import (
	"bytes"
	"compiler"
	"io/ioutil"
	"os"
   "bufio"
)

func main() {
   // the file for reading
	src, _ := ioutil.ReadFile(os.Args[1])
	reader := bytes.NewReader(src)

   // the file for writing
   dst, _ := os.Create(os.Args[2])
   writer := bufio.NewWriter(dst)
   defer dst.Close()

   // setup the parser
   p := compiler.Parser { Writer: *writer }
	p.Scanner = compiler.Scanner { Reader: *reader }

   // parse the file!
	p.SystemGoal()
}
