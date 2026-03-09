package scenes

import (
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/allanjose001/go-battleship/internal/service"
	"github.com/hajimehoshi/ebiten/v2"
)

type SelectProfileScene struct {
	root             components.LayoutWidget
	profiles         []entity.Profile
	screenSize       basic.Size
	newProfileButton *components.Button
	StackHandler
}

func (s *SelectProfileScene) GetMusic() string {
	return "menus"
}

func (s *SelectProfileScene) OnEnter(prev Scene, size basic.Size) {
	s.profiles = service.GetProfiles()
	s.screenSize = size
	s.buildUI(size)
	s.stack.ctx.CanPopOrPush = true

}

func (s *SelectProfileScene) OnExit(next Scene) {
	s.stack.ctx.CanPopOrPush = false
}

func (s *SelectProfileScene) Update() error {
	s.newProfileButton.SetDisabled(len(s.profiles) == 5) //desabilita caso n maximo de perfis salvos
	s.root.Update(basic.Point{})
	return nil
}

func (s *SelectProfileScene) Draw(screen *ebiten.Image) {
	s.root.Draw(screen)

}

func (s *SelectProfileScene) buildUI(size basic.Size) {
	title := components.NewText(
		basic.Point{},
		"Jogadores",
		colors.White,
		35,
	)

	spacer := components.NewContainer(
		basic.Point{},
		basic.Size{W: size.W, H: 20},
		0,
		colors.Transparent,
		basic.Center,
		basic.Center,
		nil,
	)

	listScreenSize := basic.Size{W: size.W * 0.7, H: 440}

	// Lista de perfis
	profileList := s.buildProfileRows(size.W*0.7, listScreenSize)

	profilesWrappler := components.NewContainer(
		basic.Point{},
		listScreenSize,
		0,
		colors.Transparent,
		basic.Center,
		basic.Center,
		profileList,
	)

	backButton := components.NewButton(
		basic.Point{},
		basic.Size{W: 220, H: 55},
		"Voltar",
		colors.Dark,
		nil,
		func(b *components.Button) {
			s.ctx.SoundService.PlaySFX("backclick", 0.8)
			s.stack.Pop()
		},
	)

	s.newProfileButton = components.NewButton(
		basic.Point{},
		basic.Size{W: 220, H: 55},
		"Novo Jogador",
		colors.Dark,
		nil,
		func(b *components.Button) {
			s.ctx.SoundService.PlaySFX("click", 0.8)
			s.stack.Push(&CreateProfileScene{})
		},
	)

	buttonRowWrappler := components.NewContainer(
		basic.Point{},
		basic.Size{W: size.W, H: 80},
		0,
		colors.Transparent,
		basic.Center,
		basic.Center,
		components.NewRow(
			basic.Point{},
			40,
			basic.Size{W: size.W, H: 80},
			basic.Center,
			basic.Center,
			[]components.Widget{backButton, s.newProfileButton},
		),
	)

	s.root = components.NewColumn(
		basic.Point{X: 0, Y: 0},
		20,
		size,
		basic.Start,
		basic.Center,
		[]components.Widget{spacer, title, spacer, profilesWrappler, spacer, buttonRowWrappler},
	)

	_ = s.Update()
}

// buildProfileRows cria as linhas de perfil
func (s *SelectProfileScene) buildProfileRows(width float32, parentSize basic.Size) components.Widget {
	rows := make([]components.Widget, len(s.profiles))
	iconSize := basic.Size{W: 45, H: 45}

	for i, p := range s.profiles {
		rows[i] = s.createProfileRow(&p, width, iconSize)
	}

	return components.NewColumn(
		basic.Point{X: 0, Y: 0},
		45,
		parentSize,
		basic.Center,
		basic.Center,
		rows,
	)
}

// createProfileRow cria uma única linha de perfil
func (s *SelectProfileScene) createProfileRow(p *entity.Profile, width float32, iconSize basic.Size) components.Widget {
	// Captura o perfil para o closure
	profile := p

	// Botão com o nome do jogador
	nameBtn := components.NewButton(
		basic.Point{},
		basic.Size{W: width * 0.5, H: 50},
		profile.Username,
		colors.PlayerInput,
		nil,
		func(b *components.Button) {
			s.ctx.SoundService.PlaySFX("click", 0.8)
			s.stack.ctx.Profile = p
			s.stack.Push(&ProfileScene{})

		},
	)

	// Ícone de deletar
	deleteBtn := components.NewDeleteIconButton(basic.Point{}, iconSize, func() {
		_ = service.RemoveProfile(profile.Username)
		s.profiles = service.GetProfiles()
		s.ctx.SoundService.PlaySFX("click", 0.8)
		s.buildUI(s.screenSize)
	})

	// Ícone de jogar
	playBtn := components.NewPlayIconButton(basic.Point{}, iconSize, func() {
		s.stack.ctx.Profile = p
		s.ctx.SoundService.PlaySFX("click", 0.8)

		s.stack.Push(&ModeSelectionScene{})
	})

	// Monta a linha: [delete] [nome] [play]
	rowWidgets := []components.Widget{}
	if deleteBtn != nil {
		rowWidgets = append(rowWidgets, deleteBtn)
	}
	rowWidgets = append(rowWidgets, nameBtn)
	if playBtn != nil {
		rowWidgets = append(rowWidgets, playBtn)
	}

	return components.NewContainer(
		basic.Point{},
		basic.Size{W: width * 0.5, H: 50},
		0,
		colors.Transparent,
		basic.Center,
		basic.Center,
		components.NewRow(
			basic.Point{},
			10, // espaçamento entre elementos
			basic.Size{W: width * 0.5, H: 50},
			basic.Center,
			basic.Center,
			rowWidgets,
		),
	)
}
