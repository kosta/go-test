package doublewriter;

import (
  "fmt"
  "strings"
  "io"
  "os"
)

type DoubleWriter struct {
  io.Writer
  shouldBuffer bool
}

func NewDoubleWriter(w io.Writer, shouldBuffer bool) *DoubleWriter {
  return &DoubleWriter{w, shouldBuffer}
}

func (r DoubleWriter) Write(p []byte)(n int, err os.Error) {
  count := 0
  in := string.NewReader(p)
  for rune, size, err := in.ReadRune(); err != nil; {
  }
  
}