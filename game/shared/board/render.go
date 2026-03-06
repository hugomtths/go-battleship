package board

//responsabilidade de desenhar o tabuleiro.

import (
	"image/color"
	"strconv"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

// cache da face para evitar recriar a cada frame
var (
	boardFace     font.Face
	boardFaceOnce sync.Once
)

func getBoardFace(labelSize float64) font.Face {
	boardFaceOnce.Do(func() {
		tt, _ := opentype.Parse(goregular.TTF)
		boardFace, _ = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    labelSize,
			DPI:     72,
			Hinting: font.HintingFull,
		})
	})
	return boardFace
}

func (b *Board) Draw(screen *ebiten.Image) {
	cellSize := b.Size / Cols

	if b.BackgroundImage != nil {
		op := &ebiten.DrawImageOptions{}
		imgW, imgH := b.BackgroundImage.Size()
		op.GeoM.Scale(b.Size/float64(imgW), b.Size/float64(imgH))
		op.GeoM.Translate(b.X, b.Y)
		screen.DrawImage(b.BackgroundImage, op)
	} else {
		// Fallback background
		ebitenutil.DrawRect(screen, b.X, b.Y, b.Size, b.Size, color.RGBA{20, 30, 60, 255})
	}

	// Draw Grid Lines (White)
	gridColor := color.White
	for i := 0; i <= Rows; i++ {
		y := b.Y + float64(i)*cellSize
		ebitenutil.DrawLine(screen, b.X, y, b.X+b.Size, y, gridColor)
	}
	for j := 0; j <= Cols; j++ {
		x := b.X + float64(j)*cellSize
		ebitenutil.DrawLine(screen, x, b.Y, x, b.Y+b.Size, gridColor)
	}

	// labels
	labelSize := cellSize * 0.5
	if labelSize < 12 {
		labelSize = 12
	}

	// reutiliza face cacheada — zero alocações por frame
	face := getBoardFace(labelSize)
	labelColor := color.White

	// topo: letras A-H
	for j := 0; j < Cols; j++ {
		ch := string(rune('A' + j))
		x := int(b.X + float64(j)*cellSize + cellSize*0.3)
		y := int(b.Y - 5)
		text.Draw(screen, ch, face, x, y, labelColor)
	}

	// esquerda: números 1-7
	for i := 0; i < Rows; i++ {
		num := strconv.Itoa(i + 1)
		x := int(b.X - cellSize*0.4)
		y := int(b.Y + float64(i)*cellSize + cellSize*0.7)
		text.Draw(screen, num, face, x, y, labelColor)
	}
}
