package components

import (
	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// BattleInput gerencia a entrada do usuário na cena de batalha.
// Ele detecta cliques do mouse e converte coordenadas de tela para índices de célula do tabuleiro.
type BattleInput struct {
	// board é a referência ao tabuleiro onde os cliques serão verificados.
	board *board.Board
}

// NewBattleInput cria uma nova instância do controlador de entrada de batalha.
func NewBattleInput(b *board.Board) *BattleInput {
	return &BattleInput{board: b}
}

// ClickedCell verifica se houve um clique válido no tabuleiro neste frame.
// Retorna a linha e coluna clicada, e um booleano indicando se o clique foi válido.
func (i *BattleInput) ClickedCell() (row, col int, ok bool) {
	// Se não houver tabuleiro, não há como calcular cliques.
	if i.board == nil {
		return 0, 0, false
	}

	// Verifica se o botão esquerdo do mouse foi pressionado neste exato frame.
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return 0, 0, false
	}

	// Obtém a posição atual do cursor na tela.
	mx, my := ebiten.CursorPosition()
	mouseX, mouseY := float64(mx), float64(my)

	// Verifica se o clique ocorreu dentro dos limites do retângulo do tabuleiro.
	if mouseX < i.board.X || mouseX > i.board.X+i.board.Size ||
		mouseY < i.board.Y || mouseY > i.board.Y+i.board.Size {
		return 0, 0, false
	}

	// Calcula o tamanho de cada célula para determinar a linha e coluna.
	cellSize := i.board.Size / float64(board.Cols)
	col = int((mouseX - i.board.X) / cellSize)
	row = int((mouseY - i.board.Y) / cellSize)

	// Validação final para garantir que os índices estão dentro da matriz do tabuleiro.
	if col < 0 || col >= board.Cols || row < 0 || row >= board.Rows {
		return 0, 0, false
	}

	// Retorna os índices da célula clicada e true indicando sucesso.
	return row, col, true
}
