package main

import (
	"net"
	"fmt"
  "os/signal"
)

type channelAndConn struct {
  ch chan<- string
  conn net.Conn
}

func makeSenderChan(conn net.Conn, broken chan<- channelAndConn) chan<- string {
  ch := make(chan string)
  go func() {
    //close channel when returning
    defer close(ch)
    //close conn when returning
    defer conn.Close()
    defer func() { broken <- channelAndConn{ch,conn} }()
    
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
  brokenConnections := make(chan channelAndConn)
  
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
	      fmt.Printf("Connected: %s %s\n",
	   			conn.RemoteAddr().Network(),
	 	   		conn.RemoteAddr())
  	    fmt.Fprint(conn, "Hello.\n")	    
	      connections[makeSenderChan(conn, brokenConnections)] = true
	    case cc := <- brokenConnections:
	      fmt.Printf("Disconnected: %s %s\n", 
	       cc.conn.RemoteAddr().Network(),
	       cc.conn.RemoteAddr())
	 	   	//remove from map
	 	    connections[cc.ch] = false, false
	 	}
	}
}
