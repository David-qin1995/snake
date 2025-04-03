<template>
  <div class="app">
    <div class="game-container">
      <h1>Snake Game</h1>
      <div class="game-board" :style="{ gridTemplateColumns: `repeat(${width}, 1fr)` }">
        <div
          v-for="(cell, index) in board"
          :key="index"
          :class="['cell', cell.type]"
        >
          <span v-if="cell.type === 'snake'" class="emoji">{{ getSnakePartEmoji(cell.part) }}</span>
          <span v-else-if="cell.type === 'food'" class="emoji">{{ getFoodEmoji(cell.foodType) }}</span>
        </div>
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

const FOOD_EMOJIS = {
  watermelon: '🍉',
  mango: '🥭',
  tomato: '🍅',
  apple: '🍎',
  coconut: '🥥'
};

const SNAKE_EMOJIS = {
  head: '🐍',
  body: '🟢',
  tail: '🟡'
};

export default {
  name: 'App',
  data() {
    return {
      width: 40,
      height: 20,
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
        console.log(`Moving ${direction}`);
        const response = await axios.post(`${API_BASE_URL}/move`, { direction });
        console.log('Move response:', response.status);
      } catch (error) {
        console.error('Error moving:', error);
        if (error.response) {
          console.error('Response data:', error.response.data);
          console.error('Response status:', error.response.status);
        }
      }
    },
    async updateGame() {
      try {
        console.log('Updating game state');
        const response = await axios.get(`${API_BASE_URL}/state`);
        console.log('State response:', response.data);
        const { board, score, gameOver } = response.data;
        
        // 更新游戏状态
        this.board = board.flat().map(cell => {
          if (cell.type === 'food') {
            return { type: 'food', foodType: cell.foodType };
          }
          return { type: cell.type };
        });
        this.score = score;
        this.gameOver = gameOver;
        
        if (gameOver) {
          console.log('Game over detected');
          this.stopGame();
        }
      } catch (error) {
        console.error('Error updating game:', error);
        if (error.response) {
          console.error('Response data:', error.response.data);
          console.error('Response status:', error.response.status);
        }
      }
    },
    async restartGame() {
      try {
        console.log('Restarting game');
        const response = await axios.post(`${API_BASE_URL}/restart`);
        console.log('Restart response:', response.status);
        this.score = 0;
        this.gameOver = false;
        this.lastDirection = null;
        this.startGame();
      } catch (error) {
        console.error('Error restarting game:', error);
        if (error.response) {
          console.error('Response data:', error.response.data);
          console.error('Response status:', error.response.status);
        }
      }
    },
    getFoodEmoji(foodType) {
      return FOOD_EMOJIS[foodType] || '🍎';
    },
    getSnakePartEmoji(part) {
      return SNAKE_EMOJIS[part] || '🟢';
    }
  }
};
</script>

<style>
html, body {
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.app {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  width: 100vw;
  height: 100vh;
  margin: 0;
  padding: 0;
  background-color: #f0f0f0;
  box-sizing: border-box;
  overflow: hidden;
}

.game-container {
  background-color: white;
  padding: 0.75rem;
  border-radius: 10px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  width: 90%;
  max-width: 800px;
  margin: 0 auto;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

h1 {
  text-align: center;
  color: #333;
  margin: 0 0 0.25rem 0;
  font-size: clamp(1rem, 2.5vw, 1.5rem);
}

.game-board {
  display: grid;
  gap: 1px;
  background-color: #ddd;
  padding: 3px;
  border-radius: 5px;
  width: 100%;
  aspect-ratio: 2;
  margin: 0 auto;
  box-sizing: border-box;
}

.cell {
  aspect-ratio: 1;
  background-color: white;
  border-radius: 1px;
  min-width: 0;
  min-height: 0;
  transition: background-color 0.2s;
  position: relative;
  overflow: hidden;
  display: flex;
  justify-content: center;
  align-items: center;
}

.emoji {
  font-size: 1.2em;
  line-height: 1;
  display: block;
  width: 100%;
  height: 100%;
  text-align: center;
  display: flex;
  justify-content: center;
  align-items: center;
}

.cell.snake {
  background-color: transparent;
}

.cell.snake .emoji {
  font-size: 1.2em;
}

.cell.food {
  background-color: transparent;
}

.game-info {
  margin-top: 0.25rem;
  text-align: center;
}

.score {
  font-size: clamp(0.7rem, 1.8vw, 1rem);
  font-weight: bold;
  margin-bottom: 0.15rem;
}

.game-over {
  color: #f44336;
  margin-top: 0.25rem;
  font-size: clamp(0.7rem, 1.8vw, 1rem);
}

.controls {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 0.25rem;
  gap: 0.15rem;
}

.horizontal-controls {
  display: flex;
  gap: 0.15rem;
}

button {
  padding: 0.2rem 0.4rem;
  font-size: clamp(0.7rem, 1.8vw, 1rem);
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 2px;
  cursor: pointer;
  transition: all 0.3s;
  min-width: 30px;
}

button:hover {
  background-color: #45a049;
  transform: scale(1.05);
}

button:active {
  transform: scale(0.95);
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .game-container {
    padding: 0.5rem;
  }
}

@media (max-width: 768px) {
  .game-container {
    padding: 0.4rem;
    width: 40%;
  }
  
  .controls {
    margin-top: 0.15rem;
  }
  
  button {
    padding: 0.15rem 0.3rem;
    font-size: 0.9rem;
  }
}

@media (max-width: 480px) {
  .app {
    padding: 3px;
  }
  
  .game-container {
    padding: 0.2rem;
    width: 50%;
  }
  
  .game-board {
    padding: 2px;
  }
  
  .controls {
    gap: 0.1rem;
  }
  
  button {
    padding: 0.15rem 0.25rem;
    font-size: 0.7rem;
    min-width: 25px;
  }
}

/* 横屏优化 */
@media (max-height: 600px) {
  .app {
    padding: 3px;
  }
  
  .game-container {
    padding: 0.2rem;
  }
  
  .game-board {
    aspect-ratio: 2;
  }
  
  h1 {
    margin-bottom: 0.15rem;
  }
  
  .game-info {
    margin-top: 0.15rem;
  }
  
  .controls {
    margin-top: 0.15rem;
  }
}
</style> 