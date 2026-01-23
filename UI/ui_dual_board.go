package ui

import (
	"image/color"
	"image/gif"
	"math"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

const (
	gap = 100
)

type DualBoardUI struct {
	Rows, Cols                 int
	CellSize                   int
	LeftOffsetX, LeftOffsetY   int
	RightOffsetX, RightOffsetY int
	LabelFace                  font.Face
	LabelSize                  int
	TitleFace                  font.Face
	TitleSize                  int
	bgFrames                   []*ebiten.Image
	bgDelays                   []time.Duration
	bgIndex                    int
	bgReady                    bool
	bgElapsed                  time.Duration
	lastUpdate                 time.Time
}

func NewDualBoardUI(rows, cols int) *DualBoardUI {
	return &DualBoardUI{Rows: rows, Cols: cols}
}

func (b *DualBoardUI) Update() error {
	if b.lastUpdate.IsZero() {
		b.lastUpdate = time.Now()
	}
	now := time.Now()
	dt := now.Sub(b.lastUpdate)
	b.lastUpdate = now

	if b.bgReady && len(b.bgFrames) > 0 {
		b.bgElapsed += dt
		for b.bgElapsed >= b.bgDelays[b.bgIndex] {
			b.bgElapsed -= b.bgDelays[b.bgIndex]
			b.bgIndex = (b.bgIndex + 1) % len(b.bgFrames)
		}
	}
	return nil
}

func (b *DualBoardUI) Draw(screen *ebiten.Image) {
	b.ensureBackgroundGIF("assets/Temas Para Festa Na Piscina.gif")
	if b.bgReady {
		drawGIFBackground(screen, b)
	}

	availW := float64(screenWidth - 2*margin - gap)
	availH := float64(screenHeight - 2*margin)
	perBoardW := availW / 2
	cellW := perBoardW / float64(b.Cols)
	cellH := availH / float64(b.Rows)
	cellSize := int(math.Floor(math.Min(cellW, cellH)))
	cellSize = int(float64(cellSize) * 0.9)
	if cellSize < 8 {
		cellSize = 8
	}
	b.CellSize = cellSize

	boardW := cellSize * b.Cols
	boardH := cellSize * b.Rows

	baseX := (screenWidth - (2*boardW + gap)) / 2
	offsetY := (screenHeight - boardH) / 2
	shiftX := -8
	shiftY := -12
	baseX += shiftX
	offsetY += shiftY

	b.LeftOffsetX = baseX
	b.LeftOffsetY = offsetY
	b.RightOffsetX = baseX + boardW + gap
	b.RightOffsetY = offsetY

	newLabelSize := int(math.Max(12, math.Min(20, float64(cellSize)*0.4)))
	if b.LabelFace == nil || b.LabelSize != newLabelSize {
		tt, _ := opentype.Parse(goregular.TTF)
		face, _ := opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    float64(newLabelSize),
			DPI:     72,
			Hinting: font.HintingFull,
		})
		b.LabelFace = face
		b.LabelSize = newLabelSize
	}
	newTitleSize := int(math.Max(14, math.Min(22, float64(cellSize)*0.45)))
	if b.TitleFace == nil || b.TitleSize != newTitleSize {
		tt, _ := opentype.Parse(goregular.TTF)
		face, _ := opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    float64(newTitleSize),
			DPI:     72,
			Hinting: font.HintingFull,
		})
		b.TitleFace = face
		b.TitleSize = newTitleSize
	}

	drawGrid(screen, b.LeftOffsetX, b.LeftOffsetY, b.Rows, b.Cols, cellSize)
	drawAxisLabels(screen, b.LeftOffsetX, b.LeftOffsetY, b.Rows, b.Cols, cellSize, b.LabelFace)
	drawGrid(screen, b.RightOffsetX, b.RightOffsetY, b.Rows, b.Cols, cellSize)
	drawAxisLabels(screen, b.RightOffsetX, b.RightOffsetY, b.Rows, b.Cols, cellSize, b.LabelFace)

	labelColor := color.Black
	text.Draw(screen, "Jogador", b.TitleFace, b.LeftOffsetX, b.LeftOffsetY-40, labelColor)
	text.Draw(screen, "Computador", b.TitleFace, b.RightOffsetX, b.RightOffsetY-40, labelColor)
}

func (b *DualBoardUI) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func drawGrid(screen *ebiten.Image, ox, oy, rows, cols, cellSize int) {
	lineColor := color.RGBA{0x00, 0x00, 0x00, 0x88}

	boardW := cellSize * cols
	boardH := cellSize * rows

	for i := 0; i <= cols; i++ {
		x := float64(ox + i*cellSize)
		ebitenutil.DrawLine(screen, x, float64(oy), x, float64(oy+boardH), lineColor)
	}
	for i := 0; i <= rows; i++ {
		y := float64(oy + i*cellSize)
		ebitenutil.DrawLine(screen, float64(ox), y, float64(ox+boardW), y, lineColor)
	}
}

func drawAxisLabels(screen *ebiten.Image, ox, oy, rows, cols, cellSize int, face font.Face) {
	colColor := color.Black
	rowColor := color.Black
	for c := 0; c < cols && c < 26; c++ {
		ch := string(rune('A' + c))
		bounds := text.BoundString(face, ch)
		w := bounds.Max.X - bounds.Min.X
		x := ox + c*cellSize + (cellSize-w)/2
		y := oy - 8
		text.Draw(screen, ch, face, x, y, colColor)
	}
	for r := 0; r < rows; r++ {
		lbl := ""
		if r+1 == 10 {
			lbl = "10"
		} else {
			lbl = string(rune('0' + (r + 1)))
		}
		bounds := text.BoundString(face, lbl)
		w := bounds.Max.X - bounds.Min.X
		h := bounds.Max.Y - bounds.Min.Y
		x := ox - w - 8
		y := oy + r*cellSize + (cellSize+h)/2 - 2
		text.Draw(screen, lbl, face, x, y, rowColor)
	}
}

func (b *DualBoardUI) ensureBackgroundGIF(path string) {
	if b.bgReady || len(b.bgFrames) > 0 {
		return
	}
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	g, err := gif.DecodeAll(f)
	if err != nil {
		return
	}
	for i := range g.Image {
		img := ebiten.NewImageFromImage(g.Image[i])
		b.bgFrames = append(b.bgFrames, img)
		delay := time.Duration(g.Delay[i]*10) * time.Millisecond
		if delay <= 0 {
			delay = 80 * time.Millisecond
		}
		b.bgDelays = append(b.bgDelays, delay)
	}
	if len(b.bgFrames) > 0 {
		b.bgReady = true
	}
}

func drawGIFBackground(screen *ebiten.Image, b *DualBoardUI) {
	if len(b.bgFrames) == 0 {
		return
	}
	frame := b.bgFrames[b.bgIndex]
	w, h := frame.Size()
	sx := float64(screenWidth) / float64(w)
	sy := float64(screenHeight) / float64(h)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(sx, sy)
	screen.DrawImage(frame, op)
}
