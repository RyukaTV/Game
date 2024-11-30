package main

import (
	"encoding/json"
	"log"
	"net/http"
	"my-game/internal/game"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	gameState  = game.GameState{Players: make(map[string]*game.Player)}
	broadcast  = make(chan []byte) // Canal pour diffuser l'état du jeu
	connections = make(map[*websocket.Conn]string) // Associe une connexion WebSocket à un ID de joueur
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erreur de connexion WebSocket :", err)
		return
	}
	defer conn.Close()

	// Générer un ID de joueur unique (l'adresse du client)
	playerID := r.RemoteAddr
	connections[conn] = playerID
	gameState.AddPlayer(playerID)

	log.Printf("Nouveau joueur connecté : %s", playerID)

	broadcastGameState()
	for {
		var move struct {
			X int `json:"x"`
			Y int `json:"y"`
		}
		err := conn.ReadJSON(&move)
		if err != nil {
			log.Printf("Erreur de lecture (%s) : %v", playerID, err)
			break
		}

		// Mettre à jour la position du joueur
		gameState.UpdatePlayer(playerID, move.X, move.Y)
		broadcastGameState()
	}
	delete(connections, conn)
	gameState.RemovePlayer(playerID)
	broadcastGameState()
	log.Printf("Joueur déconnecté : %s", playerID)
}

func broadcastGameState() {
	state := gameState.GetSnapshot()
	data, err := json.Marshal(state)
	if err != nil {
		log.Println("Erreur de sérialisation de l'état :", err)
		return
	}
	broadcast <- data
}

func handleBroadcast() {
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

func main() {
	http.HandleFunc("/ws", handleConnections)
	go handleBroadcast()

	log.Println("Serveur démarré sur :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Erreur serveur :", err)
	}
}
