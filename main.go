package main

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

func main() {

}
