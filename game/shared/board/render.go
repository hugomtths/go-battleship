package board

//responsabilidade de desenhar o tabuleiro.

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

func (b *Board) Draw(screen *ebiten.Image) {
	cellSize := b.Size / Cols

	// Draw Background
	if b.BackgroundImage != nil {
		op := &ebiten.DrawImageOptions{}
		imgW, imgH := b.BackgroundImage.Size()
		op.GeoM.Scale(b.Size/float64(imgW), b.Size/float64(imgH))
		op.GeoM.Translate(b.X, b.Y)
		screen.DrawImage(b.BackgroundImage, op)
	} else {
		// Fallback background
		ebitenutil.DrawRect(screen, b.X, b.Y, b.Size, b.Size, color.RGBA{48, 67, 103, 255})
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
	labelSize := float64(cellSize) * 0.5
	if labelSize < 12 {
		labelSize = 12
	}
	tt, _ := opentype.Parse(goregular.TTF)
	face, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    labelSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	labelColor := color.White

	// topo: letras A-H
	for j := 0; j < Cols; j++ {
		ch := string(rune('A' + j))
		bounds := text.BoundString(face, ch)
		w := bounds.Dx()
		x := b.X + float64(j)*cellSize + cellSize/2 - float64(w)/2
		y := b.Y - 10
		text.Draw(screen, ch, face, int(x), int(y), labelColor)
	}

	// esquerda: números 1-7
	for i := 0; i < Rows; i++ {
		num := strconv.Itoa(i + 1)
		bounds := text.BoundString(face, num)
		h := bounds.Dy()
		x := b.X - 25
		y := b.Y + float64(i)*cellSize + cellSize/2 + float64(h)/2 - 2
		text.Draw(screen, num, face, int(x), int(y), labelColor)
	}
}
