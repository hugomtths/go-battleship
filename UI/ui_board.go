package ui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	margin       = 20
)

type BoardUI struct {
	Rows, Cols int
	CellSize   int
	OffsetX    int
	OffsetY    int
}

func NewBoardUI(rows, cols int) *BoardUI { return &BoardUI{Rows: rows, Cols: cols} }

func (b *BoardUI) Update() error { return nil }

func (b *BoardUI) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x1f, 0x2b, 0x3a, 0xff})

	availW := screenWidth - 2*margin
	availH := screenHeight - 2*margin
	cellW := float64(availW) / float64(b.Cols)
	cellH := float64(availH) / float64(b.Rows)
	cellSize := int(math.Floor(math.Min(cellW, cellH)))
	b.CellSize = cellSize

	boardW := cellSize * b.Cols
	boardH := cellSize * b.Rows
	b.OffsetX = (screenWidth - boardW) / 2
	b.OffsetY = (screenHeight - boardH) / 2

	cellFill := color.RGBA{0xe0, 0xf0, 0xff, 0xff}
	lineColor := color.Black

	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Cols; c++ {
			x := float64(b.OffsetX + c*cellSize)
			y := float64(b.OffsetY + r*cellSize)
			ebitenutil.DrawRect(screen, x, y, float64(cellSize), float64(cellSize), cellFill)
		}
	}

	for i := 0; i <= b.Cols; i++ {
		x := float64(b.OffsetX + i*cellSize)
		ebitenutil.DrawLine(screen, x, float64(b.OffsetY), x, float64(b.OffsetY+boardH), lineColor)
	}
	for i := 0; i <= b.Rows; i++ {
		y := float64(b.OffsetY + i*cellSize)
		ebitenutil.DrawLine(screen, float64(b.OffsetX), y, float64(b.OffsetX+boardW), y, lineColor)
	}
}

func (b *BoardUI) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
