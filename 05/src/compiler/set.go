// Author: Michael Hunsinger
// Date:   Sept 27 2014
// File:   set.go
// Implementation of a set in golang

package compiler

import (
	"fmt"
)

type Set struct {
    set map[string]bool
}

func (s *Set) Add(str string) bool {

    _, found := s.set[str]
    s.set[str] = true
    return !found
}

func (s *Set) Print() {

	for k, _ := range s.set {
		fmt.Printf("%s\n", k)
	}
}
