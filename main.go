package main

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const (
	defaultColor = termbox.ColorDefault
	bgColor      = termbox.ColorDefault
	snakeColor   = termbox.ColorGreen
)

type Point struct {
	X, Y int
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
	return &Game{
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
}

func (g *Game) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w, h := termbox.Size()
	midY := h / 2
	left := (w - g.Width) / 2
	top := midY - (g.Height / 2)
	bottom := midY + (g.Height / 2) + 1
	g.renderArea(top, bottom, left)
	g.renderSnake(left, bottom)
	g.renderInfo(left, bottom)
	termbox.Flush()
}

func (g *Game) renderSnake(left, bottom int) {
	for _, b := range g.Snake {
		termbox.SetCell(left+b.X, bottom-b.Y, ' ', snakeColor, snakeColor)
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
		fmt.Println(err.Error())
		return
	}
	defer termbox.Close()
	game := NewGame(60, 25)
	game.draw()
}
