// Author: Michael Hunsinger
// Date:   Aug 24 2014
// File:   compiler.go
// The compiler class, primarily uses the Scanner function to read a file
// to return a token

package compiler

import (
   "bufio"
   "bytes"
   "io"
   "unicode"
   "fmt"
)

type Compiler struct {
   Reader bufio.Reader
   tokenBuffer bytes.Buffer
}

// 
func (c *Compiler) Scanner() Token {
   // clear the buffer
   c.ClearBuffer()

   if c.Eof() {
      return EofSym
   } else {
      // start reading the next chunk of bytes
      for !c.Eof() {

         // read the char current character
         var currChar byte
         c.Read(&currChar)

         switch {
         case unicode.IsSpace(rune(currChar)):
            break

         case unicode.IsLetter(rune(currChar)):
            c.BufferChar(currChar)

            for {
               if nextChar, _ := c.Inspect();
                  unicode.IsLetter(rune(nextChar)) ||
                  unicode.IsDigit(rune(nextChar)) ||
                  nextChar == '_' {
                     c.BufferChar(nextChar)
                     c.Advance()
               } else {
                  return c.CheckReserved()
               }
            }

         case unicode.IsDigit(rune(currChar)):
            c.BufferChar(currChar)

            for {
               if nextChar, _ := c.Inspect(); unicode.IsDigit(rune(nextChar)) {
                  c.BufferChar(nextChar)
                  c.Advance()
               } else {
                  return IntLiteral
               }
            }

         case currChar == '(': return LParen
         case currChar == ')': return RParen
         case currChar == ';': return SemiColon
         case currChar == ',': return Comma
         case currChar == '+': return PlusOp
         case currChar == '=': return EqualityOp

         case currChar == ':':
            if nextChar, _ := c.Inspect(); nextChar == '=' {
               c.Advance()
               return AssignOp
            } else {
               panic(fmt.Errorf("Lexical error when reading %c", currChar))
            }
            
         case currChar == '-':
            if nextChar, _ := c.Inspect(); nextChar == '-' {
               err := c.Read(&currChar)

               for currChar != '\n' && err == nil {
                  c.Read(&currChar)
               }
            } else {
               return MinusOp
            }

         case currChar == '*':
            if nextChar, _ := c.Inspect(); nextChar == '*' {
               c.Advance()
               return ExpOp
            } else {
               panic(fmt.Errorf("Lexical error when reading %c", currChar))
            }

         default:
            panic(fmt.Errorf("Lexical error when reading %c", currChar))
         }
      }
   }
   return EofSym
}

// Reads the next byte, will return an error if EOF
func (c *Compiler) Read(char *byte) error {
   var err error
   
   if !c.Eof() {
      *char, err = c.Reader.ReadByte()
      return nil
   } else {
      return err
   }
}

// Returns the next character but does not advance the cursor
func (c *Compiler) Inspect() (char byte, err error) {
   if char, err := c.Reader.ReadByte(); err == nil {
      c.Reader.UnreadByte()

      return char, nil
   } else {
      return 0, err
   }
}

// Advances the cursor, does not return the character
func (c *Compiler) Advance() {
   c.Reader.ReadByte()
}

// Adds the character to the tokenBuffer
func (c *Compiler) BufferChar(char byte) {
   c.tokenBuffer.WriteByte(char)
}

// Clears out the tokenBuffer
func (c *Compiler) ClearBuffer() {
   c.tokenBuffer.Reset()
}

// Determines if the tokenBuffer is a keyword or an Id
func (c *Compiler) CheckReserved() Token {
   // define a dictionary of the value in the buffer to Tokens
   dictionary := map[string]Token {
      "BEGIN": BeginSym,
      "END": EndSym,
      "READ": ReadSym,
      "WRITE": WriteSym,
   }

   buf := c.tokenBuffer.String()
   if value, exists := dictionary[buf]; exists {
      return value
   } else {
      return Id
   }
}

// Will return true if its the end of file, false if not
func (c *Compiler) Eof() bool {
   if _, err := c.Reader.ReadByte(); err == io.EOF {
      c.Reader.UnreadByte();

      return true
   } else {
      c.Reader.UnreadByte();

      return false
   }
}
