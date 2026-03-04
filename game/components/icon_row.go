package components

import (
	"image/color"

	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/hajimehoshi/ebiten/v2"
)

type IconRow struct {
	pos, currentPos basic.Point
	body            Widget
}

func NewIconRow(path, label, data string, size basic.Size, pos basic.Point, statusColor color.Color) (Widget, error) {

	icon, err := NewImage(path, basic.Point{}, basic.Size{
		W: 40,
		H: 40,
	})

	return &IconRow{
		pos: pos,
		body: NewContainer(
			basic.Point{}, size, 0,
			colors.Transparent, basic.Center, basic.Center,
			NewRow(
				basic.Point{}, 10, size,
				basic.Start, basic.Center,
				[]Widget{
					icon,
					NewText(basic.Point{}, label, colors.White, 35),
					NewText(basic.Point{}, data, statusColor, 35),
				},
			),
		),
	}, err

}

func (i *IconRow) GetPos() basic.Point {
	return i.pos
}

func (i *IconRow) SetPos(point basic.Point) {
	i.pos = point
}

func (i *IconRow) GetSize() basic.Size {
	return i.body.GetSize()
}

func (i *IconRow) Update(offset basic.Point) {
	i.currentPos = i.pos.Add(offset)
	i.body.Update(i.currentPos)
}

func (i *IconRow) Draw(screen *ebiten.Image) {
	i.body.Draw(screen)
}
