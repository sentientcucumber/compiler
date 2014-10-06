// Author: Michael Hunsinger
// Date:   Aug 24 2014
// File:   main.go
// Reads the specified file and prints out a list of tokens
// Read the README.pdf for more information on compiling and running the file

package main

import (
	"fmt"
	"compiler"
	"io/ioutil"
	"os"
	"bytes"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("You must pass in a micro file on the command line\n")
		os.Exit(1)
	}

	src, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("'%s' is not a valid file name\n", os.Args[1])
	}

	reader := bytes.NewReader(src)

	a := compiler.Analyzer { Reader: *reader }
	a.ReadGrammar()
	compiler.Predict()
}
