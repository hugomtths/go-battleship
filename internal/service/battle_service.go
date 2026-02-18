package service

import (
	"time"

	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/game/shared/placement"
	"github.com/allanjose001/go-battleship/game/state"
	"github.com/allanjose001/go-battleship/internal/ai"
	"github.com/allanjose001/go-battleship/internal/entity"
)

// BattleService concentra toda a lógica da fase de batalha:
// controle de turno, cliques do jogador, turno da IA e estatísticas.
type BattleService struct {
	// state guarda os tabuleiros e dados globais do jogo
	state       *state.GameState
	attack      *AttackService
	battleSetup *BattleSetupService
	// aiPlayer é a inteligência artificial que decide os tiros da IA
	aiPlayer *ai.AIPlayer
	// entityBoard e entityFleet representam a visão interna da IA
	entityBoard *entity.Board
	entityFleet *entity.Fleet

	playerShips []*placement.ShipPlacement

	totalShipCells int

	// contadores de tentativas e acertos do jogador e da IA
	playerAttempts int
	playerHits     int
	aiAttempts     int
	aiHits         int

	// indica de quem é o turno atual
	isPlayerTurn bool

	// controle de atraso entre o tiro do jogador e o da IA
	aiTurnPending bool
	aiTurnAt      time.Time

	// vencedor atual ("", "Jogador 1" ou "Jogador 2")
	winner string
}

// NewBattleService cria um serviço de batalha baseado em um GameState já existente.
// Se o GameService não for passado, é criada uma instância padrão.
func NewBattleService(gs *state.GameState, game *GameService, playerShips []*placement.ShipPlacement) *BattleService {
	b := &BattleService{
		state:        gs,
		playerShips:  playerShips,
		attack:       NewAttackService(),
		battleSetup:  NewBattleSetupService(),
		isPlayerTurn: true,
	}

	b.totalShipCells = calculateTotalCells(playerShips)

	b.aiPlayer, b.entityBoard, b.entityFleet = b.battleSetup.InitBattleAI(playerShips)

	return b
}

// HandlePlayerClick trata o clique do jogador no tabuleiro da IA.
// Converte a posição do mouse em célula, aplica o ataque e atualiza o turno.
// Retorna o nome do vencedor (caso a jogada finalize a partida).
func (b *BattleService) HandlePlayerClick(mouseX, mouseY float64) string {
	if !b.isPlayerTurn || b.winner != "" {
		return b.winner
	}

	aiBoard := b.state.AIBoard

	// Ignora cliques fora da área do tabuleiro da IA
	if mouseX < aiBoard.X || mouseX > aiBoard.X+aiBoard.Size || mouseY < aiBoard.Y || mouseY > aiBoard.Y+aiBoard.Size {
		return b.winner
	}

	// Traduz coordenadas de tela para linha/coluna do tabuleiro
	cellSize := aiBoard.Size / float64(board.Cols)
	col := int((mouseX - aiBoard.X) / cellSize)
	row := int((mouseY - aiBoard.Y) / cellSize)

	if col < 0 || col >= board.Cols || row < 0 || row >= board.Rows {
		return b.winner
	}

	var gameOver bool

	b.playerAttempts, b.playerHits, _, gameOver = b.attack.PlayerAttack(aiBoard, row, col, b.playerAttempts, b.playerHits, b.totalShipCells)

	if gameOver {
		b.winner = "Jogador 1"
		return b.winner
	}

	// Passa o turno para a IA com um pequeno atraso visual
	b.isPlayerTurn = false
	b.aiTurnPending = true
	b.aiTurnAt = time.Now().Add(500 * time.Millisecond)

	return b.winner
}

// Update processa o turno da IA quando for a vez dela jogar.
// Também verifica se a partida terminou.
func (b *BattleService) Update() string {
	if b.winner != "" {
		return b.winner
	}

	if b.aiTurnPending && time.Now().After(b.aiTurnAt) {
		b.aiTurnPending = false

		// Se algo falhar na inicialização da IA, devolve o turno ao jogador
		if b.aiPlayer == nil || b.entityBoard == nil {
			b.isPlayerTurn = true
			return b.winner
		}

		var gameOver bool

		b.aiAttempts, b.aiHits, gameOver = b.attack.AITurn(b.aiPlayer, b.entityBoard, b.state.PlayerBoard, b.aiAttempts, b.aiHits, b.totalShipCells)

		if gameOver {
			b.winner = "Jogador 2"
			return b.winner
		}

		// Terminado o turno da IA, devolve o turno ao jogador
		b.isPlayerTurn = true
	}

	return b.winner
}

// Stats retorna estatísticas básicas da partida e se é turno do jogador.
func (b *BattleService) Stats() (int, int, int, int, bool) {
	return b.playerAttempts, b.playerHits, b.aiAttempts, b.aiHits, b.isPlayerTurn
}

// PlayerBoard retorna o tabuleiro do jogador.
func (b *BattleService) PlayerBoard() *board.Board {
	return b.state.PlayerBoard
}

// AIBoard retorna o tabuleiro da IA.
func (b *BattleService) AIBoard() *board.Board {
	return b.state.AIBoard
}

func (b *BattleService) PlayerShips() []*placement.ShipPlacement {
	return b.playerShips
}
func calculateTotalCells(ships []*placement.ShipPlacement) int {
	total := 0
	for _, s := range ships {
		if s != nil {
			total += s.Size
		}
	}
	return total
}
