package util

import "github.com/gorilla/websocket"

func KeyInMap(connMap *map[*websocket.Conn]*string, key *string) *bool {
	var result bool = false
	for _, currentKey := range *connMap {
		if *key == *currentKey {
			result = true
			return &result
		}
	}
	return &result
}

func GetKeys(connMap *map[*websocket.Conn]*string) []string {
	var conns []string
	for _, conn := range *connMap {
		conns = append(conns, *conn)
	}
	return conns
}
