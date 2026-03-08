package scenes

import (
	"image/color"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/allanjose001/go-battleship/internal/service"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type DynamicBattlePhase int

const (
	PhasePlayerMove DynamicBattlePhase = iota
	PhasePlayerAttack
	PhaseEnemyTurn
)

// DynamicBattleScene representa a tela de batalha do modo dinâmico.
type DynamicBattleScene struct {
	battleSvc  service.BattleService
	dynamicSvc service.DynamicBattleService

	backButtonContainer components.Widget
	assets              *BattleAssets

	playerHUD *components.BattleHUD
	aiHUD     *components.BattleHUD
	inputCtrl *components.BattleInput // Controle de input para ataque (tabuleiro inimigo)

	// Input para movimento (tabuleiro do jogador)
	selectedShip   *entity.Ship
	selRow, selCol int // Posição do cursor de seleção no tabuleiro do jogador

	// Dragging logic
	isDragging                 bool
	dragStartRow, dragStartCol int

	boardView *components.BattleBoardView
	divider   *components.VerticalDivider

	phase DynamicBattlePhase

	statusText *components.Text

	StackHandler
}

func (s *DynamicBattleScene) GetMusic() string {
	return "battle"
}

func NewDynamicBattleScene() *DynamicBattleScene {
	return &DynamicBattleScene{
		phase: PhasePlayerMove,
	}
}

func (s *DynamicBattleScene) OnEnter(prev Scene, size basic.Size) {
	if s.ctx == nil || s.ctx.Match == nil {
		return
	}

	match := s.ctx.Match
	playerBoard := match.PlayerBoard
	aiBoard := match.EnemyBoard

	if s.ctx.BattleService != nil {
		s.battleSvc = s.ctx.BattleService
	} else if svc, err := service.NewBattleServiceFromMatch(match, s.ctx.IsCampaign); err == nil {
		s.battleSvc = svc
		s.ctx.SetBattleService(svc)
	}

	if s.ctx.DynamicBattleService != nil {
		s.dynamicSvc = s.ctx.DynamicBattleService
	} else {
		s.dynamicSvc = service.NewDynamicBattleService(match)
		s.ctx.SetDynamicBattleService(s.dynamicSvc)
	}

	// Fundo compartilhado
	bg, _, errBg := ebitenutil.NewImageFromFile("assets/images/Mask group.png")
	if errBg == nil {
		playerBoard.BackgroundImage = bg
		aiBoard.BackgroundImage = bg
	}

	// Carrega assets e boardView
	if s.assets == nil {
		s.assets = LoadBattleAssets()
	}
	s.boardView = components.NewBattleBoardView(s.assets.FireAnimation, s.assets.HitImage, s.assets.MissImage)

	lineX := 640.0
	lineY1 := playerBoard.Y
	lineY2 := playerBoard.Y + playerBoard.Size
	s.divider = components.NewVerticalDivider(lineX, lineY1, lineY2, colors.White)

	// Botões de saída
	row := components.NewRow(
		basic.Point{},
		20,
		basic.Size{W: 400, H: 50},
		basic.Center,
		basic.Center,
		[]components.Widget{
			components.NewButton(
				basic.Point{},
				basic.Size{W: 150, H: 50},
				"Desistir",
				colors.Red,
				colors.White,
				func(b *components.Button) {
					SwitchTo(&SelectProfileScene{})
				},
			),
		},
	)

	s.backButtonContainer = components.NewContainer(
		basic.Point{X: 440, Y: 650},
		basic.Size{W: 400, H: 50},
		0,
		colors.Transparent,
		basic.Center,
		basic.Center,
		row,
	)

	if s.assets == nil {
		s.assets = LoadBattleAssets()
	}

	s.boardView = components.NewBattleBoardView(s.assets.FireAnimation, s.assets.HitImage, s.assets.MissImage)

	// HUDs
	playerBaseX := playerBoard.X
	playerBaseY := playerBoard.Y + playerBoard.Size + 40
	aiBaseX := aiBoard.X
	aiBaseY := aiBoard.Y + aiBoard.Size + 40

	playerNameLabel := components.NewText(
		basic.Point{X: float32(playerBaseX + 30), Y: float32(playerBaseY)},
		"Você",
		colors.White,
		20,
	)
	playerAttemptsLabel := components.NewText(
		basic.Point{X: float32(playerBaseX + 30), Y: float32(playerBaseY + 30)},
		"",
		colors.White,
		16,
	)
	playerHitsLabel := components.NewText(
		basic.Point{X: float32(playerBaseX + 30), Y: float32(playerBaseY + 55)},
		"",
		colors.White,
		16,
	)

	aiNameLabel := components.NewText(
		basic.Point{X: float32(aiBaseX + 30), Y: float32(aiBaseY)},
		"IA (Hard)",
		colors.White,
		20,
	)
	aiAttemptsLabel := components.NewText(
		basic.Point{X: float32(aiBaseX + 30), Y: float32(aiBaseY + 30)},
		"",
		colors.White,
		16,
	)
	aiHitsLabel := components.NewText(
		basic.Point{X: float32(aiBaseX + 30), Y: float32(aiBaseY + 55)},
		"",
		colors.White,
		16,
	)

	s.playerHUD = components.NewBattleHUD(playerNameLabel, playerAttemptsLabel, playerHitsLabel, s.battleSvc, components.SidePlayer)
	s.aiHUD = components.NewBattleHUD(aiNameLabel, aiAttemptsLabel, aiHitsLabel, s.battleSvc, components.SideAI)

	s.inputCtrl = components.NewBattleInput(aiBoard)

	s.statusText = components.NewText(basic.Point{X: 400, Y: 30}, "Sua vez: Mova um navio (Clique para selecionar, Setas para mover)", colors.GoldMedal, 20)
}

func (s *DynamicBattleScene) OnExit(next Scene) {}

func (s *DynamicBattleScene) Update() error {
	s.backButtonContainer.Update(basic.Point{})
	if s.playerHUD != nil {
		s.playerHUD.Update(basic.Point{})
	}
	if s.aiHUD != nil {
		s.aiHUD.Update(basic.Point{})
	}

	if s.battleSvc == nil {
		return nil
	}

	switch s.phase {
	case PhasePlayerMove:
		s.statusText.Text = "SUA VEZ: Mova um navio (opcional) ou Ataque!"
		// Permite atacar diretamente (pular movimento)
		s.handlePlayerAttackInput()
		// Se o ataque resultou em troca de turno, encerra o update desta fase
		if s.phase == PhaseEnemyTurn {
			return nil
		}
		// Se ainda é a vez do jogador, processa input de movimento
		s.handlePlayerMoveInput()
	case PhasePlayerAttack:
		s.statusText.Text = "FASE DE ATAQUE: Dispare no inimigo!"
		s.handlePlayerAttackInput()
	case PhaseEnemyTurn:
		s.statusText.Text = "TURNO DO INIMIGO..."
		s.handleEnemyTurn()
	}

	return nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (s *DynamicBattleScene) handlePlayerMoveInput() {
	if s.ctx == nil || s.ctx.Match == nil {
		return
	}

	// Seleção de navio mover
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		row, col := s.getBoardCell(mx, my, s.ctx.Match.PlayerBoard)
		if row != -1 && col != -1 {
			// Tenta encontrar navio na posição (usando o board lógico do player)
			entityBoard := s.ctx.Match.PlayerEntityBoard
			if entityBoard != nil {
				ship := entity.GetShipReference(entityBoard.Positions[row][col])
				if ship != nil && !ship.IsDestroyed() {
					s.selectedShip = ship
					s.selRow = row
					s.selCol = col
					s.isDragging = true
					s.dragStartRow = row
					s.dragStartCol = col
				} else {
					s.selectedShip = nil
				}
			}
		}
	}

	// Fim do Arraste (Movimento)
	if s.isDragging && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		s.isDragging = false
		mx, my := ebiten.CursorPosition()
		endRow, endCol := s.getBoardCell(mx, my, s.ctx.Match.PlayerBoard)

		// Se soltou fora do board ou na mesma casa, não move (apenas selecionou)
		if endRow != -1 && endCol != -1 && (endRow != s.dragStartRow || endCol != s.dragStartCol) {
			// Calcula direção do movimento
			dr := endRow - s.dragStartRow
			dc := endCol - s.dragStartCol

			var dir entity.Direction = -1
			// Prioriza eixo com maior deslocamento
			if abs(dr) >= abs(dc) {
				if dr > 0 {
					dir = entity.Down
				} else {
					dir = entity.Up
				}
			} else {
				if dc > 0 {
					dir = entity.Right
				} else {
					dir = entity.Left
				}
			}

			if dir != -1 {
				// Tenta mover
				err := s.dynamicSvc.HandlePlayerMove(s.selRow, s.selCol, dir)
				if err == nil {
					s.selectedShip = nil
					s.phase = PhasePlayerAttack
				} else {
					// Se falhar (ex: colisão), tenta com findShipPosition (backup)
					r, c := s.findShipPosition(s.selectedShip)
					if r != -1 {
						err = s.dynamicSvc.HandlePlayerMove(r, c, dir)
						if err == nil {
							s.selectedShip = nil
							s.phase = PhasePlayerAttack
						}
					}
				}
			}
		}
	}

	// Movimento com setas
	if s.selectedShip != nil {
		var dir entity.Direction = -1
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			dir = entity.Up
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			dir = entity.Down
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			dir = entity.Left
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			dir = entity.Right
		}

		if dir != -1 {
			// Tenta mover usando a posição clicada/selecionada
			err := s.dynamicSvc.HandlePlayerMove(s.selRow, s.selCol, dir)

			// Se falhar (ex: clicou em parte danificada), tenta achar uma parte intacta automaticamente
			if err != nil {
				r, c := s.findShipPosition(s.selectedShip)
				if r != -1 {
					err = s.dynamicSvc.HandlePlayerMove(r, c, dir)
				}
			}

			if err == nil {
				// Sucesso, passa para ataque
				s.selectedShip = nil
				s.phase = PhasePlayerAttack
			}
		}
	}

	// Botão Pular Movimento (Espaço)
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.selectedShip = nil
		s.phase = PhasePlayerAttack
	}
}

func (s *DynamicBattleScene) findShipPosition(ship *entity.Ship) (int, int) {
	board := s.ctx.Match.PlayerEntityBoard
	// 1. Tenta encontrar uma parte INTACTA do navio
	for r := 0; r < 10; r++ {
		for c := 0; c < 10; c++ {
			if entity.GetShipReference(board.Positions[r][c]) == ship && !entity.IsAttacked(board.Positions[r][c]) {
				return r, c
			}
		}
	}
	// 2. Se não encontrar (improvável se !IsDestroyed), retorna qualquer parte
	for r := 0; r < 10; r++ {
		for c := 0; c < 10; c++ {
			if entity.GetShipReference(board.Positions[r][c]) == ship {
				return r, c
			}
		}
	}
	return -1, -1
}

func (s *DynamicBattleScene) handlePlayerAttackInput() {
	if s.inputCtrl != nil {
		row, col, ok := s.inputCtrl.ClickedCell()
		if ok {
			if res, err := s.battleSvc.HandlePlayerClick(row, col); err == nil && res != nil {
				s.handleMatchEnd(res)
				return
			}
			// Se atacou com sucesso (sem game over), passa a vez
			// O HandlePlayerClick retorna nil se o jogo continua, mas precisamos saber se o turno mudou
			// O BattleService.Stats retorna isPlayerTurn. Se mudar, mudamos a fase.
			_, _, _, _, isPlayerTurn := s.battleSvc.Stats()
			if !isPlayerTurn {
				s.phase = PhaseEnemyTurn
			}
		}
	}
}

func (s *DynamicBattleScene) handleEnemyTurn() {
	// 1. Move
	s.dynamicSvc.HandleEnemyMove()

	// 2. Attack
	if res, err := s.battleSvc.HandleEnemyTurn(); err == nil && res != nil {
		s.handleMatchEnd(res)
		return
	}

	// Verifica se voltou a ser a vez do player
	_, _, _, _, isPlayerTurn := s.battleSvc.Stats()
	if isPlayerTurn {
		s.phase = PhasePlayerMove
	}
}

func (s *DynamicBattleScene) handleMatchEnd(res *entity.MatchResult) {
	winner := s.battleSvc.WinnerName()
	actionLabel := "Voltar ao Menu"
	onAction := func() {
		SwitchTo(&SelectProfileScene{})
	}
	SwitchTo(NewGameOverScene(winner, res, actionLabel, onAction))
}

func (s *DynamicBattleScene) getBoardCell(mx, my int, b *board.Board) (int, int) {
	cellSize := float64(b.Size) / 10.0
	relX := float64(mx) - float64(b.X)
	relY := float64(my) - float64(b.Y)

	if relX >= 0 && relX < float64(b.Size) && relY >= 0 && relY < float64(b.Size) {
		col := int(relX / cellSize)
		row := int(relY / cellSize)
		return row, col
	}
	return -1, -1
}

func (s *DynamicBattleScene) Draw(screen *ebiten.Image) {
	if s.ctx == nil || s.ctx.Match == nil {
		return
	}
	match := s.ctx.Match
	playerBoard := match.PlayerBoard
	aiBoard := match.EnemyBoard

	playerBoard.Draw(screen)
	aiBoard.Draw(screen)

	if s.divider != nil {
		s.divider.Draw(screen)
	}
	if s.playerHUD != nil {
		s.playerHUD.Draw(screen, playerBoard)
	}
	if s.aiHUD != nil {
		s.aiHUD.Draw(screen, aiBoard)
	}

	if s.boardView != nil {
		s.boardView.DrawBoard(screen, playerBoard, match.PlayerShips)
		s.boardView.DrawBoard(screen, aiBoard, nil)
	}

	// Desenha highlight do navio selecionado
	if s.selectedShip != nil && s.phase == PhasePlayerMove {
		// Percorre o board para desenhar highlight sobre o navio selecionado
		cellSize := float64(playerBoard.Size) / 10.0
		for r := 0; r < 10; r++ {
			for c := 0; c < 10; c++ {
				if entity.GetShipReference(match.PlayerEntityBoard.Positions[r][c]) == s.selectedShip {
					// Draw rect
					x := float64(playerBoard.X) + float64(c)*cellSize
					y := float64(playerBoard.Y) + float64(r)*cellSize
					ebitenutil.DrawRect(screen, x, y, cellSize, cellSize, color.RGBA{255, 255, 0, 100})
				}
			}
		}
	}

	s.backButtonContainer.Draw(screen)
	s.statusText.Draw(screen)
}
