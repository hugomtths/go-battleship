package components

import (
	"fmt"
	"image/color"

	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/internal/service"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// BattleSide define qual lado da batalha o HUD está representando.
type BattleSide int

const (
	// SidePlayer indica o lado do jogador humano.
	SidePlayer BattleSide = iota
	// SideAI indica o lado do adversário (IA).
	SideAI
)

// BattleHUD é o componente responsável por exibir as informações de status da batalha.
// Ele gerencia os labels de nome, tentativas e acertos, e se atualiza consultando o BattleService.
type BattleHUD struct {
	// nameLabel exibe o nome do jogador ou IA.
	nameLabel *Text
	// attemptsLabel exibe o número de tiros disparados.
	attemptsLabel *Text
	// hitsLabel exibe o número de acertos confirmados.
	hitsLabel *Text
	// battleSvc é a referência ao serviço de batalha para consultar estatísticas em tempo real.
	battleSvc service.BattleService
	// side indica se este HUD pertence ao jogador ou à IA, para filtrar as estatísticas corretas.
	side BattleSide
}

// NewBattleHUD cria uma nova instância do HUD de batalha.
// Recebe os componentes de texto pré-configurados e as dependências necessárias.
func NewBattleHUD(nameLabel, attemptsLabel, hitsLabel *Text, svc service.BattleService, side BattleSide) *BattleHUD {
	return &BattleHUD{
		nameLabel:     nameLabel,
		attemptsLabel: attemptsLabel,
		hitsLabel:     hitsLabel,
		battleSvc:     svc,
		side:          side,
	}
}

// Update propaga a atualização para os componentes de texto filhos.
// Isso permite que animações ou efeitos nos textos continuem funcionando.
func (h *BattleHUD) Update(offset basic.Point) {
	if h.nameLabel != nil {
		h.nameLabel.Update(offset)
	}
	if h.attemptsLabel != nil {
		h.attemptsLabel.Update(offset)
	}
	if h.hitsLabel != nil {
		h.hitsLabel.Update(offset)
	}
}

// Draw renderiza o HUD na tela.
// Ele consulta as estatísticas atuais do serviço e atualiza os textos antes de desenhar.
func (h *BattleHUD) Draw(screen *ebiten.Image, b *board.Board) {
	// Se não houver serviço, não há dados para mostrar.
	if h.battleSvc == nil {
		return
	}

	// Consulta as estatísticas mais recentes da batalha.
	pShots, pHits, aiShots, aiHits, isPlayerTurn := h.battleSvc.Stats()

	attempts := 0
	hits := 0
	isTurn := false

	// Seleciona as estatísticas relevantes com base no lado configurado para este HUD.
	if h.side == SidePlayer {
		attempts = pShots
		hits = pHits
		isTurn = isPlayerTurn // É a vez do jogador?
	} else {
		attempts = aiShots
		hits = aiHits
		isTurn = !isPlayerTurn // Se não é a vez do jogador, é a vez da IA.
	}

	// Define a posição base para desenhar os indicadores, logo abaixo do tabuleiro.
	baseX := b.X
	baseY := b.Y + b.Size + 20

	// Define a cor do indicador de turno (Verde se for a vez deste lado, Vermelho caso contrário).
	indicatorColor := color.RGBA{255, 0, 0, 255}
	if isTurn {
		indicatorColor = color.RGBA{0, 255, 0, 255}
	}

	// Desenha um pequeno quadrado indicando de quem é a vez.
	ebitenutil.DrawRect(screen, baseX, baseY, 20, 20, indicatorColor)

	// Desenha o nome do jogador/IA.
	h.nameLabel.Draw(screen)

	// Atualiza e desenha o contador de tentativas.
	h.attemptsLabel.Text = fmt.Sprintf("Tentativa: %d", attempts)
	h.attemptsLabel.Draw(screen)

	// Atualiza e desenha o contador de acertos.
	h.hitsLabel.Text = fmt.Sprintf("Acertos: %d", hits)
	h.hitsLabel.Draw(screen)
}
