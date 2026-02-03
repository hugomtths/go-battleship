package scenes

import (
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/hajimehoshi/ebiten/v2"
)

type ButtonScene struct {
	button1 *components.Button
	button2 *components.Button
	row     components.Row
}

func (b *ButtonScene) OnEnter(prev Scene) {
	b.button1 = components.NewButton(
		basic.Point{},
		basic.Size{W: 200.0, H: 70.0},
		"jogar",
		colors.Blue,
		nil,
		func(bt *components.Button) {
		},
	)
	b.button2 = components.NewButton(
		basic.Point{},
		basic.Size{W: 200.0, H: 70.0},
		"voltar",
		colors.Dark,
		nil,
		func(bt *components.Button) {
		},
	)
	b.row = *components.NewRow(
		basic.Point{},
		100,
		basic.Size{
			W: 1280,
			H: 720,
		},
		basic.Center,
		basic.Center,
		[]components.Widget{
			b.button1,
			b.button2,
		},
	)
}

func (b *ButtonScene) OnExit(next Scene) {

}

func (b *ButtonScene) Update() error {
	b.row.Update()
	return nil
}

func (b *ButtonScene) Draw(screen *ebiten.Image) {
	b.row.Draw(screen)
}
