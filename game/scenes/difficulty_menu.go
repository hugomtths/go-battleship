package scenes

import (
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

type DifficultyScene struct {
	layout components.LayoutWidget
	StackHandler
}

func (d *DifficultyScene) GetMusic() string {
	return "menus"
}

func (d *DifficultyScene) init(size basic.Size) {

	btnSize := basic.Size{W: 450, H: 50}

	btnRecruta := components.NewButton(
		basic.Point{},
		btnSize,
		"Recruta",
		colors.Blue,
		colors.White,
		func(b *components.Button) {
			d.ctx.SoundService.PlaySFX("click", 0.8)
			d.selectDifficulty("easy")
		},
	)

	btnImediato := components.NewButton(
		basic.Point{},
		btnSize,
		"Imediato",
		colors.Blue,
		colors.White,
		func(b *components.Button) {
			d.ctx.SoundService.PlaySFX("click", 0.8)
			d.selectDifficulty("medium")
		},
	)

	btnAlmirante := components.NewButton(
		basic.Point{},
		btnSize,
		"Almirante",
		colors.Blue,
		colors.White,
		func(b *components.Button) {
			d.ctx.SoundService.PlaySFX("click", 0.8)
			d.selectDifficulty("hard")
		},
	)

	btnVoltar := components.NewButton(
		basic.Point{},
		basic.Size{W: 220, H: 50},
		"Voltar",
		colors.Dark,
		colors.White,
		func(b *components.Button) {
			d.ctx.SoundService.PlaySFX("backclick", 0.8)
			d.stack.Pop()
		},
	)

	screenSize := basic.Size{W: size.W, H: size.H}

	spacer := components.NewContainer(
		basic.Point{}, basic.Size{W: 1, H: 20}, 0,
		colors.Transparent, basic.Center, basic.Center,
		nil,
	)
	spacer2 := components.NewContainer(
		basic.Point{}, basic.Size{W: 1, H: 100}, 0,
		colors.Transparent, basic.Center, basic.Center,
		nil,
	)

	d.layout = components.NewColumn(
		basic.Point{X: 0, Y: 0},
		20,
		screenSize,
		basic.Start,
		basic.Center,
		[]components.Widget{
			spacer,
			components.NewText(basic.Point{}, "Seleção de Dificuldade", colors.White, 35),
			spacer2,
			btnRecruta,
			spacer,
			btnImediato,
			spacer,
			btnAlmirante,
			spacer2,
			btnVoltar,
		},
	)
	_ = d.Update()
}

func (d *DifficultyScene) selectDifficulty(diff string) {

	var prof *entity.Profile

	if d.ctx != nil {
		d.ctx.SetDifficulty(diff)
		d.ctx.IsCampaign = false
		prof = d.ctx.Profile
	}

	d.stack.Push(NewPlacementSceneWithProfile(prof))
}

func (d *DifficultyScene) OnEnter(prev Scene, size basic.Size) {
	d.init(size)
	d.stack.ctx.CanPopOrPush = true
}

func (d *DifficultyScene) OnExit(next Scene) {
	d.stack.ctx.CanPopOrPush = false
}

func (d *DifficultyScene) Update() error {
	if d.layout != nil {
		d.layout.Update(basic.Point{X: 0, Y: 0})
	}
	return nil
}

func (d *DifficultyScene) Draw(screen *ebiten.Image) {
	if d.layout != nil {
		d.layout.Draw(screen)
	}
}
