package game

import (
	"encoding/json"
	"log"
	"github.com/gorilla/websocket"
)

func BroadcastGameState(gameState *GameState, broadcast chan []byte, connections map[*websocket.Conn]string) {
	state := gameState.GetSnapshot()
	data, err := json.Marshal(state)
	if err != nil {
		log.Println("Erreur de sérialisation de l'état :", err)
		return
	}
	broadcast <- data
}

func HandleBroadcast(broadcast chan []byte, connections map[*websocket.Conn]string) {
	for {
		data := <-broadcast
		for conn := range connections {
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println("Erreur de diffusion :", err)
				conn.Close()
				delete(connections, conn)
			}
		}
	}
}
