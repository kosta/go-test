package main

import (
	"net"
	"fmt"
  "os/signal"
)

func broadcast(connections map[net.Conn]bool, msg string) {
  for conn, _ := range(connections) {
    go func() {
      fmt.Fprint(conn, msg)
    }()
  }
}

func main() {
  newConnections := make(chan *net.TCPConn)
  connections := make(map[net.Conn]bool)
  
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
	      connections[conn] = true
  	    fmt.Fprint(conn, "Hello.\n")
	  }
	}
}
