package service

import (
	"time"

	"github.com/allanjose001/go-battleship/internal/ai"
	"github.com/allanjose001/go-battleship/internal/entity"
)

// BattleService define a interface para interação com a lógica de batalha.
// Essa abstração permite que a camada de apresentação (Scenes/Components) não conheça
// os detalhes internos de como a batalha é processada.
type BattleService interface {
	// HandlePlayerClick processa a tentativa de ataque do jogador em uma célula (linha, coluna).
	HandlePlayerClick(row, col int) (*entity.MatchResult, error)
	// HandleEnemyTurn executa o turno da IA (computador).
	HandleEnemyTurn() (*entity.MatchResult, error)
	// Stats retorna as estatísticas atuais da partida (tiros, acertos e de quem é a vez).
	Stats() (playerShots, playerHits, enemyShots, enemyHits int, isPlayerTurn bool)
	// WinnerName retorna o nome do vencedor caso a partida tenha terminado.
	WinnerName() string
}

// battleService é a implementação concreta da interface BattleService.
// Ele orquestra o fluxo do jogo usando o MatchService e mantém o estado da partida.
type battleService struct {
	// matchSvc é o serviço de domínio que aplica as regras do jogo.
	matchSvc *MatchService
	// match mantém o estado atual da partida (tabuleiros, turnos, pontuação).
	match *entity.Match
	// aiPlayer é a instância da inteligência artificial que joga contra o humano.
	aiPlayer *ai.AIPlayer
	// profile é o perfil do jogador humano, usado para registrar estatísticas de vitória/derrota.
	profile *entity.Profile
}

// NewBattleServiceFromMatch inicializa o serviço a partir de um Match existente no contexto.
// Se o Match ainda não foi inicializado (runtime), ele configura a IA e inicia o jogo.
func NewBattleServiceFromMatch(match *entity.Match) (BattleService, error) {
	// Cria o serviço de setup (para IA) e o serviço de partida (regras)
	setupSvc := NewBattleSetupService()
	matchSvc := NewMatchService(nil, 500*time.Millisecond)

	var aiPlayer *ai.AIPlayer

	// Se PlayerEntityBoard for nil, significa que a IA ainda não foi configurada para este Match.
	if match.PlayerEntityBoard == nil {
		// Inicializa a IA, criando o tabuleiro lógico e a frota da IA
		var entityBoard *entity.Board
		var fleet *entity.Fleet
		aiPlayer, entityBoard, fleet = setupSvc.InitBattleAI(match.PlayerShips)

		// Calcula o total de células ocupadas por navios do jogador (vida do player)
		totalCells := 0
		for _, ship := range match.PlayerShips {
			if ship != nil {
				totalCells += ship.Size
			}
		}

		// Inicia a partida (Start) para configurar os ponteiros runtime no Match e mudar status para InProgress
		// Note que passamos match.PlayerBoard e match.EnemyBoard que vieram do PlacementScene
		if err := matchSvc.Start(
			match,
			time.Now(),
			match.PlayerBoard,
			match.EnemyBoard,
			entityBoard,
			fleet,
			totalCells, // Vida do Player
			totalCells, // Vida da IA (simétrica)
		); err != nil {
			return nil, err
		}
	} else {
		// Se já existe PlayerEntityBoard, estamos retomando um Match (ex: persistência futura).
		// Precisamos recriar o AIPlayer com a frota existente.
		// OBS: Se a IA tiver estado interno complexo (memória de tiros), precisaria ser restaurado aqui.
		// Por enquanto, recriamos com a frota salva no Match.
		aiPlayer = ai.NewHardAIPlayer(match.PlayerFleet)
	}

	return &battleService{
		matchSvc: matchSvc,
		match:    match,
		aiPlayer: aiPlayer,
		profile:  match.Profile,
	}, nil
}

// HandlePlayerClick processa a interação do jogador ao clicar no tabuleiro inimigo.
func (s *battleService) HandlePlayerClick(row, col int) (*entity.MatchResult, error) {
	// Verifica se a partida foi inicializada corretamente.
	if s.matchSvc == nil || s.match == nil {
		return nil, ErrMatchNotReady
	}

	// Solicita ao serviço de domínio que processe o ataque do jogador.
	ev, err := s.matchSvc.PlayerAttack(s.match, time.Now(), row, col)
	if err != nil {
		return nil, err
	}

	// Se o ataque resultou em Game Over, processa o fim de jogo.
	if ev.GameOver {
		// Gera o resultado final da partida.
		res := s.matchSvc.ResultForPlayer(s.match)
		// Registra o resultado no perfil do jogador (se existir).
		if s.profile != nil {
			_, _ = AddMatchToProfile(s.profile, res)
		}
		return &res, nil
	}

	// Retorna nil se o jogo continua.
	return nil, nil
}

// HandleEnemyTurn executa a lógica de ataque da IA.
func (s *battleService) HandleEnemyTurn() (*entity.MatchResult, error) {
	// Verifica pré-condições.
	if s.matchSvc == nil || s.match == nil || s.aiPlayer == nil {
		return nil, ErrMatchNotReady
	}

	// Executa um passo da IA (pode não fazer nada se não for a vez dela ou se estiver em delay).
	ev, err := s.matchSvc.EnemyAttackStep(s.match, time.Now(), s.aiPlayer)
	if err != nil {
		// Ignora erros esperados que indicam que a IA ainda não deve agir.
		if err == ErrActionNotReady ||
			err == ErrNoEnemyAttackSched ||
			err == ErrMatchNotInProgress ||
			err == ErrMatchFinished {
			return nil, nil
		}
		return nil, err
	}

	// Se a IA venceu, processa o fim de jogo.
	if ev.GameOver {
		res := s.matchSvc.ResultForPlayer(s.match)
		if s.profile != nil {
			_, _ = AddMatchToProfile(s.profile, res)
		}
		return &res, nil
	}

	return nil, nil
}

// Stats retorna um resumo do estado atual da partida para exibição no HUD.
func (s *battleService) Stats() (playerShots, playerHits, enemyShots, enemyHits int, isPlayerTurn bool) {
	if s.match == nil {
		return 0, 0, 0, 0, true
	}

	// Extrai dados diretos da struct Match.
	return s.match.PlayerShots,
		s.match.PlayerHits,
		s.match.EnemyShots,
		s.match.EnemyHits,
		s.match.Turn == entity.TurnPlayer
}

// WinnerName retorna o nome de quem venceu a partida.
func (s *battleService) WinnerName() string {
	if s.matchSvc == nil || s.match == nil {
		return ""
	}

	// Obtém o resultado final.
	res := s.matchSvc.ResultForPlayer(s.match)
	if res.Win {
		// Se o jogador venceu, retorna seu nome de perfil ou um padrão.
		if s.profile != nil && s.profile.Username != "" {
			return s.profile.Username
		}
		return "Jogador"
	}

	// Caso contrário, a IA venceu.
	return "IA_MAR"
}
