package main

import (
	"net"
	"fmt"
  "os/signal"
)

func main() {
  newConnections := make(chan *net.TCPConn)
  
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
			  return 
	    case conn := <- newConnections:
  	    fmt.Fprint(conn, "Hello. Goodbye.\n")
        conn.Close()
	  }
	}
}
