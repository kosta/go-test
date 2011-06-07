package main

import (
	"net"
	"fmt"
	"syscall"
	"os/signal"
	"bufio"
)

type ConnMessage struct {
	Conn    net.Conn
	Message string
}

func main() {
	connections := make(map[net.Conn]bool)

	newMessages := make(chan ConnMessage)
	go func() {
		count := 0
		for {
			newMessages <- ConnMessage{nil, fmt.Sprintf("ping%d\r\n", count)}
			count++
			err := syscall.Sleep(1e10) //10 secs
			if err != 0 {
				return
			}
		}
	}()

	tcp, tcperr := net.ListenTCP("tcp", &net.TCPAddr{})
	if tcperr != nil {
		fmt.Println("Error opening TCP socket:", tcperr)
		return
	}
	defer tcp.Close()
	newConnections := make(chan *net.TCPConn)
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
	}()

	for {
		select {
		case signal := <-signal.Incoming:
			fmt.Println("got signal:", signal)
			return
		case conn := <-newConnections:
			fmt.Printf("Got connection from %s %s\n",
				conn.RemoteAddr().Network(),
				conn.RemoteAddr())
			conn.SetNoDelay(true)
			conn.Write([]uint8("Hello\r\nHow are you?\r\n"))
			connections[conn] = true
			go func() {
				defer conn.Close()
				reader := bufio.NewReader(conn)
				for {
					line, err := reader.ReadString('\n')
					if err != nil {
						//todo: remove conn from connections
						fmt.Printf("error reading %s %s: %s\n",
							conn.RemoteAddr().Network(),
							conn.RemoteAddr(),
							err)
						connections[conn] = false, false
						return //closes conn
					}
					newMessages <- ConnMessage{conn, string(line)}
				}
			}()
		case msg := <-newMessages:
			for conn, _ := range connections {
				if conn != msg.Conn {
					_, err := fmt.Fprintln(conn, msg.Message)
					if err != nil {
						fmt.Printf("error writing%s: %s\n",
							conn,
							conn.RemoteAddr(),
							err)
						connections[conn] = false, false
					}
				}
			}
		}
	}
}
