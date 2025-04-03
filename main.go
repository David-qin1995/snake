package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type Point struct {
	X, Y int
}

type Game struct {
	width     int
	height    int
	snake     []Point
	food      Point
	direction Point
	score     int
	gameOver  bool
}

func NewGame(width, height int) *Game {
	game := &Game{
		width:     width,
		height:    height,
		snake:     []Point{{width / 2, height / 2}},
		direction: Point{1, 0},
		score:     0,
		gameOver:  false,
	}
	game.generateFood()
	return game
}

func (g *Game) generateFood() {
	for {
		g.food = Point{
			X: rand.Intn(g.width),
			Y: rand.Intn(g.height),
		}
		// Make sure food doesn't spawn on snake
		valid := true
		for _, p := range g.snake {
			if p.X == g.food.X && p.Y == g.food.Y {
				valid = false
				break
			}
		}
		if valid {
			break
		}
	}
}

func (g *Game) update() {
	if g.gameOver {
		return
	}

	// Calculate new head position
	head := g.snake[0]
	newHead := Point{
		X: (head.X + g.direction.X + g.width) % g.width,
		Y: (head.Y + g.direction.Y + g.height) % g.height,
	}

	// Check for collisions with self
	for _, p := range g.snake {
		if p.X == newHead.X && p.Y == newHead.Y {
			g.gameOver = true
			return
		}
	}

	// Add new head
	g.snake = append([]Point{newHead}, g.snake...)

	// Check if food is eaten
	if newHead.X == g.food.X && newHead.Y == g.food.Y {
		g.score++
		g.generateFood()
	} else {
		// Remove tail if no food was eaten
		g.snake = g.snake[:len(g.snake)-1]
	}
}

func (g *Game) draw() {
	// Clear screen
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Draw game board
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			// Check if current position is snake
			isSnake := false
			for _, p := range g.snake {
				if p.X == x && p.Y == y {
					isSnake = true
					break
				}
			}

			if isSnake {
				fmt.Print("■")
			} else if x == g.food.X && y == g.food.Y {
				fmt.Print("★")
			} else {
				fmt.Print("·")
			}
		}
		fmt.Println()
	}

	fmt.Printf("Score: %d\n", g.score)
	if g.gameOver {
		fmt.Println("Game Over!")
	}
}

func main() {
	// Set random seed
	rand.Seed(time.Now().UnixNano())

	// Set up HTTP routes with CORS headers
	http.HandleFunc("/api/move", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		handleMove(w, r)
	})
	http.HandleFunc("/api/state", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		handleGameState(w, r)
	})
	http.HandleFunc("/api/restart", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		handleRestart(w, r)
	})

	// Serve static files from the frontend build directory
	fs := http.FileServer(http.Dir("snake-game-frontend/dist"))
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received request for: %s\n", r.URL.Path)
		fs.ServeHTTP(w, r)
	}))

	// Start HTTP server
	fmt.Println("Starting server on http://localhost:8080")
	fmt.Println("Static files will be served from: snake-game-frontend/dist")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
