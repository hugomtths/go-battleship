package components

import (
	"image/gif"
	"os"
	"time"

	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/hajimehoshi/ebiten/v2"
)

type GIFWidget struct {
	pos, currentPos basic.Point
	frames          []*ebiten.Image
	delays          []int
	currentImg      *ebiten.Image
	totalDuration   int
	size            basic.Size
}

func NewGIFWidget(path string, pos basic.Point, scale float64) (*GIFWidget, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		return nil, err
	}

	frames := make([]*ebiten.Image, len(g.Image))
	delays := make([]int, len(g.Delay))
	totalDuration := 0

	var width, height int

	for i, img := range g.Image {
		frames[i] = ebiten.NewImageFromImage(img)
		delays[i] = g.Delay[i]
		totalDuration += g.Delay[i] * 10
		if i == 0 {
			width = img.Bounds().Dx()
			height = img.Bounds().Dy()
		}
	}

	if totalDuration == 0 {
		totalDuration = 100
	}

	gw := &GIFWidget{
		pos:           pos,
		frames:        frames,
		delays:        delays,
		totalDuration: totalDuration,
		size: basic.Size{
			W: float32(float64(width) * scale),
			H: float32(float64(height) * scale),
		},
	}

	if len(frames) > 0 {
		gw.currentImg = frames[0]
	}

	return gw, nil
}

func (g *GIFWidget) Update(offset basic.Point) {
	g.currentPos = g.pos.Add(offset)

	if len(g.frames) == 0 {
		return
	}

	now := int(time.Now().UnixMilli())
	cycleTime := now % g.totalDuration

	currentDuration := 0
	for k, d := range g.delays {
		frameDuration := d * 10
		if frameDuration == 0 {
			frameDuration = 100
		}
		if cycleTime < currentDuration+frameDuration {
			g.currentImg = g.frames[k]
			break
		}
		currentDuration += frameDuration
	}
}

func (g *GIFWidget) Draw(screen *ebiten.Image) {
	if g.currentImg == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}

	// Calcular escala baseada no tamanho desejado vs tamanho original
	w := float64(g.currentImg.Bounds().Dx())
	h := float64(g.currentImg.Bounds().Dy())

	scaleX := float64(g.size.W) / w
	scaleY := float64(g.size.H) / h

	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(g.currentPos.X), float64(g.currentPos.Y))

	screen.DrawImage(g.currentImg, op)
}

func (g *GIFWidget) GetPos() basic.Point {
	return g.pos
}

func (g *GIFWidget) SetPos(p basic.Point) {
	g.pos = p
}

func (g *GIFWidget) GetSize() basic.Size {
	return g.size
}
