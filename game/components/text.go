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

// Text widget que desenha texto com uma font.Face já definida e carregada no incio da run do jogo
type Text struct {
	Pos, currentPos basic.Point
	Color           color.Color
	Text            string
	face            font.Face
	Size            basic.Size // opcional, para layout
	fontSize        int
}

var GoldmanFont *opentype.Font // carregada uma vez

func NewText(
	pos basic.Point,
	str string,
	color color.Color,
	fontSize int,
) *Text {
	t := &Text{
		Pos:      pos,
		Text:     str,
		Color:    color,
		fontSize: fontSize,
		face:     createFace(float64(fontSize)),
	}

	t.updateSize() // calcula Size ao criar
	return t
}

// InitFonts carrega fonte (*opentype.Font) no inicio do jogo
func InitFonts() {
	data, err := os.ReadFile("assets/fonts/PixelifySans-Bold.ttf")
	if err != nil {
		panic(err)
	}

	GoldmanFont, err = opentype.Parse(data)
	if err != nil {
		panic(err)
	}
	if GoldmanFont == nil {
		panic("Pixelify Sans não inicializada. Chame InitFonts() antes de criar Text.")
	}
}

// createFace Cria uma Face de um tamanho específico
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

// Draw desenha componente com offset recebido em Update e somado a pos do Widget
func (t *Text) Draw(screen *ebiten.Image) {

	baseline := float32(t.face.Metrics().Ascent.Round())

	if t.Color == nil {
		t.Color = color.White
	}

	text.Draw(screen, t.Text, t.face,
		int(t.currentPos.X),
		int(baseline+t.currentPos.Y),
		t.Color,
	)
}

func (t *Text) Update(point basic.Point) {
	t.currentPos = t.Pos.Add(point)
}

func (t *Text) GetPos() basic.Point { return t.Pos }

func (t *Text) SetPos(p basic.Point) { t.Pos = p }

func (t *Text) GetSize() basic.Size { return t.Size }

func (t *Text) SetColor(c color.Color) { t.Color = c }

// SetFontSize re-renderiza texto com outro tamanho de fonte
func (t *Text) SetFontSize(px int) {
	if px <= 0 {
		return
	}
	t.fontSize = px
	t.face = createFace(float64(px))
	t.updateSize()
}
