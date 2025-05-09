package game

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Cell struct {
	Type     string `json:"type"`
	FoodType string `json:"foodType,omitempty"`
	Part     string `json:"part,omitempty"`
}

type Game struct {
	Width     int
	Height    int
	Snake     []Point
	Direction Point
	Food      Point
	foodType  string
	Score     int
	GameOver  bool
	Board     [][]Cell
}

func NewGame(width, height int) *Game {
	game := &Game{
		Width:     width,
		Height:    height,
		Board:     make([][]Cell, height),
		Snake:     make([]Point, 3), // 初始有头、身体和尾
		Direction: Point{X: 1, Y: 0},
		Score:     0,
		GameOver:  false,
	}

	// 初始化游戏板
	for i := range game.Board {
		game.Board[i] = make([]Cell, width)
		for j := range game.Board[i] {
			game.Board[i][j] = Cell{Type: "empty"}
		}
	}

	// 设置蛇的初始位置（水平放置）
	centerY := height / 2
	centerX := width / 2
	game.Snake[0] = Point{X: centerX, Y: centerY}     // 蛇头
	game.Snake[1] = Point{X: centerX - 1, Y: centerY} // 蛇身
	game.Snake[2] = Point{X: centerX - 2, Y: centerY} // 蛇尾

	// 在游戏板上标记蛇的位置
	game.Board[game.Snake[0].Y][game.Snake[0].X] = Cell{Type: "snake", Part: "head"}
	game.Board[game.Snake[1].Y][game.Snake[1].X] = Cell{Type: "snake", Part: "body"}
	game.Board[game.Snake[2].Y][game.Snake[2].X] = Cell{Type: "snake", Part: "tail"}

	// 生成食物
	game.GenerateFood()
	return game
}

func (g *Game) GenerateFood() {
	fmt.Println("Generating food")
	foodTypes := []string{"watermelon", "mango", "tomato", "apple", "coconut"}
	rand.Seed(time.Now().UnixNano())
	g.foodType = foodTypes[rand.Intn(len(foodTypes))]

	// 找到一个空位置放置食物
	emptyPositions := []Point{}
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if g.Board[y][x].Type == "empty" {
				emptyPositions = append(emptyPositions, Point{X: x, Y: y})
			}
		}
	}

	if len(emptyPositions) > 0 {
		// 随机选择一个空位置
		position := emptyPositions[rand.Intn(len(emptyPositions))]
		g.Board[position.Y][position.X] = Cell{Type: "food", FoodType: g.foodType}
		g.Food = position
		fmt.Printf("Food generated at position: (%d, %d)\n", position.X, position.Y)
	} else {
		fmt.Println("No empty positions found for food")
	}
}

func (g *Game) Update() {
	if g.GameOver {
		return
	}

	// 计算新的蛇头位置
	newHead := Point{
		X: g.Snake[0].X + g.Direction.X,
		Y: g.Snake[0].Y + g.Direction.Y,
	}

	// 检查是否撞墙
	if newHead.X < 0 || newHead.X >= g.Width || newHead.Y < 0 || newHead.Y >= g.Height {
		g.GameOver = true
		return
	}

	// 检查是否撞到自己
	for i := 1; i < len(g.Snake); i++ {
		if newHead.X == g.Snake[i].X && newHead.Y == g.Snake[i].Y {
			g.GameOver = true
			return
		}
	}

	// 移动蛇
	oldHead := g.Snake[0]
	g.Snake = append([]Point{newHead}, g.Snake...)

	// 检查是否吃到食物
	if newHead.X == g.Food.X && newHead.Y == g.Food.Y {
		g.Score++
		g.GenerateFood()
	} else {
		// 如果没有吃到食物，移除蛇尾
		tail := g.Snake[len(g.Snake)-1]
		g.Board[tail.Y][tail.X] = Cell{Type: "empty"}
		g.Snake = g.Snake[:len(g.Snake)-1]
	}

	// 更新游戏板
	g.Board[oldHead.Y][oldHead.X] = Cell{Type: "snake", Part: "body"}
	g.Board[newHead.Y][newHead.X] = Cell{Type: "snake", Part: "head"}
	if len(g.Snake) > 1 {
		g.Board[g.Snake[1].Y][g.Snake[1].X] = Cell{Type: "snake", Part: "body"}
	}
	if len(g.Snake) > 2 {
		g.Board[g.Snake[2].Y][g.Snake[2].X] = Cell{Type: "snake", Part: "body"}
	}
	g.Board[g.Snake[len(g.Snake)-1].Y][g.Snake[len(g.Snake)-1].X] = Cell{Type: "snake", Part: "tail"}
}

func (g *Game) Draw() {
	// Clear screen
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Draw game board
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			// Check if current position is snake
			isSnake := false
			for _, p := range g.Snake {
				if p.X == x && p.Y == y {
					isSnake = true
					break
				}
			}

			if isSnake {
				fmt.Print("■")
			} else if x == g.Food.X && y == g.Food.Y {
				fmt.Print("★")
			} else {
				fmt.Print("·")
			}
		}
		fmt.Println()
	}

	fmt.Printf("Score: %d\n", g.Score)
	if g.GameOver {
		fmt.Println("Game Over!")
	}
}
