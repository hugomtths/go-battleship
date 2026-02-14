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
