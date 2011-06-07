package main

import ("fmt";)

func write(bytes []byte) {
  for i, _ := range(bytes) {
    bytes[i] = 'a'
  }
  fmt.Printf("w %p %s\n", &bytes, bytes)
}

func writesafe(bytes []byte) {
  fmt.Printf("%p %s\n", &bytes, bytes)
}

func main() {
  var s = "test"
  fmt.Printf("s %p %s\n", &s, s)
  
  var b []byte = ([]byte)(s)
  
  write(([]byte)("foo"))
  
  fmt.Printf("b %p %s\n", &b, b)
  
  fmt.Printf("s %p %s\n", &s, s)
}