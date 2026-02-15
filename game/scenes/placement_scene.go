package scenes

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/game/shared/placement"
	"github.com/allanjose001/go-battleship/game/shared/setup"
	"github.com/allanjose001/go-battleship/game/state"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PlacementScene struct {
	board        *board.Board
	ships        []*placement.ShipPlacement
	selectedShip *placement.ShipPlacement
	activeShip   *placement.ShipPlacement // navio ativo para rotação
	orientation  board.Orientation
	actionsRow   *components.Row
	playRow      *components.Row
	playerLabel  *components.Text
	playButton   *components.Button
	StackHandler
}

func NewPlacementScene() *PlacementScene {
	return &PlacementScene{}
}

/* =======================
   Interface Scene
======================= */

func (s *PlacementScene) OnEnter(prev Scene, size basic.Size) {
	s.board = board.NewBoard(80, 100, 400)
	bg, _, err := ebitenutil.NewImageFromFile("assets/images/Mask group.png")
	if err == nil {
		s.board.BackgroundImage = bg
	}
	s.orientation = board.Horizontal

	img1, _, _ := ebitenutil.NewImageFromFile("assets/images/1 slot 1.png")
	img2, _, _ := ebitenutil.NewImageFromFile("assets/images/3 slots 2.png")
	img3, _, _ := ebitenutil.NewImageFromFile("assets/images/Frame 400.png")
	img4, _, _ := ebitenutil.NewImageFromFile("assets/images/NAVIO 4 SLOTS 1.png")

	s.ships = []*placement.ShipPlacement{
		{Image: img3, Size: 6, ListX: 800, ListY: 100},
		{Image: img4, Size: 4, ListX: 800, ListY: 180},
		{Image: img2, Size: 3, ListX: 800, ListY: 240},
		{Image: img2, Size: 3, ListX: 800, ListY: 300},
		{Image: img1, Size: 1, ListX: 800, ListY: 360},
	}

	btnColor := color.RGBA{48, 67, 103, 255}
	playBtnColor := color.RGBA{60, 120, 60, 255}

	s.actionsRow = components.NewRow(
		basic.Point{X: 80, Y: 580},
		100,
		basic.Size{W: 400, H: 50},
		basic.Start,
		basic.Center,
		[]components.Widget{
			components.NewButton(
				basic.Point{},
				basic.Size{W: 150, H: 50},
				"Aleatório",
				btnColor,
				colors.White,
				func(b *components.Button) { s.randomPlacement() },
			),
			components.NewButton(
				basic.Point{},
				basic.Size{W: 150, H: 50},
				"Rotacionar",
				btnColor,
				colors.White,
				func(b *components.Button) {
					// Determina a nova orientação baseada no navio ativo ou na global
					targetOri := board.Horizontal

					if s.activeShip != nil && s.activeShip.Placed {
						// Se tem navio no board, alterna a partir da orientação DELE
						if s.activeShip.Orientation == board.Horizontal {
							targetOri = board.Vertical
						} else {
							targetOri = board.Horizontal
						}

						// Tenta rotacionar o navio
						s.rotateShipTo(s.activeShip, targetOri)
						// Sincroniza a orientação global com o resultado (se rotacionou ou não)
						s.orientation = s.activeShip.Orientation
					} else {
						// Se não tem navio ativo, apenas alterna a orientação global para o próximo arraste
						if s.orientation == board.Horizontal {
							s.orientation = board.Vertical
						} else {
							s.orientation = board.Horizontal
						}
					}
				},
			),
		},
	)

	s.playButton = components.NewButton(
		basic.Point{},
		basic.Size{W: 150, H: 50},
		"JOGAR",
		playBtnColor,
		colors.White,
		func(b *components.Button) {
			if s.AllShipsPlaced() {
				gs := state.NewGameState()
				gs.PlayerBoard = s.board
				gs.PlayerShips = s.ships

				// Ajusta o tabuleiro da IA para ser simétrico ao do jogador
				gs.AIBoard.X = 1280 - s.board.X - s.board.Size // 1280 - 80 - 400 = 800
				gs.AIBoard.Y = s.board.Y
				gs.AIBoard.Size = s.board.Size

				// Adicionamos navios aleatórios para a IA
				setup.RandomlyPlaceAIShips(gs.AIBoard)

				SwitchTo(NewBattleScene(gs))
			}
		},
	)

	s.playRow = components.NewRow(
		basic.Point{X: 800, Y: 580},
		0,
		basic.Size{W: 400, H: 50},
		basic.Center,
		basic.Center,
		[]components.Widget{
			s.playButton,
		},
	)

	s.playerLabel = components.NewText(
		basic.Point{X: 250, Y: 520},
		"Jogador 1",
		colors.White,
		24,
	)
	// Centraliza o texto em relação ao tabuleiro
	textW := s.playerLabel.GetSize().W
	boardCenter := s.board.X + s.board.Size/2
	newX := boardCenter - float64(textW)/2
	s.playerLabel.SetPos(basic.Point{X: float32(newX), Y: 520})
}

func (s *PlacementScene) OnExit(next Scene) {}

/* =======================
   Update
======================= */

func (s *PlacementScene) Update() error {
	s.actionsRow.Update(basic.Point{})
	s.playRow.Update(basic.Point{})
	s.playerLabel.Update(basic.Point{})

	// Habilita o botão JOGAR apenas se todos os barcos estiverem posicionados

	if s.AllShipsPlaced() {
		s.playButton.ToggleDisabled()

	}

	mx, my := ebiten.CursorPosition()
	mouseX, mouseY := float64(mx), float64(my)

	//Clique para pegar navio (lista ou tabuleiro)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		//Primeiro: tentar pegar navio do tabuleiro
		for _, ship := range s.ships {
			if !ship.Placed {
				continue
			}

			cellSize := s.board.Size / float64(board.Cols)
			x := s.board.X + float64(ship.X)*cellSize
			y := s.board.Y + float64(ship.Y)*cellSize

			w := cellSize * float64(ship.Size)
			h := cellSize
			if ship.Orientation == board.Vertical {
				w, h = h, w
			}

			if mouseX >= x && mouseX <= x+w && mouseY >= y && mouseY <= y+h {
				s.removeShipFromBoard(ship)

				ship.Dragging = true
				ship.DragX = x
				ship.DragY = y
				ship.OffsetX = mouseX - x
				ship.OffsetY = mouseY - y

				s.selectedShip = ship
				s.activeShip = ship
				s.orientation = ship.Orientation //Sincroniza a orientação global com a do navio
				return nil
			}
		}

		//Depois: tentar pegar da lista
		for _, ship := range s.ships {
			if ship.Placed {
				continue
			}

			w, h := ship.Image.Size()
			// Consideramos o tamanho visual na lista (sem rotação por enquanto na lista)
			if mouseX >= ship.ListX && mouseX <= ship.ListX+float64(w) &&
				mouseY >= ship.ListY && mouseY <= ship.ListY+float64(h) {

				ship.Dragging = true
				ship.DragX = ship.ListX
				ship.DragY = ship.ListY
				ship.OffsetX = mouseX - ship.ListX
				ship.OffsetY = mouseY - ship.ListY
				s.selectedShip = ship
				s.activeShip = ship //tornar ativo ao pegar da lista
				return nil
			}
		}
	}

	//Arrastando
	if s.selectedShip != nil && s.selectedShip.Dragging {
		s.selectedShip.DragX = mouseX - s.selectedShip.OffsetX
		s.selectedShip.DragY = mouseY - s.selectedShip.OffsetY
	}

	//Soltar
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && s.selectedShip != nil {
		s.tryPlaceSelectedShip(mouseX, mouseY)
		s.selectedShip = nil
	}

	return nil
}

/* =======================
   Draw
======================= */

func (s *PlacementScene) Draw(screen *ebiten.Image) {
	s.board.Draw(screen)

	for _, ship := range s.ships {
		op := &ebiten.DrawImageOptions{}

		if ship.Placed {
			cellSize := s.board.Size / float64(board.Cols)
			x := s.board.X + float64(ship.X)*cellSize
			y := s.board.Y + float64(ship.Y)*cellSize

			iw, ih := ship.Image.Size()
			op.GeoM.Scale(
				(cellSize*float64(ship.Size))/float64(iw),
				cellSize/float64(ih),
			)

			if ship.Orientation == board.Vertical {
				op.GeoM.Rotate(math.Pi / 2)
				op.GeoM.Translate(cellSize, 0)
			}

			op.GeoM.Translate(x, y)

			// Destaque visual para o navio ativo (no tabuleiro)
			if s.activeShip == ship {
				var highlightColor color.RGBA = color.RGBA{255, 255, 0, 160} // Amarelo mais visível
				highlightW := cellSize * float64(ship.Size)
				highlightH := cellSize
				if ship.Orientation == board.Vertical {
					highlightW, highlightH = highlightH, highlightW
				}
				ebitenutil.DrawRect(screen, x-2, y-2, highlightW+4, highlightH+4, highlightColor)
			}

			screen.DrawImage(ship.Image, op)
			continue
		}

		if ship.Dragging {
			cellSize := s.board.Size / float64(board.Cols)
			iw, ih := ship.Image.Size()

			op.GeoM.Scale(
				(cellSize*float64(ship.Size))/float64(iw),
				cellSize/float64(ih),
			)

			if s.orientation == board.Vertical {
				op.GeoM.Rotate(math.Pi / 2)
				op.GeoM.Translate(cellSize, 0)
			}

			op.GeoM.Translate(ship.DragX, ship.DragY)

			// 🔥 Destaque visual para o navio ativo (enquanto arrasta)
			if s.activeShip == ship {
				var highlightColor color.RGBA = color.RGBA{255, 255, 0, 160}
				highlightW := cellSize * float64(ship.Size)
				highlightH := cellSize
				if s.orientation == board.Vertical {
					highlightW, highlightH = highlightH, highlightW
				}
				ebitenutil.DrawRect(screen, ship.DragX-2, ship.DragY-2, highlightW+4, highlightH+4, highlightColor)
			}
		} else {
			op.GeoM.Translate(ship.ListX, ship.ListY)

			//  Destaque visual para o navio ativo (na lista)
			if s.activeShip == ship {
				w, h := ship.Image.Size()
				var highlightColor color.RGBA = color.RGBA{255, 255, 0, 100}
				ebitenutil.DrawRect(screen, ship.ListX-2, ship.ListY-2, float64(w)+4, float64(h)+4, highlightColor)
			}
		}

		screen.DrawImage(ship.Image, op)
	}

	s.actionsRow.Draw(screen)
	s.playRow.Draw(screen)
	s.playerLabel.Draw(screen)

	// Linha vertical separando tabuleiro e navios
	lineX := 640.0
	lineY1 := s.board.Y
	lineY2 := float64(s.actionsRow.GetPos().Y) + float64(s.actionsRow.GetSize().H)
	ebitenutil.DrawLine(screen, lineX, lineY1, lineX, lineY2, colors.White)
}

/* =======================
   Helpers
======================= */

func (s *PlacementScene) tryPlaceSelectedShip(mouseX, mouseY float64) {
	ship := s.selectedShip
	ship.Dragging = false

	cellSize := s.board.Size / float64(board.Cols)

	// Calcula a posição pretendida do topo-esquerdo do navio no tabuleiro
	// baseada na posição atual de arraste
	targetX := ship.DragX
	targetY := ship.DragY

	// Se estiver perto o suficiente do tabuleiro, tenta encaixar
	if targetX+cellSize*float64(ship.Size) > s.board.X && targetX < s.board.X+s.board.Size &&
		targetY+cellSize*float64(ship.Size) > s.board.Y && targetY < s.board.Y+s.board.Size {

		col := int(math.Round((targetX - s.board.X) / cellSize))
		row := int(math.Round((targetY - s.board.Y) / cellSize))

		if s.board.CanPlace(ship.Size, row, col, s.orientation) {
			s.board.PlaceShip(ship.Size, row, col, s.orientation)
			ship.Placed = true
			ship.X = col
			ship.Y = row
			ship.Orientation = s.orientation
			s.activeShip = ship
			return
		}
	}

	//Drop inválido ou fora do tabuleiro -> volta pra lista
	ship.Placed = false
	s.activeShip = ship // Mesmo na lista, ele pode ser o ativo para rotacionar antes de puxar
}

func (s *PlacementScene) removeShipFromBoard(target *placement.ShipPlacement) {
	target.Placed = false
	s.board.Clear()

	for _, ship := range s.ships {
		if ship != target && ship.Placed {
			s.board.PlaceShip(ship.Size, ship.Y, ship.X, ship.Orientation)
		}
	}
}

func (s *PlacementScene) randomPlacement() {
	s.board.Clear()
	rand.Seed(time.Now().UnixNano())

	for _, ship := range s.ships {
		ship.Placed = false
	}

	for _, ship := range s.ships {
		for {
			row := rand.Intn(board.Rows)
			col := rand.Intn(board.Cols)
			or := board.Orientation(rand.Intn(2))

			if s.board.CanPlace(ship.Size, row, col, or) {
				s.board.PlaceShip(ship.Size, row, col, or)
				ship.Placed = true
				ship.X = col
				ship.Y = row
				ship.Orientation = or
				s.activeShip = ship //último navio colocado vira ativo
				break
			}
		}
	}
}

func (s *PlacementScene) findLastPlacedShip() {
	// Procura o último navio colocado (mais recente)
	for i := len(s.ships) - 1; i >= 0; i-- {
		if s.ships[i].Placed {
			s.activeShip = s.ships[i]
			return
		}
	}
	s.activeShip = nil
}

func (s *PlacementScene) rotateShipTo(ship *placement.ShipPlacement, newOri board.Orientation) {
	if ship == nil || !ship.Placed {
		return
	}

	// Se já estiver na orientação certa, não faz nada
	if ship.Orientation == newOri {
		return
	}

	s.removeShipFromBoard(ship)

	if s.board.CanPlace(ship.Size, ship.Y, ship.X, newOri) {
		s.board.PlaceShip(ship.Size, ship.Y, ship.X, newOri)
		ship.Orientation = newOri
		ship.Placed = true
	} else {
		// Se não couber, volta para a orientação original
		s.board.PlaceShip(ship.Size, ship.Y, ship.X, ship.Orientation)
		ship.Placed = true

		// Opcional: Se não coube na orientação global,
		// talvez devêssemos reverter a s.orientation global?
		// Por enquanto deixamos assim para o usuário saber que ele tentou mudar.
	}
}

func (s *PlacementScene) AllShipsPlaced() bool {
	for _, ship := range s.ships {
		if !ship.Placed {
			return false
		}
	}
	return true
}

var _ Scene = (*PlacementScene)(nil)
