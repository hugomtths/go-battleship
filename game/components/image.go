package components

import (
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Image struct {
	image           *ebiten.Image
	pos, currentPos basic.Point
	size            basic.Size
}

func NewImage(path string, pos basic.Point, size basic.Size) (*Image, error) {
	img, err := loadImage(path)

	if err != nil {
		return nil, err
	}
	return &Image{
		image: img,
		pos:   pos,
		size:  size,
	}, nil
}

func loadImage(path string) (*ebiten.Image, error) {

	image, _, err := ebitenutil.NewImageFromFile(path)

	return image, err
}

func (i *Image) GetPos() basic.Point {
	return i.pos
}

func (i *Image) SetPos(point basic.Point) {
	i.pos = point
}

// GetSize retorna tamanho da imagem 1:1 ou tamanho settado
func (i *Image) GetSize() basic.Size {
	if i.size.W == 0 || i.size.H == 0 {
		return basic.Size{
			W: float32(i.image.Bounds().Dx()),
			H: float32(i.image.Bounds().Dy()),
		}
	}
	return i.size
}

func (i *Image) Update(offset basic.Point) {
	i.currentPos = i.pos.Add(offset)
}

// Draw renderiza a imagem com o tamanho especificado no construtor
func (i *Image) Draw(screen *ebiten.Image) {
	if i.image == nil {
		return
	}

	w := float32(i.image.Bounds().Dx())
	h := float32(i.image.Bounds().Dy())

	size := i.GetSize()
	scaleX := size.W / w
	scaleY := size.H / h

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(float64(scaleX), float64(scaleY))
	opts.GeoM.Translate(float64(i.currentPos.X), float64(i.currentPos.Y))

	screen.DrawImage(i.image, opts)
}
