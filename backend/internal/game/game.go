package game

import (
	"log"
)

type GameState struct {
	Players map[string]*Player `json:"players"` // ID du joueur => Position
}

func (gs *GameState) AddPlayer(playerID string) {
	// Ajouter un joueur avec une position initiale
	gs.Players[playerID] = &Player{ID: playerID, X: 0, Y: 0}
	log.Printf("Joueur ajouté : %s", playerID)
}

func (gs *GameState) UpdatePlayer(playerID string, x, y int) {
	if player, exists := gs.Players[playerID]; exists {
		player.X = x
		player.Y = y
		log.Printf("Position mise à jour pour %s: (%d, %d)", playerID, x, y)
	} else {
		log.Printf("Le joueur %s n'existe pas", playerID)
	}
}

func (gs *GameState) RemovePlayer(playerID string) {
	delete(gs.Players, playerID)
	log.Printf("Joueur supprimé : %s", playerID)
}

// GetSnapshot retourne une copie de l'état actuel du jeu
func (gs *GameState) GetSnapshot() *GameState {
	// Crée une copie des joueurs pour éviter la modification directe
	playersCopy := make(map[string]*Player)
	for id, player := range gs.Players {
		playersCopy[id] = &Player{ID: player.ID, X: player.X, Y: player.Y}
	}
	return &GameState{Players: playersCopy}
}
