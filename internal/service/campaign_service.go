package service

import (
	"fmt"
	"time"

	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/internal/ai"
	"github.com/allanjose001/go-battleship/internal/entity"
)

// CampaignService gerencia o fluxo do modo campanha.
type CampaignService struct {
	matchService *MatchService
}

func NewCampaignService(ms *MatchService) *CampaignService {
	return &CampaignService{
		matchService: ms,
	}
}

// StartCampaignMatch configura a partida de campanha.
func (cs *CampaignService) StartCampaignMatch(
	username string,
	fleet *entity.Fleet,
	playerBoard *board.Board,
	enemyBoard *board.Board,
	playerEntityBoard *entity.Board,
	enemyEntityBoard *entity.Board,
	enemyFleet *entity.Fleet,
	enemyShipCells int,
	playerShipCells int,
) (*ai.AIPlayer, error) {
	profile, err := FindProfile(username)
	if err != nil {
		return nil, err
	}

	// Inicializa campanha se for o primeiro acesso do usuário
	if profile.CurrentCampaign == nil {
		profile.CurrentCampaign = &entity.Campaign{
			ID:             fmt.Sprintf("camp_%s_%d", username, time.Now().Unix()),
			DifficultyStep: make(map[string]entity.MatchResult),
			IsActive:       true,
		}
	}

	// Identifica a dificuldade atual (Easy -> Medium -> Hard)
	diff, finished := cs.GetNextDifficulty(profile.CurrentCampaign)
	if finished {
		return nil, fmt.Errorf("campanha para %s já foi concluída", username)
	}

	// Cria a IA correspondente ao nível atual
	opponent := cs.selectAI(diff, fleet)

	// Cria o objeto de partida em memória
	match := cs.matchService.Create(profile.CurrentCampaign.ID, diff)

	// Inicia os estados da partida (Turnos, Timers, etc)
	err = cs.matchService.Start(
		match,
		time.Now(),
		playerBoard,
		enemyBoard,
		playerEntityBoard,
		enemyEntityBoard,
		fleet,
		enemyFleet,
		enemyShipCells,
		playerShipCells,
	)

	if err != nil {
		return nil, err
	}

	return opponent, nil
}

// HandleCampaignResult processa o fim da partida
func (cs *CampaignService) HandleCampaignResult(username string, diff string, currentMatchResult *entity.MatchResult, playerWins, enemyWins int) (*entity.MatchResult, bool, error) {
	profile, err := FindProfile(username)
	if err != nil {
		return nil, false, err
	}

	if profile.CurrentCampaign == nil {
		return nil, false, fmt.Errorf("nenhuma campanha ativa para %s", username)
	}

	// 1. Acumulação de Estatísticas
	isFirstMatch := (playerWins + enemyWins) == 1
	var accumulated entity.MatchResult

	if isFirstMatch {
		accumulated = *currentMatchResult
		accumulated.Win = false // Série ainda não vencida
	} else {
		if prev, ok := profile.CurrentCampaign.DifficultyStep[diff]; ok {
			accumulated = prev
			accumulated.PlayerShots += currentMatchResult.PlayerShots
			accumulated.Hits += currentMatchResult.Hits
			accumulated.LostShips += currentMatchResult.LostShips
			accumulated.KilledShips += currentMatchResult.KilledShips
			accumulated.Duration += currentMatchResult.Duration
			accumulated.Score += currentMatchResult.Score
			if currentMatchResult.HigherHitSequence > accumulated.HigherHitSequence {
				accumulated.HigherHitSequence = currentMatchResult.HigherHitSequence
			}
		} else {
			accumulated = *currentMatchResult
		}
	}

	// Salva o estado intermediário (acumulado) no perfil
	profile.CurrentCampaign.DifficultyStep[diff] = accumulated
	if err := UpdateProfile(*profile); err != nil {
		return nil, false, err
	}

	// 2. Verifica se a série terminou
	isSeriesOver := playerWins >= 2 || enemyWins >= 2
	if !isSeriesOver {
		return nil, false, nil
	}

	// 3. Finalização da Série
	finalRes := accumulated
	finalRes.Win = playerWins >= 2
	finalRes.Mode = "Campanha"
	finalRes.Difficulty = diff

	// Atualiza o passo final com o status de vitória correto
	profile.CurrentCampaign.DifficultyStep[diff] = finalRes

	// 5. Persistência no histórico
	_, err = AddMatchToProfile(profile, finalRes)
	return &finalRes, true, err
}

// GetNextDifficulty retorna qual o próximo passo da campanha.
func (cs *CampaignService) GetNextDifficulty(c *entity.Campaign) (string, bool) {
	if c.DifficultyStep == nil {
		return "easy", false
	}
	steps := []string{"easy", "medium", "hard"}
	for _, step := range steps {
		res, ok := c.DifficultyStep[step]
		if !ok || !res.Win {
			return step, false
		}
	}
	return "", true
}

// selectAI instancia a IA correta para o nível.
func (cs *CampaignService) selectAI(diff string, fleet *entity.Fleet) *ai.AIPlayer {
	switch diff {
	case "medium":
		return ai.NewMediumAIPlayer(fleet)
	case "hard":
		return ai.NewHardAIPlayer(fleet)
	default:
		return ai.NewEasyAIPlayer()
	}
}
