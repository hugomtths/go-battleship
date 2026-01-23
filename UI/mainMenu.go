package ui

import (

	"github.com/hajimehoshi/ebiten/v2"
)

type App struct {
	current Screen
}

func NewApp() *App {
	// troque rows/cols dinamicamente conforme quiser
	return &App{current: NewDualBoardUI(10, 10)}
}

func (a *App) Update() error              { return a.current.Update() }
func (a *App) Draw(screen *ebiten.Image)  { a.current.Draw(screen) }
func (a *App) Layout(w, h int) (int, int) { return a.current.Layout(w, h) }

func Run() error {
    ebiten.SetWindowSize(screenWidth, screenHeight)
    ebiten.SetWindowTitle("Battleship - Preview")
    return ebiten.RunGame(NewApp())
}
