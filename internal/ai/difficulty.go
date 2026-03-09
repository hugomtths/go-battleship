package ai

import "github.com/allanjose001/go-battleship/internal/entity"

func NewEasyAIPlayer() *AIPlayer {
	return &AIPlayer{
		Strategies: []Strategy{
			&RandomStrategy{},
		},
	}
}

func NewMediumAIPlayer(enemyFleet *entity.Fleet) *AIPlayer {
	return &AIPlayer{
		enemyFleet: enemyFleet,
		Strategies: []Strategy{
			&PartialLineStrategy{},
			&DiscoveryStrategy{},
			&RandomStrategy{},
		},
	}
}

func NewHardAIPlayer(enemyFleet *entity.Fleet) *AIPlayer {
	return &AIPlayer{
		enemyFleet: enemyFleet,
		Strategies: []Strategy{
			&StrategicSearchStrategy{},
			&FullLineStrategy{},
			&DiscoveryStrategy{},
			&RandomStrategy{},
		},
	}
}

func NewDynamicAIPlayer(enemyFleet *entity.Fleet, ownBoard *entity.Board) *AIPlayer {
	return &AIPlayer{
		enemyFleet:   enemyFleet,
		ownBoard:     ownBoard,
		evasionQueue: make([]*entity.Ship, 0), // inicializa fila vazia
		Strategies: []Strategy{
			//&EvasionStrategy{},
			&RandomMoveStrategy{Chance: 40}, // 2º: move aleatoriamente 40% das vezes
			&StrategicSearchStrategy{},
			&FullLineStrategy{},
			&DiscoveryStrategy{},
			&RandomStrategy{},
		},
	}
}
