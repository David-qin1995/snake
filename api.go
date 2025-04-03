package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var gameLock sync.Mutex

func handleState(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling state request")
	gameLock.Lock()
	defer gameLock.Unlock()

	// 将游戏状态转换为JSON
	response := map[string]interface{}{
		"board":    game.Board,
		"score":    game.Score,
		"gameOver": game.GameOver,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func handleMove(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling move request")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var move struct {
		Direction string `json:"direction"`
	}
	if err := json.NewDecoder(r.Body).Decode(&move); err != nil {
		fmt.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gameLock.Lock()
	defer gameLock.Unlock()

	// 更新方向
	switch move.Direction {
	case "up":
		if game.Direction.Y != 1 { // 防止直接反向移动
			game.Direction = Point{X: 0, Y: -1}
		}
	case "down":
		if game.Direction.Y != -1 {
			game.Direction = Point{X: 0, Y: 1}
		}
	case "left":
		if game.Direction.X != 1 {
			game.Direction = Point{X: -1, Y: 0}
		}
	case "right":
		if game.Direction.X != -1 {
			game.Direction = Point{X: 1, Y: 0}
		}
	}

	w.WriteHeader(http.StatusOK)
}

func handleRestart(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling restart request")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameLock.Lock()
	defer gameLock.Unlock()

	game = NewGame(40, 20)
	w.WriteHeader(http.StatusOK)
}
