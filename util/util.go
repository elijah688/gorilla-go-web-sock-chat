package util

import "github.com/gorilla/websocket"

func KeyInMap(connMap *map[string]*websocket.Conn, key *string) *bool {
	var result bool = false
	for currentKey := range *connMap {
		if *key == currentKey {
			result = true
			return &result
		}
	}
	return &result
}

func GetKeys(connMap *map[string]*websocket.Conn) []string {
	var conns []string
	for conn := range *connMap {
		conns = append(conns, conn)
	}
	return conns
}
