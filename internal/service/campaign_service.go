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
	match := cs.matchService.Create(profile.CurrentCampaign.ID)

	// Inicia os estados da partida (Turnos, Timers, etc)
	err = cs.matchService.Start(
		match,
		time.Now(),
		playerBoard,
		enemyBoard,
		playerEntityBoard,
		fleet,
		enemyShipCells,
		playerShipCells,
	)

	if err != nil {
		return nil, err
	}

	return opponent, nil
}

// HandleCampaignResult processa o fim da partida
func (cs *CampaignService) HandleCampaignResult(username string, diff string, m *entity.Match) error {
	profile, err := FindProfile(username)
	if err != nil {
		return err
	}

	if profile.CurrentCampaign == nil {
		return fmt.Errorf("nenhuma campanha ativa para %s", username)
	}

	// 1. Extrai o resultado final da partida
	result := m.Result()

	// 2. Atualiza o progresso da campanha em memória
	profile.CurrentCampaign.DifficultyStep[diff] = result

	// 3. Lógica de progressão: Se venceu o Hard, encerra a campanha ativa
	_, finished := cs.GetNextDifficulty(profile.CurrentCampaign)
	if finished && result.Win {
		profile.CurrentCampaign.IsActive = false
		profile.Campaigns = append(profile.Campaigns, *profile.CurrentCampaign)
		profile.CurrentCampaign = nil
	}

	// 4. Persistência com profile_scene
	_, err = AddMatchToProfile(profile, result)
	return err
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
