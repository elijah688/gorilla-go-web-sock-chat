package Broadcaster

import (
	"chat/packages/Util"
	"fmt"
	"log"
	"sync"

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
		(*broadcasterPacket).Payload.Clients = Util.GetClients(wsConnMap)
		err := conn.WriteJSON(*broadcasterPacket.Payload)
		if err != nil {
			log.Println(err)
			conn.Close()
			delete(*wsConnMap, conn)
		}

	}
	mutex.Unlock()
}

func SafeBroadcastUserHasDupName(
	wsConnMap *map[*websocket.Conn]*string,
	mutex *sync.Mutex,
	broadcasterPacket *BroadcasterPacket,
) {
	broadcasterPacket.Payload.Message = fmt.Sprintf("The name \"%s\" is taken!", broadcasterPacket.Payload.Name)

	mutex.Lock()
	delete(*wsConnMap, broadcasterPacket.WsConnection)
	broadcasterPacket.Payload.Clients = Util.GetClients(wsConnMap)
	mutex.Unlock()

	err := (*broadcasterPacket).WsConnection.WriteJSON(broadcasterPacket.Payload)
	if err != nil {
		log.Println(err)
	}
}

func Broadcaster(
	broacPackChan *chan *BroadcasterPacket,
	wsConnMap *map[*websocket.Conn]*string,
	mutex *sync.Mutex,
) {
	for {
		var broadcasterPacket BroadcasterPacket = *<-*broacPackChan

		switch broadcasterPacket.Payload.Action {
		case "connected":
			if *Util.KeyInMap(wsConnMap, &broadcasterPacket.Payload.Name) {
				SafeBroadcastUserHasDupName(
					wsConnMap,
					mutex,
					&broadcasterPacket,
				)
			} else {

				mutex.Lock()
				(*wsConnMap)[broadcasterPacket.WsConnection] = &broadcasterPacket.Payload.Name
				mutex.Unlock()

				broadcasterPacket.Payload.Message = fmt.Sprintf("%s, connected....", broadcasterPacket.Payload.Name)

				SafeBroadcast(
					wsConnMap,
					mutex,
					&broadcasterPacket,
				)
			}

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
