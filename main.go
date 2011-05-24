package main

import (
	"./socketio"
	"http"
	"fmt"
)

func main() {
	sio := socketio.NewSocketIO(nil)

	sio.OnConnect(func(c *socketio.Conn) {
		sio.Broadcast("connected: " + c.String())
	})

	sio.OnDisconnect(func(c *socketio.Conn) {
		sio.BroadcastExcept(c, "disconnected: " + c.String())
	})

	sio.OnMessage(func(c *socketio.Conn, msg socketio.Message) {
		sio.BroadcastExcept(c, msg.Data())
	})

	mux := sio.ServeMux()
	mux.Handle("/", http.FileServer("www/", "/"))
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
