package main

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"math/rand"
)

const (
	defaultColor = termbox.ColorDefault
	bgColor      = termbox.ColorDefault
	snakeColor   = termbox.ColorGreen
)

var (
	UP    = Point{X: 0, Y: -1}
	DOWN  = Point{X: 0, Y: 1}
	LEFT  = Point{X: -1, Y: 0}
	RIGHT = Point{X: 1, Y: 0}
)

type Point struct {
	X, Y int
}

func (p Point) ToRune() rune {
	if p == UP {
		return '▲'
	}
	if p == DOWN {
		return '▼'
	}
	if p == LEFT {
		return '◀'
	}
	if p == RIGHT {
		return '▶'
	}
	return '●'
}

type Game struct {
	Snake         []Point //координаты змейки, первый элемент голова
	Food          Point   // координаты еды
	Malware       []Point // позиции вреда
	Dir           Point   // направление змейки
	Score         int
	Level         int
	GameOver      bool
	Width, Height int
	Quit          chan struct{} //сигнал для окончания игры
}

func NewGame(width int, height int) *Game {
	snake := make([]Point, 1)
	snake[0] = Point{width / 2, height / 2}
	malware := make([]Point, 0)
	dir := Point{1, 0}
	quit := make(chan struct{})
	game := Game{
		Snake:    snake,
		Malware:  malware,
		Dir:      dir,
		Score:    0,
		Level:    1,
		GameOver: false,
		Width:    width,
		Height:   height,
		Quit:     quit,
	}
	game.placeMalware()
	game.placeFood()
	return &game
}

func (g *Game) handleInput(ev termbox.Event) {
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyArrowUp:
			if g.Dir == DOWN {
				return
			}
			g.Dir = UP
		case termbox.KeyArrowDown:
			if g.Dir == UP {
				return
			}
			g.Dir = DOWN
		case termbox.KeyArrowLeft:
			if g.Dir == RIGHT {
				return
			}
			g.Dir = LEFT
		case termbox.KeyArrowRight:
			if g.Dir == LEFT {
				return
			}
			g.Dir = RIGHT
		case termbox.KeyEsc:
			close(g.Quit)
		}
		switch ev.Ch {
		case 'w':
			if g.Dir == DOWN {
				return
			}
			g.Dir = UP
		case 's':
			if g.Dir == UP {
				return
			}
			g.Dir = DOWN
		case 'a':
			if g.Dir == RIGHT {
				return
			}
			g.Dir = LEFT
		case 'd':
			if g.Dir == LEFT {
				return
			}
			g.Dir = RIGHT
		case 'q':
			close(g.Quit)
		}
	}
}

func (g *Game) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	left, top, bottom := g.getBoundaries()
	g.renderArea(top, bottom, left)
	g.renderSnake(left, bottom)
	g.renderInfo(left, bottom)
	termbox.SetCell(left+g.Food.X, top+g.Food.Y, '●', termbox.ColorGreen, bgColor)
	for _, m := range g.Malware {
		termbox.SetCell(left+m.X, top+m.Y, '✗', termbox.ColorRed, bgColor)
	}
	termbox.Flush()
}

func (g *Game) getBoundaries() (int, int, int) {
	w, h := termbox.Size()
	midY := h / 2
	left := (w - g.Width) / 2
	top := midY - (g.Height / 2)
	bottom := midY + (g.Height / 2) + 1
	return left, top, bottom
}

func (g *Game) renderSnake(left, bottom int) {
	for _, b := range g.Snake {
		termbox.SetCell(left+b.X, bottom-b.Y, g.Dir.ToRune(), snakeColor, bgColor)
	}
}

func (g *Game) renderArea(top, bottom, left int) {
	for i := top; i < bottom; i++ {
		termbox.SetCell(left-1, i, '│', defaultColor, bgColor)
		termbox.SetCell(left+g.Width, i, '│', defaultColor, bgColor)
	}
	termbox.SetCell(left-1, top, '┌', defaultColor, bgColor)
	termbox.SetCell(left-1, bottom, '└', defaultColor, bgColor)
	termbox.SetCell(left+g.Width, top, '┐', defaultColor, bgColor)
	termbox.SetCell(left+g.Width, bottom, '┘', defaultColor, bgColor)
	fill(left, top, g.Width, 1, termbox.Cell{Ch: '─', Fg: defaultColor, Bg: bgColor})
	fill(left, bottom, g.Width, 1, termbox.Cell{Ch: '─', Fg: defaultColor, Bg: bgColor})
}

func (g *Game) renderInfo(left, bottom int) {
	score := fmt.Sprintf("Score: %v Level: %v", g.Score, g.Level)
	tbprint(left, bottom+1, defaultColor, defaultColor, score)
}

func (g *Game) isOnSnake(p Point) bool {
	for _, s := range g.Snake {
		if p.X == s.X && p.Y == s.Y {
			return true
		}
	}
	return false
}

func (g *Game) isOnFood(p Point) bool {
	if p.X == g.Food.X && p.Y == g.Food.Y {
		return true
	}
	return false
}

func (g *Game) isOnMalware(p Point) bool {
	for _, m := range g.Malware {
		if p.X == m.X && p.Y == m.Y {
			return true
		}
	}
	return false
}

func (g *Game) placeFood() {
	for {
		x := rand.Intn(g.Width)
		y := rand.Intn(g.Height)
		food := Point{X: x, Y: y}
		if !g.isOnMalware(food) && !g.isOnSnake(food) && !g.isOnFood(food) {
			g.Food = food
			return
		}
	}
}

func (g *Game) placeMalware() {
	for {
		x := rand.Intn(g.Width)
		y := rand.Intn(g.Height)
		mal := Point{X: x, Y: y}
		if !g.isOnMalware(mal) && !g.isOnSnake(mal) && !g.isOnFood(mal) {
			g.Malware = append(g.Malware, mal)
			return
		}
	}
}

func (g *Game) isOutOfBounds(p Point) bool {
	if p.X < 1 || p.X >= g.Width-1 || p.Y < 1 || p.Y >= g.Height-1 {
		return true
	}
	return false
}

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		return
	}
	defer termbox.Close()
	game := NewGame(60, 25)
	game.draw()
	mainCh := make(chan termbox.Event, 1)
	go func() {
		for {
			mainCh <- termbox.PollEvent()
		}
	}()
	for {
		select {
		case ev := <-mainCh:
			game.handleInput(ev)
			game.draw()
		case <-game.Quit:
			return
		}
	}
}
