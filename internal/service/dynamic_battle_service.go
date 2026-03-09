package service

import (
	"fmt"
	"time"

	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/internal/ai"
	"github.com/allanjose001/go-battleship/internal/entity"
)

// DynamicBattleService estende BattleService para incluir movimentação de navios.
type DynamicBattleService interface {
	BattleService
	MovePlayerShip(ship *entity.Ship, newRow, newCol int) error
}

type dynamicBattleService struct {
	*battleService
	dynamicMatchSvc *DynamicMatchService
}

// NewDynamicBattleServiceFromMatch inicializa o serviço de batalha dinâmica.
func NewDynamicBattleServiceFromMatch(match *entity.Match, isCampaign bool) (DynamicBattleService, error) {
	match.IsDynamicMode = true // Força a flag de modo dinâmico no objeto Match
	// Usamos DynamicMatchService em vez do MatchService comum
	dynamicMatchSvc := NewDynamicMatchService(NewAttackService(), 500*time.Millisecond)

	var aiPlayer *ai.AIPlayer

	if match.PlayerEntityBoard == nil {
		// Inicialização para novo jogo: precisamos converter PlayerShips (visual) para PlayerEntityBoard/Fleet (lógico)
		playerFleet := entity.NewFleet()
		playerEntityBoard := &entity.Board{}

		// Mapeamento dos navios posicionados para a estrutura lógica
		usedShips := make(map[int]bool)
		for _, ps := range match.PlayerShips {
			if ps == nil || !ps.Placed {
				continue
			}

			var entShip *entity.Ship
			for i, ship := range playerFleet.Ships {
				if !usedShips[i] && ship.Size == ps.Size {
					entShip = ship
					usedShips[i] = true
					break
				}
			}

			if entShip != nil {
				entShip.Horizontal = ps.Orientation == board.Horizontal
				playerEntityBoard.PlaceShip(entShip, ps.Y, ps.X)
			}
		}

		match.PlayerFleet = playerFleet
		match.PlayerEntityBoard = playerEntityBoard

		// Inicialização da IA
		aiFleet := entity.NewFleet()
		aiBoard := &entity.Board{}

		// Posiciona navios da IA aleatoriamente
		aiFleetSvc := NewAIFleetService()
		aiFleetSvc.PositionShipsRandomly(aiBoard, aiFleet)

		match.EnemyFleet = aiFleet
		match.EnemyEntityBoard = aiBoard

		// Cria o DynamicAIPlayer APÓS atribuir EnemyEntityBoard ao match
		aiPlayer = ai.NewDynamicAIPlayer(match.PlayerFleet, match.EnemyEntityBoard)

		// Adicione este log para confirmar
		if match.EnemyEntityBoard == nil {
			return nil, fmt.Errorf("EnemyEntityBoard nil após atribuição")
		}

		totalPlayerCells := 0
		for _, ship := range playerFleet.GetFleetShips() {
			if ship != nil {
				totalPlayerCells += ship.Size
			}
		}

		totalEnemyCells := 0
		for _, ship := range aiFleet.GetFleetShips() {
			if ship != nil {
				totalEnemyCells += ship.Size
			}
		}

		if err := dynamicMatchSvc.Start(
			match,
			time.Now(),
			match.PlayerBoard,
			match.EnemyBoard,
			playerEntityBoard,
			aiBoard,
			playerFleet,
			aiFleet,
			totalEnemyCells,
			totalPlayerCells,
		); err != nil {
			return nil, err
		}
	} else {
		// Caso já esteja inicializado (retomada de estado)
		aiPlayer = ai.NewDynamicAIPlayer(match.PlayerFleet, match.EnemyEntityBoard)
	}

	baseSvc := &battleService{
		matchSvc:   dynamicMatchSvc.MatchService,
		match:      match,
		aiPlayer:   aiPlayer,
		profile:    match.Profile,
		isCampaign: isCampaign,
	}

	return &dynamicBattleService{
		battleService:   baseSvc,
		dynamicMatchSvc: dynamicMatchSvc,
	}, nil
}

func (s *dynamicBattleService) MovePlayerShip(ship *entity.Ship, newRow, newCol int) error {
	if s.match == nil || s.dynamicMatchSvc == nil {
		return ErrMatchNotReady
	}

	err := s.dynamicMatchSvc.MovePlayerShip(s.match, ship, newRow, newCol, time.Now())
	if err != nil {
		return err
	}

	return nil
}
