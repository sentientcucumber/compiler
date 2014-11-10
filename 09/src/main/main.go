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
	"bufio"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: ./main grammar_file program_file\n")
		os.Exit(1)
	}

	gmr, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("'%s' is not a valid file name\n", os.Args[1])
	}
	
	pgm, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		fmt.Printf("'%s' is not a valid file name\n", os.Args[1])
	}

	// the file for writing
	dst, _ := os.Create(os.Args[3])
	writer := bufio.NewWriter(dst)
	defer dst.Close()

	// First file should be the grammar, second is the file being parsed
	gmrReader := bytes.NewReader(gmr)
	pgmReader := bytes.NewReader(pgm)

	// Get the grammar from the analyzer
	a := compiler.Analyzer { Reader: *gmrReader }
	grammar := a.ReadGrammar()

	// Create a generator, necessary for table
	g := compiler.Generator { Grammar: grammar }

	// Setup the parser
	p := compiler.Parser { Grammar: a.ReadGrammar(), Reader: *pgmReader, Writer: *writer }
	p.Scanner = compiler.Scanner { Reader: *pgmReader }
	p.Table = g.GetTable()

	p.Compiler()
}
