// Author: Michael Hunsinger
// Date:   Oct 18 2014
// File:   stack.go
// Implementation of a stack, credit to https://gist.github.com/bemasher/1777766 

package compiler

type Stack struct {
	top   *Node
	size  int
}

type Node struct {
	value   interface{}
	next    *Node
}

func (s *Stack) Len() int {
	return s.size
}

func (s *Stack) Push (v interface{}) {
	s.top = &Node { v, s.top }
	s.size++
}

func (s *Stack) Pop() (v interface{}) {
	if s.size > 0 {
		v, s.top = s.top.value, s.top.next
		s.size--
		return
	}

	return nil
}

func (s *Stack) Peek() (v interface{}) {
	if s.size > 0 {
		v = s.top.value
		return
	}

	return nil
}

func (s *Stack) Empty() bool {
	return s.size < 1
}
