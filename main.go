package main

import (
	"chat/util"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WebSocketPayload struct {
	Name    string   `json:"name"`
	Action  string   `json:"action"`
	Message string   `json:"message"`
	Clients []string `json:"clients"`
}

type BroadcasterPacket struct {
	WsConnection *websocket.Conn
	Payload      *WebSocketPayload
}

func SafeBroadcast(
	wsConnMap *map[*websocket.Conn]*string,
	mutex *sync.Mutex,
	broadcasterPacket *BroadcasterPacket,
) {
	mutex.Lock()
	for conn := range *wsConnMap {
		(*broadcasterPacket).Payload.Clients = util.GetKeys(wsConnMap)
		err := conn.WriteJSON(*broadcasterPacket.Payload)
		if err != nil {
			log.Println(err)
			conn.Close()
			delete(*wsConnMap, conn)
		}
	}
	mutex.Unlock()
}

func broadcaster(
	broacPackChan *chan *BroadcasterPacket,
	wsConnMap *map[*websocket.Conn]*string,
	mutex *sync.Mutex,
) {
	for {
		var broadcasterPacket BroadcasterPacket = *<-*broacPackChan

		switch broadcasterPacket.Payload.Action {
		case "connected":
			if !*util.KeyInMap(wsConnMap, &broadcasterPacket.Payload.Name) {

				mutex.Lock()
				(*wsConnMap)[broadcasterPacket.WsConnection] = &broadcasterPacket.Payload.Name
				mutex.Unlock()

				broadcasterPacket.Payload.Message = fmt.Sprintf("%s, connected....", broadcasterPacket.Payload.Name)
			}
			SafeBroadcast(
				wsConnMap,
				mutex,
				&broadcasterPacket,
			)

		case "broadcast":
			SafeBroadcast(
				wsConnMap,
				mutex,
				&broadcasterPacket,
			)
		case "disconnected":
			delete(*wsConnMap, broadcasterPacket.WsConnection)
			SafeBroadcast(
				wsConnMap,
				mutex,
				&broadcasterPacket,
			)
		}
	}
}

func ConnecitonManager(
	conn *websocket.Conn,
	broacPackChan *chan *BroadcasterPacket,
	wsConnMap *map[*websocket.Conn]*string,
) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	for {
		var wsPayload WebSocketPayload
		err := conn.ReadJSON(&wsPayload)
		if err != nil {
			log.Println(err)
		}

		*broacPackChan <- &BroadcasterPacket{
			Payload:      &wsPayload,
			WsConnection: conn,
		}
	}
}

func Handler(
	w *http.ResponseWriter,
	r *http.Request,
	upgrader *websocket.Upgrader,
	broadPackChan *chan *BroadcasterPacket,
	wsConnMap *map[*websocket.Conn]*string,
) {
	conn, err := upgrader.Upgrade(*w, r, nil)
	if err != nil {
		log.Println(err)
	}
	go ConnecitonManager(
		conn,
		broadPackChan,
		wsConnMap,
	)

}
func main() {
	var wsConnMap *map[*websocket.Conn]*string = &map[*websocket.Conn]*string{}
	var broadPackChan chan *BroadcasterPacket = make(chan *BroadcasterPacket)
	var upgrader *websocket.Upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	var mutex = &sync.Mutex{}
	go broadcaster(&broadPackChan, wsConnMap, mutex)

	mux := mux.NewRouter()

	mux.HandleFunc(
		"/ws",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			Handler(&w, r, upgrader, &broadPackChan, wsConnMap)
		}).Methods("Get")
	fs := http.FileServer(http.Dir("./static/"))

	mux.PathPrefix("/").Handler(fs)

	http.ListenAndServe(":8080", mux)
}
