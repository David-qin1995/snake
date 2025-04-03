<template>
  <div class="app">
    <div class="game-container">
      <h1>Snake Game</h1>
      <div class="game-board" :style="{ gridTemplateColumns: `repeat(${width}, 1fr)` }">
        <div
          v-for="(cell, index) in board"
          :key="index"
          :class="['cell', cell.type]"
        ></div>
      </div>
      <div class="game-info">
        <div class="score">Score: {{ score }}</div>
        <div v-if="gameOver" class="game-over">
          Game Over!
          <button @click="restartGame">Play Again</button>
        </div>
      </div>
      <div class="controls">
        <button @click="move('up')">↑</button>
        <div class="horizontal-controls">
          <button @click="move('left')">←</button>
          <button @click="move('right')">→</button>
        </div>
        <button @click="move('down')">↓</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api';

export default {
  name: 'App',
  data() {
    return {
      width: 20,
      height: 10,
      board: [],
      score: 0,
      gameOver: false,
      gameInterval: null,
      lastDirection: null
    };
  },
  created() {
    this.initializeBoard();
    this.startGame();
  },
  beforeUnmount() {
    this.stopGame();
  },
  methods: {
    initializeBoard() {
      this.board = Array(this.width * this.height).fill({ type: 'empty' });
    },
    startGame() {
      this.gameInterval = setInterval(this.updateGame, 100);
      document.addEventListener('keydown', this.handleKeyPress);
    },
    stopGame() {
      clearInterval(this.gameInterval);
      document.removeEventListener('keydown', this.handleKeyPress);
    },
    handleKeyPress(event) {
      if (this.gameOver) return;
      
      switch(event.key) {
        case 'ArrowUp':
          if (this.lastDirection !== 'down') {
            this.move('up');
            this.lastDirection = 'up';
          }
          break;
        case 'ArrowDown':
          if (this.lastDirection !== 'up') {
            this.move('down');
            this.lastDirection = 'down';
          }
          break;
        case 'ArrowLeft':
          if (this.lastDirection !== 'right') {
            this.move('left');
            this.lastDirection = 'left';
          }
          break;
        case 'ArrowRight':
          if (this.lastDirection !== 'left') {
            this.move('right');
            this.lastDirection = 'right';
          }
          break;
      }
    },
    async move(direction) {
      try {
        await axios.post(`${API_BASE_URL}/move`, { direction });
      } catch (error) {
        console.error('Error moving:', error);
      }
    },
    async updateGame() {
      try {
        const response = await axios.get(`${API_BASE_URL}/state`);
        const { board, score, gameOver } = response.data;
        
        // Flatten the 2D board array
        this.board = board.flat().map(type => ({ type }));
        this.score = score;
        this.gameOver = gameOver;
        
        if (gameOver) {
          this.stopGame();
        }
      } catch (error) {
        console.error('Error updating game:', error);
      }
    },
    async restartGame() {
      try {
        await axios.post(`${API_BASE_URL}/restart`);
        this.score = 0;
        this.gameOver = false;
        this.lastDirection = null;
        this.startGame();
      } catch (error) {
        console.error('Error restarting game:', error);
      }
    }
  }
};
</script>

<style>
.app {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f0f0f0;
}

.game-container {
  background-color: white;
  padding: 2rem;
  border-radius: 10px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

h1 {
  text-align: center;
  color: #333;
  margin-bottom: 1rem;
}

.game-board {
  display: grid;
  gap: 2px;
  background-color: #ddd;
  padding: 10px;
  border-radius: 5px;
}

.cell {
  aspect-ratio: 1;
  background-color: white;
  border-radius: 2px;
}

.cell.snake {
  background-color: #4CAF50;
}

.cell.food {
  background-color: #f44336;
}

.game-info {
  margin-top: 1rem;
  text-align: center;
}

.score {
  font-size: 1.2rem;
  font-weight: bold;
  margin-bottom: 0.5rem;
}

.game-over {
  color: #f44336;
  margin-top: 1rem;
}

.controls {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 1rem;
  gap: 0.5rem;
}

.horizontal-controls {
  display: flex;
  gap: 0.5rem;
}

button {
  padding: 0.5rem 1rem;
  font-size: 1.2rem;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s;
}

button:hover {
  background-color: #45a049;
}

button:active {
  transform: scale(0.95);
}
</style> 