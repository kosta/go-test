package main

import (
	"net"
	"fmt"
)

func main() {
	tcp, tcperr := net.ListenTCP("tcp", &net.TCPAddr{})
	if tcperr != nil {
		fmt.Println("Error opening TCP socket:", tcperr)
		return
	}
	defer tcp.Close()
	fmt.Println("Listening on", tcp.Addr().Network(), tcp.Addr())
	for {
		conn, connerr := tcp.AcceptTCP()
		if connerr != nil {
			fmt.Println("Error accepting TCP connection:", connerr)
		} else {
		  //same as conn.Write([]byte("Hello. Goodbye.\n"))
			fmt.Fprint(conn, "Hello. Goodbye.\n")
			conn.Close()
		}
	}
}
