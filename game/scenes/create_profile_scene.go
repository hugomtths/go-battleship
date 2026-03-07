package scenes

import (
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/allanjose001/go-battleship/internal/service"
	"github.com/hajimehoshi/ebiten/v2"
)

type CreateProfileScene struct {
	root       components.LayoutWidget
	nameField  *components.TextField
	errorText  *components.Text
	saveButton *components.Button
	StackHandler
}

func (s *CreateProfileScene) GetMusic() string {
	return "menus"
}

func (s *CreateProfileScene) OnEnter(prev Scene, size basic.Size) {
	label := components.NewText(
		basic.Point{},
		"Digite o seu nome",
		colors.White,
		22,
	)

	fieldSize := basic.Size{W: size.W * 0.6, H: 50}
	s.nameField = components.NewTextField(
		basic.Point{},
		fieldSize,
		"",
	)

	s.errorText = components.NewText(
		basic.Point{},
		"",
		colors.Red,
		18,
	)

	backButton := components.NewButton(
		basic.Point{},
		basic.Size{W: 200, H: 55},
		"voltar",
		colors.Dark,
		nil,
		func(b *components.Button) {
			s.stack.Pop()
		},
	)

	s.saveButton = components.NewButton(
		basic.Point{},
		basic.Size{W: 200, H: 55},
		"salvar",
		colors.Dark,
		nil,
		func(b *components.Button) {
			username := s.nameField.Text

			profile := entity.Profile{
				Username: username,
			}

			err := service.SaveProfile(profile)
			if err != nil {
				s.errorText.Text = "Erro ao salvar perfil"
				return
			}
			s.stack.Pop()
		},
	)

	buttonRow := components.NewContainer(
		basic.Point{},
		basic.Size{W: size.W * 0.55, H: 50},
		0, colors.Transparent,
		basic.Center, basic.Center,
		components.NewRow(
			basic.Point{},
			40, basic.Size{W: size.W * 0.55, H: 50},
			basic.Center,
			basic.Center,
			[]components.Widget{
				backButton,
				s.saveButton,
			},
		),
	)

	s.root = components.NewColumn(
		basic.Point{},
		30,
		size,
		basic.Center,
		basic.Center,
		[]components.Widget{
			label,
			s.nameField,
			components.NewContainer(
				basic.Point{},
				basic.Size{W: size.W * 0.55, H: 40},
				0, colors.Transparent,
				basic.Center, basic.Center,
				s.errorText,
			),
			buttonRow,
		},
	)
	_ = s.Update()
	s.stack.ctx.CanPopOrPush = true

}

func (s *CreateProfileScene) OnExit(next Scene) {
	s.stack.ctx.CanPopOrPush = false

}

func (s *CreateProfileScene) Update() error {
	s.saveButton.SetDisabled(s.nameField.Text == "")
	if s.root != nil {
		s.root.Update(basic.Point{})
	}
	return nil
}

func (s *CreateProfileScene) Draw(screen *ebiten.Image) {
	if s.root != nil {
		s.root.Draw(screen)
	}
}
