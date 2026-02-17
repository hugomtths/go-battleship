package placement

import (
	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/hajimehoshi/ebiten/v2"
)

type ShipPlacement struct {
	Image        *ebiten.Image
	Size         int
	Placed       bool
	X, Y         int
	Orientation  board.Orientation
	ListX, ListY float64

	Dragging         bool
	DragX, DragY     float64
	OffsetX, OffsetY float64
}
