package compiler

import (
   "bufio"
   "os"
   "bytes"
)

type Compiler struct {
   Reader bufio.Reader
   buffer bytes.Buffer
}

// Sets up a file for reading, will stop program if the file doesn't exist
func (c *Compiler) Compiler(filename string) {
   // open file, if an error is returned, stop all processes
   file, fileError := os.Open(filename)
   if fileError != nil {
      panic(fileError)
   }
   
   // finally, create the Reader to be used later
   reader := bufio.NewReader(file)
   reader.Peek(0)
}

func (c *Compiler) Scanner() Token {
   // clear the buffer
   c.ClearBuffer()

   // read the next character
   char, _, readError := c.Reader.ReadRune()

   if readError != nil {
      return EofSym
      os.Exit(0)
   } else {
      for readError == nil {
         c.Read(char)
      }
   }

   return SemiColon
}

func (c *Compiler) Read(char rune) {
   
}

// func (s *Scanner) Inspect() rune {

// }

// func (s *Scanner) Advance() {

// }

// func (s *Scanner) Eof() bool {

// }

// func (s *Scanner) BufferChar(c rune) {

// }

func (c *Compiler) ClearBuffer() {
   c.buffer.Reset()
}

// func (s *Scanner) CheckReserved() {

// }
