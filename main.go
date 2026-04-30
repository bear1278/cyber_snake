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

func NewGame(snake []Point, malware []Point, food Point, dir Point, score int, level int, gameOver bool, width int, height int, quit chan struct{}) *Game {
	return &Game{Snake: snake, Malware: malware, Food: food, Dir: dir, Score: score, Level: level, GameOver: gameOver, Width: width, Height: height, Quit: quit}
}

func main() {

}
