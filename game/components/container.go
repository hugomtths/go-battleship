package components

import (
	"image/color"

	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Container retangulo que comporta um widget (tambem pode ser column ou row)
type Container struct {
	Pos    basic.Point
	Size   basic.Size
	Radius float32 //opcional
	Color  color.Color
	Child  Widget

	MainAlign  basic.Align
	CrossAlign basic.Align

	OnUpdate func(c *Container) //para injetar algum comportamento
}

func NewContainer(pos basic.Point, size basic.Size,
	radius float32, color color.Color, mainAlign basic.Align, crossAlign basic.Align,
	child Widget, callBack func(c *Container)) *Container {
	c := &Container{
		Pos:        pos,
		Size:       size,
		Radius:     radius,
		Color:      color,
		Child:      child,
		MainAlign:  mainAlign,
		CrossAlign: crossAlign,
		OnUpdate:   callBack,
	}

	c.alignChild() // aplica alinhamento no filho

	return c

}

func (c *Container) Update() {
	//fazer logica caso queira alterar o container na tela
	if c.OnUpdate != nil {
		c.OnUpdate(c)
	}
	if c.Child != nil {
		c.Child.Update()
	}
}

// Draw desenha o container e seu filho (opcional)
func (c *Container) Draw(screen *ebiten.Image) {
	c.draw(screen, basic.Point{})
}

func (c *Container) draw(screen *ebiten.Image, offset basic.Point) {
	final := c.Pos.Add(offset)
	drawRoundedRect(screen, final, c.Size, c.Radius, c.Color)
	if c.Child != nil {
		c.Child.draw(screen, final)
	}
}

func (c *Container) GetPos() basic.Point {
	return c.Pos
}

func (c *Container) SetPos(p basic.Point) {
	c.Pos = p
}

func (c *Container) GetSize() basic.Size {
	return c.Size
}

func (c *Container) SetSize(s basic.Size) {
	c.Size = s
}

// alinhamento de child dentro do container
func (c *Container) alignChild() {
	if c.Child == nil {
		return
	}

	childSize := c.Child.GetSize()
	pos := basic.Point{}

	// horizontal
	switch c.MainAlign {
	case basic.Start:
		pos.X = 0
	case basic.Center:
		pos.X = (c.Size.W - childSize.W) / 2
	case basic.End:
		pos.X = c.Size.W - childSize.W
	}

	// vertical
	switch c.CrossAlign {
	case basic.Start:
		pos.Y = 0
	case basic.Center:
		pos.Y = (c.Size.H - childSize.H) / 2
	case basic.End:
		pos.Y = c.Size.H - childSize.H
	}

	c.Child.SetPos(pos)
}

// NewButton Construtor.
// w e h podem ser 0 para usar tamanho padrão.
// OnClick e OnHover podem ser nil.

// drawRoundedRect desenha um retangulo com borda arredondada
func drawRoundedRect(dst *ebiten.Image, pos basic.Point, size basic.Size, r float32, c color.Color) {
	if c == nil { //nesse caso serve apenas para referencia de tamanho (semelhante a sizedbox)
		c = colors.Transparent
	}

	//quando não usar borda arredondada, apenas desenha rect normal
	if r == 0.0 {
		img := ebiten.NewImage(int(size.W), int(size.H))
		img.Fill(c)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(pos.X), float64(pos.Y))

		dst.DrawImage(img, op)
	}

	w := size.W
	h := size.H
	x := pos.X
	y := pos.Y

	var p vector.Path

	p.MoveTo(x+r, y)
	p.LineTo(x+w-r, y)
	p.QuadTo(x+w, y, x+w, y+r)

	p.LineTo(x+w, y+h-r)
	p.QuadTo(x+w, y+h, x+w-r, y+h)

	p.LineTo(x+r, y+h)
	p.QuadTo(x, y+h, x, y+h-r)

	p.LineTo(x, y+r)
	p.QuadTo(x, y, x+r, y)
	p.Close()

	cr, cg, cb, ca := c.RGBA()

	drawOpts := &vector.DrawPathOptions{
		AntiAlias: true,
	}
	drawOpts.ColorScale.Scale(
		float32(cr)/0xffff,
		float32(cg)/0xffff,
		float32(cb)/0xffff,
		float32(ca)/0xffff,
	)

	fillOpts := &vector.FillOptions{
		FillRule: vector.FillRuleNonZero,
	}

	vector.FillPath(dst, &p, fillOpts, drawOpts)
}
