package ui

import "github.com/hajimehoshi/ebiten/v2"

type Screen interface {
    Update() error
    Draw(screen *ebiten.Image)
    Layout(outsideWidth, outsideHeight int) (int, int)
}