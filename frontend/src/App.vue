<template>
  <div id="app">
    <h1>Jeu Multijoueur</h1>
    <div>
      <span>Nombre de joueurs : {{ playerCount }}</span>
    </div>

    <div>
      <h2>Liste des joueurs</h2>
      <ul>
        <li v-for="(player, id) in players" :key="id">
          Joueur {{ player.id }}: Position ({{ player.x }}, {{ player.y }})
        </li>
      </ul>
    </div>

    <div class="grid">
      <div
        v-for="(cell, index) in grid"
        :key="index"
        :class="['cell', { player: cell !== null }]"
      >
        {{ cell }}
      </div>
    </div>

    <div class="controls">
      <button @click="move('up')">↑</button>
      <button @click="move('left')">←</button>
      <button @click="move('down')">↓</button>
      <button @click="move('right')">→</button>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      gridSize: 10,
      grid: Array(100).fill(null),
      playerPosition: { x: 0, y: 0 },
      socket: null,
      players: {},
    };
  },
  computed: {
    playerCount() {
      return Object.keys(this.players).length;
    },
  },
  methods: {
    // Vérifie si une cellule est occupée par un joueur
    isPlayer(cell) {
      return cell !== null;
    },
    updateGrid() {
      // Réinitialise la grille (remet toutes les cellules à null)
      this.grid = Array(this.gridSize * this.gridSize).fill(null);
      const playerIndex = this.playerPosition.y * this.gridSize + this.playerPosition.x;
      this.grid[playerIndex] = "P";  // Afficher "P" pour le joueur principal

      // Place les autres joueurs dans la grille
      Object.values(this.players).forEach(player => {
        const playerIndex = player.y * this.gridSize + player.x;
        this.grid[playerIndex] = player.id;  // Afficher l'ID du joueur dans la grille
      });
    },

    move(direction) {
      if (direction === "up" && this.playerPosition.y > 0) this.playerPosition.y--;
      if (direction === "down" && this.playerPosition.y < this.gridSize - 1) this.playerPosition.y++;
      if (direction === "left" && this.playerPosition.x > 0) this.playerPosition.x--;
      if (direction === "right" && this.playerPosition.x < this.gridSize - 1) this.playerPosition.x++;
      this.updateGrid();

      if (this.socket) {
        this.socket.send(JSON.stringify({ x: this.playerPosition.x, y: this.playerPosition.y }));
      }
    },
    connectToServer() {
      this.socket = new WebSocket("ws://localhost:8080/ws");
      this.socket.onmessage = (event) => {
        const serverState = JSON.parse(event.data);
        console.log("État synchronisé depuis le serveur :", serverState);
        this.players = serverState.players;
        this.updateGrid();
      };
    },
  },
  mounted() {
    this.updateGrid();
    this.connectToServer();
  },
};
</script>

<style>
.grid {
  display: grid;
  grid-template-columns: repeat(10, 50px);
  grid-template-rows: repeat(10, 50px);
  gap: 2px;
}
.cell {
  width: 50px;
  height: 50px;
  background-color: #ccc;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: bold;
  border: 1px solid #999;
}
.cell.player {
  background-color: #f00;
}
.controls {
  margin-top: 20px;
}
button {
  margin: 5px;
  padding: 10px;
  font-size: 16px;
}
</style>
