package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"snake_game/internal/game"
)

var (
	Game     *game.Game
	GameLock sync.Mutex
)

// InitGame 初始化游戏实例
func InitGame(width, height int) {
	Game = game.NewGame(width, height)
	fmt.Println("Game initialized in API module")
}

// HandleState 处理获取游戏状态的请求
func HandleState(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling state request")
	GameLock.Lock()
	defer GameLock.Unlock()

	// 将游戏状态转换为JSON
	response := map[string]interface{}{
		"board":    Game.Board,
		"score":    Game.Score,
		"gameOver": Game.GameOver,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// HandleMove 处理蛇移动的请求
func HandleMove(w http.ResponseWriter, r *http.Request) {
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

	GameLock.Lock()
	defer GameLock.Unlock()

	// 更新方向
	switch move.Direction {
	case "up":
		if Game.Direction.Y != 1 { // 防止直接反向移动
			Game.Direction = game.Point{X: 0, Y: -1}
		}
	case "down":
		if Game.Direction.Y != -1 {
			Game.Direction = game.Point{X: 0, Y: 1}
		}
	case "left":
		if Game.Direction.X != 1 {
			Game.Direction = game.Point{X: -1, Y: 0}
		}
	case "right":
		if Game.Direction.X != -1 {
			Game.Direction = game.Point{X: 1, Y: 0}
		}
	}

	w.WriteHeader(http.StatusOK)
}

// HandleRestart 处理重启游戏的请求
func HandleRestart(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling restart request")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	GameLock.Lock()
	defer GameLock.Unlock()

	Game = game.NewGame(40, 20)
	w.WriteHeader(http.StatusOK)
}

// EnableCors 添加CORS头
func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Origin")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Content-Type", "application/json")
}

// CorsMiddleware CORS中间件
func CorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Handling %s request to %s\n", r.Method, r.URL.Path)
		EnableCors(&w)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}
