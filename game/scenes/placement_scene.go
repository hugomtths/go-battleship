package scenes

import (
	"image/color"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/game/shared/placement"
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
	// container com a linha de botões sob o tabuleiro (Aleatório, Rotacionar)
	leftButtons components.Widget
	// playerLabel mostra o texto "Jogador 1" alinhado ao tabuleiro
	playerLabel *components.Text
	// playButton é o botão que dispara a transição para a batalha
	playButton *components.Button
	// container com a linha de start sob a coluna de navios
	rightButtons components.Widget
	StackHandler
}

// NewPlacementScene cria uma cena de posicionamento vazia.
// A configuração completa é feita em OnEnter.
func NewPlacementScene() *PlacementScene {
	return &PlacementScene{}
}

// OnEnter é chamado quando a cena entra em foco.
// Aqui criamos o tabuleiro, carregamos imagens, configuramos navios,
// serviços e os botões de interface.
func (s *PlacementScene) OnEnter(prev Scene, size basic.Size) {
	// Cria o tabuleiro do jogador na tela
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

	// Cria a lista de navios disponíveis na lateral direita
	ships := []*placement.ShipPlacement{
		{Image: img3, Size: 6, ListX: 800, ListY: 100},
		{Image: img4, Size: 4, ListX: 800, ListY: 180},
		{Image: img2, Size: 3, ListX: 800, ListY: 240},
		{Image: img2, Size: 3, ListX: 800, ListY: 300},
		{Image: img1, Size: 1, ListX: 800, ListY: 360},
	}

	s.board = b
	s.ships = ships

	// Cria o serviço de placement, responsável apenas por posicionamento visual
	s.svc = service.NewPlacementService(b, ships)

	// Cores dos botões
	btnColor := color.RGBA{48, 67, 103, 255}
	playBtnColor := color.RGBA{60, 120, 60, 255}

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
			SwitchTo(NewBattleScene(gs))
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
		0,
		basic.Size{W: 200, H: 50},
		basic.Center,
		basic.Center,
		[]components.Widget{
			s.playButton,
		},
	)
	s.rightButtons = components.NewContainer(
		basic.Point{X: float32(minListX), Y: float32(y + sizeX + 80)},
		basic.Size{W: 200, H: 50},
		25,
		nil,
		basic.Center,
		basic.Center,
		rightRow,
	)

	// Cria o rótulo com o nome do jogador
	s.playerLabel = components.NewText(
		basic.Point{X: 250, Y: 520},
		"Jogador 1",
		colors.White,
		24,
	)
	// Centraliza o texto em relação ao tabuleiro
	textW := s.playerLabel.GetSize().W
	boardCenter := x + sizeX/2
	newX := boardCenter - float64(textW)/2
	s.playerLabel.SetPos(basic.Point{X: float32(newX), Y: 520})
}

// OnExit é chamado ao sair da cena de placement.
// Não há limpeza especial necessária neste caso.
func (s *PlacementScene) OnExit(next Scene) {}

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

	// Linha vertical separando tabuleiro e navios da lateral
	lineX := 640.0
	_, y, sizeY := s.svc.BoardRect()
	lineY1 := y
	lineY2 := y + sizeY
	ebitenutil.DrawLine(screen, lineX, lineY1, lineX, lineY2, colors.White)
}

// Verifica em tempo de compilação se PlacementScene implementa Scene.
var _ Scene = (*PlacementScene)(nil)
