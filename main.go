package main

import (
	"chat/packages/Broadcaster"
	"chat/packages/WebSocket"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {
	var wsConnMap *map[*websocket.Conn]*string = &map[*websocket.Conn]*string{}
	var broadPackChan chan *Broadcaster.BroadcasterPacket = make(chan *Broadcaster.BroadcasterPacket)
	var upgrader *websocket.Upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	var mutex = &sync.Mutex{}
	go Broadcaster.Broadcaster(&broadPackChan, wsConnMap, mutex)

	mux := mux.NewRouter()

	mux.HandleFunc(
		"/ws",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			WebSocket.Handler(&w, r, upgrader, &broadPackChan, wsConnMap)
		}).Methods("Get")
	fs := http.FileServer(http.Dir("./static/"))

	mux.PathPrefix("/").Handler(fs)

	http.ListenAndServe(":8080", mux)
}
