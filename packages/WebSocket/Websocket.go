package WebSocket

import (
	"chat/packages/Broadcaster"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func Manager(
	conn *websocket.Conn,
	broacPackChan *chan *Broadcaster.BroadcasterPacket,
	wsConnMap *map[*websocket.Conn]*string,
) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	for {
		var wsPayload Broadcaster.WebSocketPayload
		err := conn.ReadJSON(&wsPayload)
		if err != nil {
			log.Println(err)
		}

		*broacPackChan <- &Broadcaster.BroadcasterPacket{
			Payload:      &wsPayload,
			WsConnection: conn,
		}
	}
}

func Handler(
	w *http.ResponseWriter,
	r *http.Request,
	upgrader *websocket.Upgrader,
	broadPackChan *chan *Broadcaster.BroadcasterPacket,
	wsConnMap *map[*websocket.Conn]*string,
) {
	conn, err := upgrader.Upgrade(*w, r, nil)
	if err != nil {
		log.Println(err)
	}
	go Manager(
		conn,
		broadPackChan,
		wsConnMap,
	)

}
