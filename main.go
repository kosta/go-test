package main

import (
	"./socketio"
	"http"
	"fmt"
	"syscall"
)

func main() {
	sio := socketio.NewSocketIO(nil)

	sio.OnConnect(func(c *socketio.Conn) {
		sio.Broadcast("connected: socket.io/" + c.String())
	})

	sio.OnDisconnect(func(c *socketio.Conn) {
		sio.BroadcastExcept(c, "disconnected: socket.io/" + c.String())
	})

	sio.OnMessage(func(c *socketio.Conn, msg socketio.Message) {
		sio.BroadcastExcept(c, msg.Data())
	})
	
	go func() {
	  count := 0
	  for {
	    sio.Broadcast(fmt.Sprintf("ping%d", count))
	    count++
	    syscall.Sleep(1e9)
	  }
	}()

	mux := sio.ServeMux()
	mux.Handle("/", http.FileServer("www/", "/"))
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
