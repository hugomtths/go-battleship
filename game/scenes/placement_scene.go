package scenes

import (
	"fmt"
	"image/color"
	"time"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/game/shared/placement"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/allanjose001/go-battleship/internal/service"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// placementRenderer é responsável por desenhar o tabuleiro de posicionamento
// e todos os navios do jogador, delegando o desenho de cada navio
// para o componente DrawShip.
type placementRenderer struct {
	screen *ebiten.Image
}

// Draw desenha o tabuleiro e percorre a lista de navios,
// indicando qual deles está ativo para ser destacado.
func (r *placementRenderer) Draw(b *board.Board, ships []*placement.ShipPlacement, active *placement.ShipPlacement, orientation board.Orientation) {
	b.Draw(r.screen)
	for _, ship := range ships {
		components.DrawShip(r.screen, b, ship, active == ship, orientation)
	}
}

// PlacementScene controla a tela onde o jogador posiciona seus navios
// antes de iniciar a batalha. Ela orquestra a interação com o serviço
// de placement e os componentes de interface.
type PlacementScene struct {
	// svc encapsula toda a regra de negócio de posicionamento
	svc service.PlacementService
	// board e ships são usados para construir o GameState de batalha
	board *board.Board
	ships []*placement.ShipPlacement
	// perfil do jogador selecionado na tela anterior
	playerProfile *entity.Profile
	// container com a linha de botões sob o tabuleiro (Aleatório, Rotacionar)
	leftButtons components.Widget
	// playerLabel mostra o texto "Jogador 1" alinhado ao tabuleiro
	playerLabel *components.Text
	// playButton é o botão que dispara a transição para a batalha
	playButton *components.Button
	// container com a linha de start sob a coluna de navios
	rightButtons components.Widget

	// Elementos decorativos (ex: Título da partida em modo campanha)
	decorations []components.Widget

	// Estado da Série (Melhor de 3)
	matchIndex        int // 1, 2 ou 3
	seriesScorePlayer int
	seriesScoreEnemy  int
	StackHandler
}

func (s *PlacementScene) GetMusic() string {
	return "menus" //TODO Procurar musica para placement
}

// NewPlacementScene cria uma cena de posicionamento vazia.
// A configuração completa é feita em OnEnter.
func NewPlacementScene() *PlacementScene { return &PlacementScene{} }

// NewPlacementSceneWithProfile cria a cena já com o perfil do jogador selecionado
func NewPlacementSceneWithProfile(p *entity.Profile) *PlacementScene {
	return &PlacementScene{playerProfile: p}
}

// SetSeriesState configura o estado da série de partidas (ex: partida 2 de 3)
func (s *PlacementScene) SetSeriesState(index, pWins, eWins int) {
	s.matchIndex = index
	s.seriesScorePlayer = pWins
	s.seriesScoreEnemy = eWins
}

// OnEnter é chamado quando a cena entra em foco.
// Aqui criamos o tabuleiro, carregamos imagens, configuramos navios,
// serviços e os botões de interface.
func (s *PlacementScene) OnEnter(prev Scene, size basic.Size) {
	// Cria o tabuleiro do jogador na tela
	s.decorations = []components.Widget{}
	b := board.NewBoard(80, 100, 400)

	// Tenta carregar a imagem de fundo do tabuleiro
	bg, _, err := ebitenutil.NewImageFromFile("assets/images/Mask group.png")
	if err == nil {
		b.BackgroundImage = bg
	}

	// Carrega os sprites dos navios com tamanhos diferentes
	img1, _, _ := ebitenutil.NewImageFromFile("assets/images/1 slot 1.png")
	img2, _, _ := ebitenutil.NewImageFromFile("assets/images/3 slots 2.png")
	img3, _, _ := ebitenutil.NewImageFromFile("assets/images/Frame 400.png")
	img4, _, _ := ebitenutil.NewImageFromFile("assets/images/NAVIO 4 SLOTS 1.png")

	ships := []*placement.ShipPlacement{
		{Image: img3, Size: 6, ListX: 800, ListY: 100},
		{Image: img3, Size: 6, ListX: 800, ListY: 160},
		{Image: img4, Size: 4, ListX: 800, ListY: 250},
		{Image: img4, Size: 4, ListX: 800, ListY: 320},
		{Image: img2, Size: 3, ListX: 800, ListY: 380},
		{Image: img1, Size: 1, ListX: 800, ListY: 440},
	}

	s.board = b
	s.ships = ships

	// Cria o serviço de placement, responsável apenas por posicionamento visual
	s.svc = service.NewPlacementService(b, ships)

	// Cores dos botões
	btnColor := color.RGBA{48, 67, 103, 255}
	playBtnColor := color.RGBA{60, 120, 60, 255}

	// Botão Voltar
	backButton := components.NewButton(
		basic.Point{},
		basic.Size{W: 120, H: 50},
		"Voltar",
		btnColor,
		colors.White,
		func(b *components.Button) {
			s.stack.Pop()
		},
	)

	// Botão para iniciar a batalha. Só funciona se todos os navios
	// tiverem sido posicionados no tabuleiro.
	s.playButton = components.NewButton(
		basic.Point{},
		basic.Size{W: 150, H: 50},
		"Partida",
		playBtnColor,
		colors.White,
		func(b *components.Button) {
			if !s.svc.AllShipsPlaced() {
				return
			}

			factory := service.NewGameService()
			gs := factory.NewBattleGameState(s.board, s.ships)

			matchID := fmt.Sprintf("match-%d", time.Now().UnixNano())

			diff := "easy"
			if s.stack.ctx != nil && s.stack.ctx.Difficulty != "" {
				diff = s.stack.ctx.Difficulty
			}

			match := entity.NewMatch(matchID, diff, gs.PlayerBoard, gs.AIBoard, s.ships, s.playerProfile)

			svc, err := service.NewBattleServiceFromMatch(match)
			if err != nil {
				fmt.Println("Erro ao criar serviço de batalha:", err)
				return
			}

			if s.stack.ctx != nil {
				s.stack.ctx.SetMatch(match)
				s.stack.ctx.SetBattleService(svc)
			}

			// Configura a cena de batalha com o estado da série
			battleScene := NewBattleScene()
			if s.ctx != nil && s.ctx.IsCampaign {
				battleScene.SetSeriesState(s.matchIndex, s.seriesScorePlayer, s.seriesScoreEnemy)
			}

			SwitchTo(battleScene)
		},
	)

	x, y, sizeX := s.svc.BoardRect()
	leftRow := components.NewRow(
		basic.Point{}, // pos relativo ao container
		50,
		basic.Size{W: float32(sizeX), H: 50},
		basic.Center,
		basic.Center,
		[]components.Widget{
			components.NewButton(
				basic.Point{},
				basic.Size{W: 150, H: 50},
				"Aleatório",
				btnColor,
				colors.White,
				func(b *components.Button) {
					s.svc.RandomPlacement()
				},
			),
			components.NewButton(
				basic.Point{},
				basic.Size{W: 150, H: 50},
				"Rotacionar",
				btnColor,
				colors.White,
				func(b *components.Button) {
					s.svc.Rotate()
				},
			),
		},
	)
	s.leftButtons = components.NewContainer(
		basic.Point{X: float32(x), Y: float32(y + sizeX + 80)},
		basic.Size{W: float32(sizeX), H: 50},
		25,
		nil,
		basic.Center,
		basic.Center,
		leftRow,
	)
	minListX := ships[0].ListX
	for _, sp := range ships {
		if sp.ListX < minListX {
			minListX = sp.ListX
		}
	}
	rightRow := components.NewRow(
		basic.Point{},
		20, // Espaçamento entre botões
		basic.Size{W: 300, H: 50},
		basic.Center,
		basic.Center,
		[]components.Widget{
			backButton,
			s.playButton,
		},
	)
	s.rightButtons = components.NewContainer(
		basic.Point{X: float32(minListX) - 40, Y: float32(y + sizeX + 80)}, // Ajuste X para caber melhor
		basic.Size{W: 300, H: 50},
		25,
		nil,
		basic.Center,
		basic.Center,
		rightRow,
	)

	// Tenta recuperar profile do contexto se não tiver sido passado
	if s.playerProfile == nil && s.stack.ctx != nil && s.stack.ctx.Profile != nil {
		s.playerProfile = s.stack.ctx.Profile
	}

	// Cria o rótulo com o nome do jogador
	labelText := "Jogador 1"
	if s.playerProfile != nil && s.playerProfile.Username != "" {
		labelText = s.playerProfile.Username
	}
	s.playerLabel = components.NewText(
		basic.Point{X: 250, Y: 520},
		labelText,
		colors.White,
		24,
	)
	// Centraliza o texto em relação ao tabuleiro
	textW := s.playerLabel.GetSize().W
	boardCenter := x + sizeX/2
	newX := boardCenter - float64(textW)/2
	s.playerLabel.SetPos(basic.Point{X: float32(newX), Y: 520})

	// Se estiver em modo campanha (matchIndex > 0), exibe info da série
	if s.matchIndex > 0 {
		aiName := "IA"
		if s.ctx != nil {
			switch s.ctx.Difficulty {
			case "easy":
				aiName = "Recruta Bot"
			case "medium":
				aiName = "Imediato Bot"
			case "hard":
				aiName = "Almirante Bot"
			}
		}

		pName := "Você"
		if s.playerProfile != nil && s.playerProfile.Username != "" {
			pName = s.playerProfile.Username
		}

		line1 := fmt.Sprintf("Partida %d/3", s.matchIndex)
		line2 := fmt.Sprintf("%s %d X %d %s", pName, s.seriesScorePlayer, s.seriesScoreEnemy, aiName)

		t1 := components.NewText(basic.Point{}, line1, colors.White, 24)
		t2 := components.NewText(basic.Point{}, line2, colors.GoldMedal, 28)

		centerX := float32(size.W) / 2
		yBase := float32(650)

		t1.SetPos(basic.Point{X: centerX - float32(t1.GetSize().W)/2, Y: yBase})
		t2.SetPos(basic.Point{X: centerX - float32(t2.GetSize().W)/2, Y: yBase + 30})

		s.decorations = append(s.decorations, t1, t2)
	}
	s.stack.ctx.CanPopOrPush = true
	_ = s.Update()
}

// OnExit é chamado ao sair da cena de placement.
// Não há limpeza especial necessária neste caso.
func (s *PlacementScene) OnExit(next Scene) {
	s.stack.ctx.CanPopOrPush = false
}

// Update é chamado a cada frame para tratar entradas do usuário.
// Aqui atualizamos botões, rótulo e delegamos para o serviço de placement
// as interações de clique/arraste dos navios.
func (s *PlacementScene) Update() error {
	if s.leftButtons != nil {
		s.leftButtons.Update(basic.Point{})
	}
	if s.rightButtons != nil {
		s.rightButtons.Update(basic.Point{})
	}
	s.playerLabel.Update(basic.Point{})
	for _, d := range s.decorations {
		d.Update(basic.Point{})
	}

	// Pega a posição atual do mouse em coordenadas de tela
	mx, my := ebiten.CursorPosition()
	mouseX, mouseY := float64(mx), float64(my)

	// Clique do mouse: tenta selecionar um navio no tabuleiro
	// ou na lista da lateral direita.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		if s.svc.SelectOnBoard(mouseX, mouseY) {
			return nil
		}

		if s.svc.SelectOnList(mouseX, mouseY) {
			return nil
		}
	}

	// Atualiza posição do navio que está sendo arrastado
	s.svc.UpdateDragging(mouseX, mouseY)

	// Ao soltar o botão do mouse, tenta soltar o navio no tabuleiro
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		s.svc.DropSelected()
	}

	return nil
}

// Draw é responsável por desenhar todo o conteúdo da cena:
// tabuleiro, navios, botões, rótulo do jogador e a linha divisória.
func (s *PlacementScene) Draw(screen *ebiten.Image) {
	// Renderer que sabe como desenhar tabuleiro e navios
	r := &placementRenderer{screen: screen}
	s.svc.Draw(r)

	if s.leftButtons != nil {
		s.leftButtons.Draw(screen)
	}
	if s.rightButtons != nil {
		s.rightButtons.Draw(screen)
	}
	s.playerLabel.Draw(screen)
	for _, d := range s.decorations {
		d.Draw(screen)
	}

	// Linha vertical separando tabuleiro e navios da lateral
	lineX := 640.0
	_, y, sizeY := s.svc.BoardRect()
	lineY1 := y
	lineY2 := y + sizeY
	ebitenutil.DrawLine(screen, lineX, lineY1, lineX, lineY2, colors.White)
}

// Verifica em tempo de compilação se PlacementScene implementa Scene.
var _ Scene = (*PlacementScene)(nil)
