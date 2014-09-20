// Author: Michael Hunsinger
// Date:   Sept 18 2014
// File:   scanner.go
// Definiton for a scanner. Uses the Scan function to return Tokens from a file

package compiler

import (
	"bytes"
	"fmt"
	"regexp"
)

// scanner definition
type Scanner struct {
	Reader      bytes.Reader
}

// constants used for regexp
const (
	letter      string = "[a-zA-Z]"
	digit       string = "[0-9]"
	whitespace  string = "(\t)*?( )*?(\n)*?"
	plus        string = "\\+"
	minus       string = "-"
	equals      string = "="
	colon       string = ":"
	semicolon   string = ";"
	lpar        string = "("
	rpar        string = ")"
	underscore  string = "_"
)

// Primary function of the scanner, used to scan an entire file to generate a
// list of tokens.
func (s *Scanner) Scan(tokenCode Token, tokenText bytes.Buffer) Token {
	state := StartState
	tokenText.Reset()

	for (state != EndState) {
		switch s.Action(&state, s.CurrentChar()) {

		case ActionError:

		case MoveAppend:
			state = s.NextState(&state, s.CurrentChar())
			tokenText.WriteByte(s.CurrentChar())
			s.ConsumeChar()
			
		case MoveNoAppend:
			state = s.NextState(&state, s.CurrentChar())
			s.ConsumeChar()

		case HaltAppend:
			s.LookupCode(state, s.CurrentChar(), &tokenCode)
			tokenText.WriteByte(s.CurrentChar())
			s.CheckExceptions(&tokenCode, tokenText)
			s.ConsumeChar()
			if tokenCode == TokenError {
				s.Scan(tokenCode, tokenText)
			}

			return tokenCode

		case HaltNoAppend:
			s.LookupCode(state, s.CurrentChar(), &tokenCode)
			s.CheckExceptions(&tokenCode, tokenText)
			s.ConsumeChar()
			if tokenCode == TokenError {
				s.Scan(tokenCode, tokenText)
			}

			return tokenCode

		case HaltReuse:
			s.LookupCode(state, s.CurrentChar(), &tokenCode)
			s.CheckExceptions(&tokenCode, tokenText)
			if tokenCode == TokenError {
				s.Scan(tokenCode, tokenText)
			}

			return tokenCode
		}

		fmt.Printf("end state: %v, action: %v, char %c\n", state, action, s.CurrentChar())
	}

	return EofSym
}

// Based on the state (why is the character in here?), it will determine the
// next action to perform.
func (s *Scanner) Action(state* State, char byte) Action {
	
	return MoveAppend
}

// Determine's the next state the scanner will be in. The next state will
// determine the next action via the Action(state, char) function.
func (s *Scanner) NextState(state* State, char byte) State {

	return ProcessState
}

// Consume the current character, the character is not returned.
func (s *Scanner) ConsumeChar() {
	s.Reader.ReadByte()
}


func (s *Scanner) LookupCode(state State, char byte, code* Token) {
	
}

func (s *Scanner) CheckExceptions(code* Token, text bytes.Buffer) {

}

// Looks at the next character and returns it but does not advance the reader.
func (s *Scanner) CurrentChar() byte {
	if char, err := s.Reader.ReadByte(); err == nil {
		s.Reader.UnreadByte()
		
		return char
	} else {
		return 0
	}
}
