// Author: Michael Hunsinger
// Date:   Sept 18 2014
// File:   scanner.go
// Definiton for a scanner. Uses the Scan function to return Tokens from a file

package compiler

import (
	"bytes"
	"regexp"
)

// scanner definition
type Scanner struct {
	Reader      bytes.Reader
}

// constants used for regexp
const (
	alpha       string = "[a-zA-Z]"
	numeric     string = "[0-9]"
	whitespace  string = "(?: |\n|\t)"
	plus        string = "\\+"
	dash        string = "-"
	equals      string = "="
	colon       string = ":"
	semicolon   string = ";"
	lpar        string = "\\("
	rpar        string = "\\)"
	underscore  string = "_"
	comma       string = ","
)

// Primary function of the scanner, used to scan an entire file to generate a
// list of tokens.
func (s *Scanner) Scan(tokenCode* int, tokenText* bytes.Buffer) {
	state := StartState
	tokenText.Reset()

	for state != EndState {
		currChar := s.currentChar()

		switch s.Action(state, currChar) {

		case ActionError:
			state = EndState
			*tokenCode = EofSym
			
		case MoveAppend:
			state = s.nextState(state, currChar)
			tokenText.WriteByte(currChar)
			s.consumeChar()
			
		case MoveNoAppend:
			state = s.nextState(state, currChar)
			s.consumeChar()

		case HaltAppend:
			s.lookupCode(state, currChar, tokenCode)
			tokenText.WriteByte(currChar)
			s.checkExceptions(tokenCode, *tokenText)
			s.consumeChar()
			if *tokenCode == UnknownToken {
				s.Scan(tokenCode, tokenText)
			}

			return

		case HaltNoAppend:
			s.lookupCode(state, currChar, tokenCode)
			s.checkExceptions(tokenCode, *tokenText)
			s.consumeChar()
			if *tokenCode == UnknownToken {
				s.Scan(tokenCode, tokenText)
			}

			return

		case HaltReuse:
			s.lookupCode(state, currChar, tokenCode)
			s.checkExceptions(tokenCode, *tokenText)
			if *tokenCode == UnknownToken {
				s.Scan(tokenCode, tokenText)
			}

			return
		}
	}
}

// Based on the state and current character being read, it will determine the
// next action to perform.
func (s *Scanner) Action(state State, char byte) (a Action) {

	switch state {

	case StartState:
		switch {

		case s.isAlpha(char):
			a = MoveAppend

		case s.isNumeric(char):
			a = MoveAppend

		case s.isWhitespace(char):
			a = MoveNoAppend

		case s.isPlus(char), s.isSemicolon(char), s.isLParen(char),
     		 s.isRParen(char), s.isComma(char):
			a = HaltAppend

		case s.isColon(char):
			a = MoveAppend

		case s.isEquals(char):
			a = HaltAppend

		case s.isDash(char):
			a = HaltReuse

		default:
			a = ActionError
		}

	case ScanAlpha:
		if s.isAlpha(char) || s.isNumeric(char) || s.isUnderscore(char) {
			a = MoveAppend
		} else {
			a = HaltReuse
		}

	case ScanWhitespace:
		if s.isWhitespace(char) {
			a = MoveNoAppend
		} else {
			a = MoveAppend
		}

	case ScanNumeric:
		if s.isNumeric(char) {
			a = MoveAppend
		} else {
			a = HaltReuse
		}

	case ScanColon:
		if s.isEquals(char) {
			a = HaltAppend
		} else {
			a = ActionError
		}

	case ScanDash:
		if s.isDash(char) {
			a = MoveAppend
		} else {
			a = HaltReuse
		}

	case ProcessPlusOp, ProcessSemicolon, ProcessLParen, ProcessRParen,
		ProcessComma, ProcessAssign:
		a = HaltReuse

	case ProcessComment:
		a = HaltNoAppend
		
	default:
		a = ActionError
	}

	return
}

// Determine's the next state the scanner will be in. The next state will
// determine the next action via the Action(state, char) function.
func (s *Scanner) nextState(state State, char byte) (next State) {
	
	switch state {

	case StartState:
		switch {
		
		case s.isAlpha(char):
			next = ScanAlpha

		case s.isNumeric(char):
			next = ScanNumeric

		case s.isWhitespace(char):
			next = ScanWhitespace

		case string(char) == "":
			next = EndState

		}
		
	case ScanAlpha:
		if s.isAlpha(char) || s.isNumeric(char) || s.isUnderscore(char) {
			next = ScanAlpha
		} else {
			next = ProcessAlpha
		}

	case ScanNumeric:
		if s.isNumeric(char) {
			next = ScanNumeric
		} else {
			next = ProcessNumeric
		}

	case ScanWhitespace:
		if s.isWhitespace(char) {
			next = ScanWhitespace
		} else if s.isAlpha(char) {
			next = ScanAlpha
		} else if s.isNumeric(char) {
			next = ScanNumeric
		} else if s.isPlus(char) {
			next = ProcessPlusOp
		} else if s.isSemicolon(char) {
			next = ProcessSemicolon
		} else if s.isLParen(char) {
			next = ProcessLParen
		} else if s.isRParen(char) {
			next = ProcessRParen
		} else if s.isComma(char) {
			next = ProcessComma
		} else if s.isColon(char) {
			next = ScanColon
		} else if s.isEquals(char) {
			next = ProcessAssign
		} else if s.isDash(char) {
			next = ScanDash
		}

	case ScanDash:
		if s.isDash(char) {
			next = ProcessComment
		} else {
			next = ProcessMinusOp
		}
		
	default:
		next = EndState
	}

	return 
}

// TokenCode is obtained 
func (s *Scanner) lookupCode(state State, char byte, code* int) {
	switch state {

	case ScanAlpha:
		if !s.isAlpha(char) || !s.isNumeric(char) || !s.isUnderscore(char) {
			*code = Id
		}

	case ScanNumeric:
		if !s.isNumeric(char) {
			*code = IntLiteral
		}

	case ScanColon:
		if s.isEquals(char) {
			*code = AssignOp
		} else {
			*code = UnknownToken
		}

	case ScanDash:
		if s.isDash(char) {
			*code = Comment
		} else {
			*code = MinusOp
		}

	case StartState:
		if s.isPlus(char) {
			*code = PlusOp
		} else if s.isSemicolon(char) {
			*code = SemiColon
		} else if s.isLParen(char) {
			*code = LParen
		} else if s.isRParen(char) {
			*code = RParen
		} else if s.isComma(char) {
			*code = Comma
		} else if s.isDash(char) {
			*code = MinusOp
		}

	case ProcessPlusOp:
		*code = PlusOp

	case ProcessSemicolon:
		*code = SemiColon

	case ProcessLParen:
		*code = LParen

	case ProcessRParen:
		*code = RParen

	case ProcessComma:
		*code = Comma

	case ProcessMinusOp:
		*code = MinusOp

	case ProcessComment:
		*code = Comment
		
	case EndState:
		fmt.Printf("here\n")
		*code = EofSym

	default:
		*code = 0
	}
}

// Checks to see if text is a reserved word
func (s *Scanner) checkExceptions(code* int, text bytes.Buffer) {
	switch {
	case text.String() == "BEGIN":
		*code = BeginSym

	case text.String() == "END":
		*code = EndSym

	case text.String() == "READ":
		*code = ReadSym
		
	case text.String() == "WRITE":
		*code = WriteSym
	}
}

// Consume the current character, the character is not returned.
func (s *Scanner) consumeChar() {
	s.Reader.ReadByte()
}

// Looks at the next character and returns it but does not advance the reader.
func (s *Scanner) currentChar() byte {
	if char, err := s.Reader.ReadByte(); err == nil {
		s.Reader.UnreadByte()
		return char
	} else {
		return 0
	}
}

// Determines if the character passed to it is an alpha character
func (s *Scanner) isAlpha(c byte) bool {
	re := regexp.MustCompile(alpha)
	return re.MatchString(string(c))
}

// Determines if the character passed to is numeric
func (s *Scanner) isNumeric(c byte) bool {
	re := regexp.MustCompile(numeric)
	return re.MatchString(string(c))
}

// Determines if the character passed to is whitespace
func (s *Scanner) isWhitespace(c byte) bool {
	re := regexp.MustCompile(whitespace)
	return re.MatchString(string(c))
}

// Determines if the character passed to is underscore
func (s *Scanner) isUnderscore(c byte) bool {
	re := regexp.MustCompile(underscore)
	return re.MatchString(string(c))
}

// Determines if the character passed to is plus
func (s *Scanner) isPlus(c byte) bool {
	re := regexp.MustCompile(plus)
	return re.MatchString(string(c))
}

// Determines if the character passed to is minus
func (s *Scanner) isDash(c byte) bool {
	re := regexp.MustCompile(dash)
	return re.MatchString(string(c))
}

// Determines if the character passed to is equals
func (s *Scanner) isEquals(c byte) bool {
	re := regexp.MustCompile(equals)
	return re.MatchString(string(c))
}

// Determines if the character passed to is colon
func (s *Scanner) isColon(c byte) bool {
	re := regexp.MustCompile(colon)
	return re.MatchString(string(c))
}

// Determines if the character passed to is semicolon
func (s *Scanner) isSemicolon(c byte) bool {
	re := regexp.MustCompile(semicolon)
	return re.MatchString(string(c))
}

// Determines if the character passed to is lpar
func (s *Scanner) isLParen(c byte) bool {
	re := regexp.MustCompile(lpar)
	return re.MatchString(string(c))
}

// Determines if the character passed to is rpar
func (s *Scanner) isRParen(c byte) bool {
	re := regexp.MustCompile(rpar)
	return re.MatchString(string(c))
}

// Determines if the character passed to is comma
func (s *Scanner) isComma(c byte) bool {
	re := regexp.MustCompile(comma)
	return re.MatchString(string(c))
}
