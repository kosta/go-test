package main

import (
	"net"
	"fmt"
  "os/signal"
)

func makeSenderChan(conn net.Conn) chan<- string {
  ch := make(chan string)
  go func() {
    //close channel when returning
    defer close(ch)
    //close conn when returning
    defer conn.Close()
    
    for {
      s, ok := <-ch
      if !ok { 
        return
      }
      _, err := fmt.Fprint(conn, s)
      if err != nil {
        return
      }
    }
  }()
  return ch
}

func broadcast(connections map[chan<- string]bool, msg string) {
  for conn, _ := range(connections) {
    conn <- msg
  }
}

func main() {
  newConnections := make(chan *net.TCPConn)
  connections := make(map[chan<- string]bool)
  
	tcp, tcperr := net.ListenTCP("tcp", &net.TCPAddr{})
	if tcperr != nil {
		fmt.Println("Error opening TCP socket:", tcperr)
		return
	}
	defer tcp.Close()
	fmt.Println("Listening on", tcp.Addr().Network(), tcp.Addr())
	
	go func() {
	  for {
  		conn, connerr := tcp.AcceptTCP()
  		if connerr != nil {
  			fmt.Println("Error accepting TCP connection:", connerr)
  		} else {
        newConnections <- conn
  		}
	  }
	}() //go func accepting connections
	
	for {
	  select {
	    case signal := <-signal.Incoming:
			  fmt.Println("got signal:", signal)
			  broadcast(connections, "Goodbye.\n")
			  return 
	    case conn := <- newConnections:
  	    fmt.Fprint(conn, "Hello.\n")	    
	      connections[makeSenderChan(conn)] = true
	  }
	}
}
