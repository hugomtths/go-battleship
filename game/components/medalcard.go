package components

import (
	"image/color"

	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/hajimehoshi/ebiten/v2"
)

// MedalCard encapsula um card visual de medalha reutilizável.
// Ele funciona como um widget composto: toda a renderização real é delegada
// para o campo `body`, que contém a árvore de containers/layouts internos.
//
// TODO: Para medalhas não adquiridas, a ideia é permitir estilização externa
// (ex: cor diferente) via StylableWidget.
type MedalCard struct {
	pos, currentPos basic.Point    // posição base e posição final com offset aplicado
	size            basic.Size     // tamanho total do card (usado pelo layout pai)
	body            StylableWidget // widget raiz que contém todo o layout interno
}

// NewMedal constrói a estrutura completa do card.
// O layout interno é proporcional ao tamanho recebido:
//
//	ATUALMENTE:  10% da largura -> área do ícone, 50% da largura -> área de texto
//
// O restante funciona como padding visual.
func NewMedal(icon, title, desc string, size basic.Size) *MedalCard {
	iconTxt := NewText(basic.Point{}, icon, color.RGBA{255, 200, 0, 255}, 28) // ícone da medalha
	titleTxt := NewText(basic.Point{}, title, color.RGBA{40, 40, 50, 255}, 16)
	descTxt := NewText(basic.Point{}, desc, color.RGBA{100, 100, 110, 255}, 12)

	return &MedalCard{
		size: size,
		body: NewContainer( // container pai do card inteiro
			basic.Point{},
			size,
			15, // borda mais arredondada para destacar visualmente
			colors.White,
			basic.Center,
			basic.Center,
			NewRow(
				basic.Point{}, 10,
				size,
				basic.Center,
				basic.Center,
				[]Widget{
					// Container do ícone (largura proporcional ao pai)
					NewContainer(
						basic.Point{},
						basic.Size{
							W: size.W * 0.1,
							H: size.H,
						},
						0, colors.Transparent, basic.Center, basic.Center,
						iconTxt,
					),

					// Container de texto (título + descrição)
					NewContainer(
						basic.Point{},
						basic.Size{
							W: size.W * 0.5,
							H: size.H,
						},
						0, colors.Transparent, basic.Start, basic.Start,
						NewColumn(
							basic.Point{}, 12,
							basic.Size{
								W: size.W * 0.6,
								H: size.H,
							},
							basic.Center,
							basic.Center,
							[]Widget{
								titleTxt,
								descTxt,
							},
						),
					),
				},
			),
		),
	}
}

// GetPos retorna a posição base do card.
func (m *MedalCard) GetPos() basic.Point {
	return m.pos
}

// SetPos define a posição base do card.
func (m *MedalCard) SetPos(point basic.Point) {
	m.pos = point
}

// GetSize retorna o tamanho total do card.
func (m *MedalCard) GetSize() basic.Size {
	return m.size
}

// Update aplica offset acumulado e propaga para o layout interno.
func (m *MedalCard) Update(offset basic.Point) {
	m.currentPos = m.pos.Add(offset)
	m.body.Update(m.currentPos)
}

// Draw delega o desenho para o widget raiz.
func (m *MedalCard) Draw(screen *ebiten.Image) {
	m.body.Draw(screen)
}
