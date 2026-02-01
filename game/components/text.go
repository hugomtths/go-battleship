package components

import (
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"image/color"
	"os"
)

// Text widget que desenha texto com uma font.Face
type Text struct {
	Pos   basic.Point
	Color color.Color
	Text  string
	face  font.Face
	Size  basic.Size // opcional, para layout
}

var GoldmanFont *opentype.Font // carregada uma vez

func NewText(pos basic.Point, str string, color color.Color, fontSize int) *Text {
	t := &Text{
		Pos:   pos,
		Text:  str,
		Color: color,
		face:  createFace(float64(fontSize)),
	}

	t.updateSize() // calcula Size ao criar
	return t
}

// InitFonts carrega fonte (*opentype.Font) no inicio do jogo
func InitFonts() {
	data, err := os.ReadFile("assets/fonts/Goldman.ttf")
	if err != nil {
		panic(err)
	}

	GoldmanFont, err = opentype.Parse(data)
	if err != nil {
		panic(err)
	}
	if GoldmanFont == nil {
		panic("GoldmanFont não inicializada. Chame InitFonts() antes de criar Text.")
	}
}

// Face Cria uma Face de um tamanho específico
func createFace(size float64) font.Face {
	face, _ := opentype.NewFace(GoldmanFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	return face
}

// calcula o Size baseado no texto e face usando MeasureString e Metrics
func (t *Text) updateSize() {
	if t.face != nil {

		m := t.face.Metrics()

		ascent := float32(m.Ascent.Round())
		descent := float32(m.Descent.Round())

		t.Size = basic.Size{
			W: float32(font.MeasureString(t.face, t.Text).Round()),
			H: ascent + descent,
		}
	} else {
		t.Size = basic.Size{W: 0, H: 0}
	}
}

// Draw desenha componente com offset calculado
func (t *Text) Draw(screen *ebiten.Image) {
	t.draw(screen, basic.Point{})
}

func (t *Text) draw(screen *ebiten.Image, offset basic.Point) {

	baseline := float32(t.face.Metrics().Ascent.Round())

	//serve para corrigir posição em arvores de widgets
	final := t.Pos.Add(offset)

	if t.Color == nil {
		t.Color = color.White
	}

	text.Draw(screen, t.Text, t.face,
		int(final.X),
		int(baseline+final.Y),
		t.Color,
	)
}

func (t *Text) Update() {
	// criar update caso necessário
}

func (t *Text) GetPos() basic.Point  { return t.Pos }
func (t *Text) SetPos(p basic.Point) { t.Pos = p }
func (t *Text) GetSize() basic.Size  { return t.Size }
func (t *Text) SetSize(_ basic.Size) {}
