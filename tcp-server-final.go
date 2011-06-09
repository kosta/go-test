package main

import (
	"fmt"
	"net"
	"bufio"
	"sync"
	"os/signal"
)

//slide 58
func readLinesAndSendToChan(conn net.Conn, messages chan string) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Read error: %s, closing %s\n", err, conn.RemoteAddr())
			break
		} //closes conn due to defer
		messages <- line
	}
}

//slide 60
func writeLinesFromChan(conn net.Conn,
brokenChannels chan chan<- string,
waits *sync.WaitGroup) chan<- string {
	waits.Add(1)
	ch := make(chan string, 100) //buffer of 100 strings
	go func() {
		defer func() { brokenChannels <- ch }()
		defer waits.Done()
		addr := conn.RemoteAddr()
		for msg := range ch {
			_, err := fmt.Fprint(conn, msg)
			if err != nil {
				fmt.Printf("Write error: closing %s\n", addr)
				break
			}
		}
	}()
	return ch
}

func main() {
	//open TCP port - slide 52
	tcp, tcperr := net.Listen("tcp", ":0")
	if tcperr != nil {
		fmt.Println("Error opening TCP socket:", tcperr)
		return
	}
	fmt.Println("Listening on", tcp.Addr().Network(), tcp.Addr())

	//tcp.Accept gofunc - slide 55
	newConnections := make(chan net.Conn)
	go func() {
		for {
			conn, connerr := tcp.Accept()
			if connerr != nil {
				return //skip error handling for this example
			} else {
				newConnections <- conn
			}
		}
	}()

	//slide 57 - overwritten by slide 59
	messages := make(chan string)
	/*for {
			select {
			case conn := <-newConnections:
				fmt.Fprint(conn, "Hello.\n")
	      go readLinesAndSendToChan(conn, messages)
			case msg := <-messages:
				fmt.Println("Got message:", msg)
			}
		}*/

	waits := new(sync.WaitGroup)

	//slide 59
	connections := make(map[chan<- string]bool) //used as set
	brokenChannels := make(chan chan<- string)  //channel of channels
	for {
		select {
		case conn := <-newConnections:
			go readLinesAndSendToChan(conn, messages)
			connections[writeLinesFromChan(conn, brokenChannels, waits)] = true
		case msg := <-messages:
			fmt.Print("got message: ", msg)
			for conn, _ := range connections {
				conn <- msg
			}
		case broken := <-brokenChannels:
			close(broken)
			connections[broken] = false, false //remove from map        
		case signal := <-signal.Incoming:
			fmt.Printf("got signal: '%s'. Quitting.\n", signal)
			tcp.Close()
			for conn, _ := range connections {
				close(conn)
			}
			fmt.Println("waiting to flush existing messages...")
			waits.Wait()
			close(brokenChannels)
			return
		}
	}
}
