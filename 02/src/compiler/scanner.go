// Author: Michael Hunsinger
// Date:   Aug 24 2014
// File:   scanner.go
// Definiton for a scanner. Uses the Scan function to return Tokens from a file

package compiler

import (
   "bytes"
   "io"
   "unicode"
   "fmt"
)

type Scanner struct {
   Reader bytes.Reader
   tokenBuffer bytes.Buffer
}

// returns the next token in the file
func (s *Scanner) Scan() Token {
   // clear the buffer
   s.ClearBuffer()

   if s.Eof() {
      return EofSym
   } else {
      // start reading the next chunk of bytes
      for !s.Eof() {

         // read the char current character
         var currChar byte
         s.Read(&currChar)

         switch {
         // if its a space, do nothing
         case unicode.IsSpace(rune(currChar)):
            break

         // if it's a letter, determine if it's an Id or reserved word
         case unicode.IsLetter(rune(currChar)):
            s.BufferChar(currChar)

            for {
               if nextChar, _ := s.Inspect();
                  unicode.IsLetter(rune(nextChar)) ||
                  unicode.IsDigit(rune(nextChar)) ||
                  nextChar == '_' {
                     s.BufferChar(nextChar)
                     s.Advance()
               } else {
                  return s.CheckReserved()
               }
            }

         // if it's a digit, it must be a intliteral
         case unicode.IsDigit(rune(currChar)):
            s.BufferChar(currChar)

            for {
               if nextChar, _ := s.Inspect(); unicode.IsDigit(rune(nextChar)) {
                  s.BufferChar(nextChar)
                  s.Advance()
               } else {
                  return IntLiteral
               }
            }

         // various 'simple' tokens
         case currChar == '(': return LParen
         case currChar == ')': return RParen
         case currChar == ';': return SemiColon
         case currChar == ',': return Comma
         case currChar == '+': return PlusOp
         case currChar == '=': return EqualityOp

         // determine if it's an assigment, if not lexical error
         case currChar == ':':
            if nextChar, _ := s.Inspect(); nextChar == '=' {
               s.Advance()
               return AssignOp
            } else {
               panic(fmt.Errorf("Lexical error when reading %c", currChar))
            }

         // determine if MinusOp or comment, if comment, consume and move on
         case currChar == '-':
            if nextChar, _ := s.Inspect(); nextChar == '-' {
               err := s.Read(&currChar)

               for currChar != '\n' && err == nil {
                  s.Read(&currChar)
               }
            } else {
               return MinusOp
            }

            // determine if exponentiation, if not lexical error
         case currChar == '*':
            if nextChar, _ := s.Inspect(); nextChar == '*' {
               s.Advance()
               return ExpOp
            } else {
               panic(fmt.Errorf("Lexical error when reading %c", currChar))
            }

         default:
            panic(fmt.Errorf("Lexical error when reading %c", currChar))
         }
      }
   }

   // if none of these, it must be the EofSym
   return EofSym
}

// Reads the next byte, will return an error if EOF
func (s *Scanner) Read(char *byte) error {
   var err error
   
   if !s.Eof() {
      *char, err = s.Reader.ReadByte()
      return nil
   } else {
      return err
   }
}

// Returns the next character but does not advance the cursor
func (s *Scanner) Inspect() (char byte, err error) {
   if char, err := s.Reader.ReadByte(); err == nil {
      s.Reader.UnreadByte()

      return char, nil
   } else {
      return 0, err
   }
}

// Advances the cursor, does not return the character
func (s *Scanner) Advance() {
   s.Reader.ReadByte()
}

// Adds the character to the tokenBuffer
func (s *Scanner) BufferChar(char byte) {
   s.tokenBuffer.WriteByte(char)
}

// Clears out the tokenBuffer
func (s *Scanner) ClearBuffer() {
   s.tokenBuffer.Reset()
}

// Determines if the tokenBuffer is a keyword or an Id
func (s *Scanner) CheckReserved() Token {
   // define a dictionary of the value in the buffer to Tokens
   dictionary := map[string]Token {
      "BEGIN": BeginSym,
      "END": EndSym,
      "READ": ReadSym,
      "WRITE": WriteSym,
   }

   buf := s.tokenBuffer.String()
   if value, exists := dictionary[buf]; exists {
      return value
   } else {
      return Id
   }
}

// Will return true if its the end of file, false if not
func (s *Scanner) Eof() bool {
   if _, err := s.Reader.ReadByte(); err == io.EOF {
      s.Reader.UnreadByte();

      return true
   } else {
      s.Reader.UnreadByte();

      return false
   }
}
