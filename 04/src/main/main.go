// Author: Michael Hunsinger
// Date:   Aug 24 2014
// File:   main.go
// Reads the specified file and prints out a list of tokens
// Read the README.pdf for more information on compiling and running the file

package main

import (
	"bytes"
	"compiler"
	"fmt"
	"io/ioutil"
	"os"
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
	s := compiler.Scanner{Reader: *reader}
	tokenCode := 0
	tokenArray := make([]int, 100)

	for i := 0; tokenCode != compiler.EofSym && i < cap(tokenArray); i++ {
		s.Scan(&tokenCode, bytes.NewBuffer(*new([]byte)))
		tokenArray[i] = tokenCode
	}

	PrintTokens(tokenArray)
}

func PrintTokens(t []int) {

	tokens := map[int]string{
		1:  "BeginSym",
		2:  "EndSym",
		3:  "ReadSym",
		4:  "WriteSym",
		5:  "Id",
		6:  "IntLiteral",
		7:  "LParen",
		8:  "RParen",
		9:  "SemiColon",
		10: "Comma",
		11: "AssignOp",
		12: "PlusOp",
		13: "MinusOp",
		14: "Comment",
		15: "EofSym",
	}

	for _, e := range t {
		fmt.Printf("%s ", tokens[e])
	}

	fmt.Printf("\n")
}
