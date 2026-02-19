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
	root      components.Widget
	nameField *components.TextField
	errorText *components.Text
	StackHandler
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
			if SwitchTo != nil {
				SwitchTo(&SelectProfileScene{})
			}
		},
	)

	saveButton := components.NewButton(
		basic.Point{},
		basic.Size{W: 200, H: 55},
		"salvar",
		colors.Dark,
		nil,
		func(b *components.Button) {
			username := s.nameField.Text
			if username == "" {
				s.errorText.Text = "Nome não pode ser vazio"
				return
			}

			profile := entity.Profile{
				Username: username,
			}

			err := service.SaveProfile(profile)
			if err != nil {
				s.errorText.Text = "Erro ao salvar perfil"
				return
			}

			if SwitchTo != nil {
				SwitchTo(&SelectProfileScene{})
			}
		},
	)

	buttonRow := components.NewRow(
		basic.Point{},
		40,
		basic.Size{W: size.W, H: 80},
		basic.Center,
		basic.Center,
		[]components.Widget{
			backButton,
			saveButton,
		},
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
			s.errorText,
			buttonRow,
		},
	)
}

func (s *CreateProfileScene) OnExit(next Scene) {}

func (s *CreateProfileScene) Update() error {
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
