package scenes

import (
	"fmt"
	"image/color"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/allanjose001/go-battleship/internal/service"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// DynamicBattleScene representa a tela de batalha dinâmica.
// Permite que o jogador mova seus navios durante o turno.
type DynamicBattleScene struct {
	BattleScene
	dynamicBattleSvc service.DynamicBattleService
	playerInputCtrl  *components.BattleInput

	selectedShip    *entity.Ship
	selectedShipIdx int // índice no match.PlayerShips para sincronização visual
}

func NewDynamicBattleScene() *DynamicBattleScene {
	return &DynamicBattleScene{
		selectedShipIdx: -1, // Começa com nenhum navio selecionado
	}
}

func (s *DynamicBattleScene) OnEnter(prev Scene, size basic.Size) {
	if s.ctx == nil || s.ctx.Match == nil {
		return
	}

	match := s.ctx.Match

	// Cria o DynamicBattleService ANTES de chamar BattleScene.OnEnter,
	// para que o pai não instancie um serviço comum (sem ownBoard).
	svc, err := service.NewDynamicBattleServiceFromMatch(match, s.ctx.IsCampaign)
	if err == nil {
		s.dynamicBattleSvc = svc
		s.battleSvc = svc         // compatibilidade com BattleScene
		s.ctx.BattleService = svc // injeta no contexto para BattleScene.OnEnter não recriar
	}

	// Agora o pai encontra s.ctx.BattleService != nil e reutiliza sem criar outro AIPlayer
	s.BattleScene.OnEnter(prev, size)

	// Controlador de entrada para o tabuleiro do jogador (para seleção)
	s.playerInputCtrl = components.NewBattleInput(match.PlayerBoard)
}

func (s *DynamicBattleScene) Update() error {
	// 1. Atualiza elementos base (HUDs, botões, etc)
	s.backButtonContainer.Update(basic.Point{})
	if s.playerHUD != nil {
		s.playerHUD.Update(basic.Point{})
	}
	if s.aiHUD != nil {
		s.aiHUD.Update(basic.Point{})
	}

	if s.dynamicBattleSvc == nil {
		return nil
	}

	// 2. Lógica de Turno do Jogador
	isPlayerTurn := s.ctx.Match.Turn == entity.TurnPlayer

	if isPlayerTurn {
		// A. Seleção de Navio (clique no tabuleiro do jogador)
		if row, col, ok := s.playerInputCtrl.ClickedCell(); ok {
			s.handleShipSelection(row, col)
		}

		// B. Movimentação de Navio (teclado)
		if s.selectedShip != nil {
			if err := s.handleShipMovement(); err != nil {
				fmt.Printf("Erro ao mover navio: %v\n", err)
			}
		}

		// C. Ataque (clique no tabuleiro da IA)
		if row, col, ok := s.inputCtrl.ClickedCell(); ok {
			if res, err := s.dynamicBattleSvc.HandlePlayerClick(row, col); err == nil && res != nil {
				s.handleMatchEnd(res)
				return nil
			}
		}
	}

	// 3. Lógica de Turno da IA
	if res, err := s.dynamicBattleSvc.HandleEnemyTurn(); err == nil && res != nil {
		s.handleMatchEnd(res)
		return nil
	}

	return nil
}

func (s *DynamicBattleScene) handleShipSelection(row, col int) {
	match := s.ctx.Match
	if match.PlayerEntityBoard == nil {
		return
	}

	// Busca o navio na posição clicada no board lógico
	pos := match.PlayerEntityBoard.Positions[row][col]
	ship := entity.GetShipReference(pos)

	if ship != nil {
		// Verifica se o navio foi atingido (não pode mover se foi atingido, conforme regra)
		if ship.HitCount > 0 {
			fmt.Println("Não pode mover um navio que já foi atingido!")
			s.selectedShip = nil
			return
		}

		s.selectedShip = ship

		// Encontra a posição top-left atual do navio no board lógico
		minR, minC := 10, 10
		for r := 0; r < entity.BoardSize; r++ {
			for c := 0; c < entity.BoardSize; c++ {
				if entity.GetShipReference(match.PlayerEntityBoard.Positions[r][c]) == ship {
					if r < minR {
						minR = r
					}
					if c < minC {
						minC = c
					}
				}
			}
		}

		// Encontra o índice correspondente no match.PlayerShips para sincronização visual
		for i, ps := range match.PlayerShips {
			if ps != nil && ps.X == minC && ps.Y == minR {
				s.selectedShipIdx = i
				break
			}
		}
		fmt.Printf("Navio selecionado: %s em %d,%d\n", ship.Name, minR, minC)
	} else {
		s.selectedShip = nil
		s.selectedShipIdx = -1
	}
}

func (s *DynamicBattleScene) handleShipMovement() error {
	var dr, dc int
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		dr = -1
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		dr = 1
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		dc = -1
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		dc = 1
	}

	if dr == 0 && dc == 0 {
		return nil
	}

	match := s.ctx.Match
	// Precisamos encontrar a posição top-left atual do navio para o serviço
	var minR, minC int = 10, 10
	found := false
	for r := 0; r < entity.BoardSize; r++ {
		for c := 0; c < entity.BoardSize; c++ {
			if entity.GetShipReference(match.PlayerEntityBoard.Positions[r][c]) == s.selectedShip {
				if r < minR {
					minR = r
				}
				if c < minC {
					minC = c
				}
				found = true
			}
		}
	}

	if !found {
		return fmt.Errorf("navio não encontrado no tabuleiro")
	}

	newRow := minR + dr
	newCol := minC + dc

	// Tenta mover via serviço
	if err := s.dynamicBattleSvc.MovePlayerShip(s.selectedShip, newRow, newCol); err != nil {
		return err
	}

	// Sincronização Visual: Atualiza o ShipPlacement para o DrawBoard refletir a nova posição
	if s.selectedShipIdx >= 0 && s.selectedShipIdx < len(match.PlayerShips) {
		ps := match.PlayerShips[s.selectedShipIdx]
		ps.X = newCol
		ps.Y = newRow
	}

	// Desmarca o navio após o movimento
	s.selectedShip = nil
	s.selectedShipIdx = -1

	fmt.Printf("Navio movido para %d, %d\n", newRow, newCol)
	return nil
}

func (s *DynamicBattleScene) Draw(screen *ebiten.Image) {
	// Reutiliza o Draw da BattleScene
	s.BattleScene.Draw(screen)

	// Adiciona destaque para o navio selecionado, se houver
	if s.selectedShip != nil && s.selectedShipIdx >= 0 {
		match := s.ctx.Match
		ps := match.PlayerShips[s.selectedShipIdx]
		cellSize := match.PlayerBoard.Size / float64(board.Cols)

		// Desenha um retângulo de seleção ao redor do navio
		rectW := cellSize
		rectH := cellSize
		if ps.Orientation == board.Horizontal {
			rectW *= float64(ps.Size)
		} else {
			rectH *= float64(ps.Size)
		}

		ebitenutil.DrawRect(screen,
			match.PlayerBoard.X+float64(ps.X)*cellSize,
			match.PlayerBoard.Y+float64(ps.Y)*cellSize,
			rectW, rectH,
			color.RGBA{255, 255, 255, 100})
	}
}
