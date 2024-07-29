package game

import (
	"math/rand"
	"slices"
	"time"

	"github.com/Richtermnd/snake/terminal"
	"github.com/eiannone/keyboard"
)

type point struct {
	X int
	Y int
}

const FPS = 15

// Styling.
const (
	EMPTY = ' '
	SNAKE = '$'
	FOOD  = '@'
)

var foodMap map[point]bool = make(map[point]bool)
var snake []point

// Buffer to store snake tail after eating food
// to add this to snake on next frame.
var buffer []point

// Move direction
var dir point

// Screen width and height
var width int
var height int

// Start start the game
func Start() {
	width, height = terminal.GetTerminalSize()
	ticker := time.NewTicker(time.Second / FPS)
	terminal.HideCursor()
	// quit channel to handle Ctrc+C
	quit := make(chan struct{}, 1)
	go handleKeyboard(quit)
	newGame()
	for {
		select {
		case <-ticker.C:
		case <-quit:
			terminal.ResetCursor()
			terminal.ShowCursor()
			return
		}
		move()
		eat()
		if isGameOver() {
			newGame()
		}
		render()
	}
}

func handleKeyboard(quit chan struct{}) {
	keys, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	for {
		event := <-keys
		switch event.Key {
		case keyboard.KeyArrowUp:
			changeDir(point{0, -1})
		case keyboard.KeyArrowDown:
			changeDir(point{0, 1})
		case keyboard.KeyArrowLeft:
			changeDir(point{-1, 0})
		case keyboard.KeyArrowRight:
			changeDir(point{1, 0})
		case keyboard.KeyCtrlC:
			quit <- struct{}{}
		}
	}
}

func changeDir(newDir point) {
	if newDir.X == 0 && newDir.Y == -dir.Y {
		return
	}
	if newDir.Y == 0 && newDir.X == -dir.X {
		return
	}
	dir = newDir
}

func move() {
	for i := len(snake) - 1; i > 0; i-- {
		snake[i] = snake[i-1]
	}
	snake[0].X += dir.X
	snake[0].Y += dir.Y
	if len(buffer) > 0 {
		snake = append(snake, buffer[0])
		buffer = buffer[1:] // God, save GC
	}
}

func eat() {
	if len(foodMap) == 0 {
		placeFood()
	}
	_, ok := foodMap[snake[0]]
	if !ok {
		return
	}
	delete(foodMap, snake[0])
	buffer = append(buffer, snake[len(snake)-1])
}

func placeFood() {
	var x, y int
	x = rand.Intn(width)
	y = rand.Intn(height)
	for slices.Contains(snake, point{x, y}) {
		x = rand.Intn(width)
		y = rand.Intn(height)
	}
	foodMap[point{x, y}] = true
}

func newGame() {
	snake = []point{{width / 2, height / 2}}
	dir.X = 1
	dir.Y = 0
	buffer = make([]point, 0)
	clear(foodMap)
}

func isGameOver() bool {
	head := snake[0]
	if head.X >= width || head.X < 0 {
		return true
	}
	if head.Y >= height || head.Y < 0 {
		return true
	}
	return slices.Contains(snake[1:], head)
}

func render() {
	screen := make([]byte, width*height)
	for i := range screen {
		screen[i] = EMPTY
	}
	for _, point := range snake {
		screen[point.Y*width+point.X] = SNAKE
	}
	for point := range foodMap {
		screen[point.Y*width+point.X] = FOOD
	}
	terminal.ResetAndWrite(screen)
}
