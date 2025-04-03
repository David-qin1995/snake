# Snake Game

这是一个简单的贪吃蛇游戏，包含后端和前端。

## 项目结构

```
snake_game/
├── cmd/                # 命令行入口
│   └── server/         # 服务器入口
│       └── main.go     # 主程序
├── internal/           # 内部包
│   ├── api/            # API处理
│   │   └── handlers.go # API处理函数
│   └── game/           # 游戏逻辑
│       └── game.go     # 游戏核心逻辑
├── web/                # Web相关文件
│   └── static/         # 静态文件
│       ├── css/        # CSS文件
│   │   └── js/         # JavaScript文件
│   │   └── index.html  # 主HTML文件
├── go.mod              # Go模块文件
└── README.md           # 项目说明文件
```

## 如何运行

### 后端

```bash
# 在项目根目录下运行
go run cmd/server/main.go
```

### 前端开发

```bash
# 进入frontend目录
cd snake-game-frontend

# 安装依赖
npm install

# 运行开发服务器
npm run serve

# 构建生产环境
npm run build
```

## 如何玩

1. 打开浏览器，访问 http://localhost:8080
2. 使用方向键或屏幕上的按钮控制蛇的移动
3. 吃到食物会增加蛇的长度和分数
4. 碰到墙壁或自己的身体会导致游戏结束

## 技术栈

- 后端: Go
- 前端: Vue.js, HTML, CSS

## Game Rules

- Control the snake to eat the food (★)
- Each time you eat food, your score increases and the snake grows longer
- Avoid hitting the walls or yourself
- The game ends if you hit a wall or yourself

## Game Symbols

- ■: Snake body
- ★: Food
- ·: Empty space 