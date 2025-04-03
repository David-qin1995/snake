package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

var (
	game     *Game
	gameLock sync.Mutex
	lastMove time.Time
)

type GameState struct {
	Board    [][]string `json:"board"`
	Score    int        `json:"score"`
	GameOver bool       `json:"gameOver"`
}

func init() {
	game = NewGame(20, 10)
	lastMove = time.Now()

	// Start auto-movement goroutine
	go func() {
		for {
			time.Sleep(200 * time.Millisecond) // Move every 200ms
			gameLock.Lock()
			if !game.gameOver {
				game.update()
			}
			gameLock.Unlock()
		}
	}()
}

func handleMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var move struct {
		Direction string `json:"direction"`
	}
	if err := json.NewDecoder(r.Body).Decode(&move); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gameLock.Lock()
	defer gameLock.Unlock()

	// Only allow direction changes if enough time has passed since last move
	if time.Since(lastMove) < 100*time.Millisecond {
		sendGameState(w)
		return
	}

	switch move.Direction {
	case "up":
		if game.direction.Y != 1 {
			game.direction = Point{0, -1}
		}
	case "down":
		if game.direction.Y != -1 {
			game.direction = Point{0, 1}
		}
	case "left":
		if game.direction.X != 1 {
			game.direction = Point{-1, 0}
		}
	case "right":
		if game.direction.X != -1 {
			game.direction = Point{1, 0}
		}
	}

	lastMove = time.Now()
	sendGameState(w)
}

func handleGameState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameLock.Lock()
	defer gameLock.Unlock()
	sendGameState(w)
}

func sendGameState(w http.ResponseWriter) {
	state := GameState{
		Board:    make([][]string, game.height),
		Score:    game.score,
		GameOver: game.gameOver,
	}

	for y := 0; y < game.height; y++ {
		state.Board[y] = make([]string, game.width)
		for x := 0; x < game.width; x++ {
			state.Board[y][x] = "empty"
		}
	}

	// Mark snake positions
	for _, p := range game.snake {
		state.Board[p.Y][p.X] = "snake"
	}

	// Mark food position
	state.Board[game.food.Y][game.food.X] = "food"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func handleRestart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameLock.Lock()
	defer gameLock.Unlock()

	game = NewGame(20, 10)
	lastMove = time.Now()
	sendGameState(w)
}
