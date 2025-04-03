package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"snake_game/internal/api"
)

func main() {
	// 初始化游戏
	api.InitGame(40, 20)
	fmt.Println("Game initialized")

	// 启动游戏更新循环
	go func() {
		fmt.Println("Starting game update loop")
		for {
			if !api.Game.GameOver {
				api.GameLock.Lock()
				api.Game.Update()
				api.GameLock.Unlock()
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	// 创建多路复用器
	mux := http.NewServeMux()

	// API 路由
	mux.HandleFunc("/api/move", api.CorsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Move request received")
		api.HandleMove(w, r)
	}))
	mux.HandleFunc("/api/restart", api.CorsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Restart request received")
		api.HandleRestart(w, r)
	}))
	mux.HandleFunc("/api/state", api.CorsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("State request received")
		api.HandleState(w, r)
	}))

	// 提供静态文件
	fs := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received request for path: %s\n", r.URL.Path)

		if r.Method == "OPTIONS" {
			api.EnableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.URL.Path == "/" {
			fmt.Println("Serving index.html")
			http.ServeFile(w, r, "./web/static/index.html")
			return
		}

		// 检查文件是否存在
		filePath := "./web/static" + r.URL.Path
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("File not found: %s\n", filePath)
			http.ServeFile(w, r, "./web/static/index.html")
			return
		}

		fmt.Printf("Serving file: %s\n", filePath)
		fs.ServeHTTP(w, r)
	}))

	// 创建服务器
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Starting server on http://localhost:8080")
	fmt.Println("Static files will be served from: ./web/static")
	fmt.Println("Server is running...")

	// 启动服务器
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
