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
	root             components.Widget
	profiles         []entity.Profile
	screenSize       basic.Size
	newProfileButton *components.Button
	builded          bool
	StackHandler
}

func (s *SelectProfileScene) OnEnter(prev Scene, size basic.Size) {
	s.profiles = service.GetProfiles()
	s.screenSize = size
	s.root = s.buildUI(size)
	s.builded = true
}

func (s *SelectProfileScene) OnExit(next Scene) {}

func (s *SelectProfileScene) Update() error {
	if !s.builded {
		return nil
	}
	s.newProfileButton.SetDisabled(len(s.profiles) == 5) //desabilita caso n maximo de perfis salvos
	s.root.Update(basic.Point{})

	return nil
}

func (s *SelectProfileScene) Draw(screen *ebiten.Image) {
	if s.builded {
		s.root.Draw(screen)
	}
}

func (s *SelectProfileScene) buildUI(size basic.Size) components.Widget {
	title := components.NewText(
		basic.Point{},
		"Jogadores",
		colors.White,
		22, //TODO: PADRONIZAR OS TAMANHOS DE TITULOS DE SCENES
	)

	spacer := components.NewContainer(
		basic.Point{},
		basic.Size{W: size.W, H: 30},
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

	return components.NewColumn(
		basic.Point{X: 0, Y: 0},
		16,
		size,
		basic.Start,
		basic.Center,
		[]components.Widget{spacer, title, spacer, profilesWrappler, spacer, buttonRowWrappler},
	)
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
			s.ctx.Profile = p
			s.stack.Push(&ProfileScene{})

		},
	)

	// Ícone de deletar
	deleteBtn := components.NewDeleteIconButton(basic.Point{}, iconSize, func() {
		_ = service.RemoveProfile(profile.Username)
		s.profiles = service.GetProfiles()
		s.builded = false
		s.root = s.buildUI(s.screenSize)
		s.builded = true
	})

	// Ícone de jogar
	playBtn := components.NewPlayIconButton(basic.Point{}, iconSize, func() {
		s.ctx.Profile = p
		s.stack.Push(&ModeSelectionScene{})
		//s.stack.Push(&PlacementScene{})
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
