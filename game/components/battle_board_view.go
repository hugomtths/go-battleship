package components

import (
	"image/color"

	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/game/shared/placement"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// BattleBoardView é o componente responsável pela visualização do tabuleiro durante a batalha.
// Ele encapsula a lógica de desenhar os navios (para o jogador) e os marcadores de tiro (Hit/Miss).
type BattleBoardView struct {
	// fireAnimation controla a animação de fogo para acertos (Hits).
	fireAnimation *FireAnimation
	// hitImage é a imagem estática usada para acertos caso não haja animação.
	hitImage *ebiten.Image
	// missImage é a imagem usada para tiros na água (Miss).
	missImage *ebiten.Image
}

// NewBattleBoardView cria uma nova instância do visualizador de tabuleiro.
// Recebe os assets gráficos necessários para desenhar os estados das células.
func NewBattleBoardView(fireAnimation *FireAnimation, hitImage, missImage *ebiten.Image) *BattleBoardView {
	return &BattleBoardView{
		fireAnimation: fireAnimation,
		hitImage:      hitImage,
		missImage:     missImage,
	}
}

// DrawBoard é o método principal de desenho.
// Ele orquestra o desenho dos navios e dos marcadores sobre o tabuleiro.
// Parâmetros:
// - screen: A imagem de destino onde o desenho será feito.
// - b: O tabuleiro lógico (dados das células).
// - ships: A lista de navios a serem desenhados (pode ser nil para o tabuleiro do inimigo, onde não vemos os navios).
func (v *BattleBoardView) DrawBoard(screen *ebiten.Image, b *board.Board, ships []*placement.ShipPlacement) {
	// Se o tabuleiro for nil, não há nada para desenhar.
	if b == nil {
		return
	}

	// 1. Desenha navios, se houver uma lista de navios fornecida.
	// Isso geralmente é usado apenas para o tabuleiro do próprio jogador.
	if ships != nil {
		for _, ship := range ships {
			// Só desenha se o navio estiver marcado como posicionado.
			if ship.Placed {
				// DrawShip é uma função auxiliar do pacote components que desenha um único navio.
				DrawShip(screen, b, ship, false, ship.Orientation)
			}
		}
	}

	// 2. Desenha marcadores (hits/misses) sobre o tabuleiro.
	// Isso acontece para ambos os jogadores (mostra onde já atiraram).
	v.DrawMarkers(screen, b)
}

// DrawMarkers itera sobre todas as células do tabuleiro e desenha os indicadores de tiro.
func (v *BattleBoardView) DrawMarkers(screen *ebiten.Image, b *board.Board) {
	if b == nil {
		return
	}

	// Calcula o tamanho visual de cada célula baseado no tamanho total do tabuleiro e número de colunas.
	cellSize := b.Size / float64(board.Cols)

	// Percorre todas as linhas e colunas do grid.
	for i := 0; i < board.Rows; i++ {
		for j := 0; j < board.Cols; j++ {
			cell := b.Cells[i][j]

			// Verifica se a célula tem um estado de tiro (Hit ou Miss).
			if cell.State == board.Hit || cell.State == board.Miss {
				// Calcula a posição X e Y exata para desenhar o marcador nesta célula.
				x := b.X + float64(j)*cellSize
				y := b.Y + float64(i)*cellSize

				// Lógica para desenhar um Acerto (Hit).
				if cell.State == board.Hit {
					var img *ebiten.Image

					// Tenta pegar o quadro atual da animação de fogo.
					if v.fireAnimation != nil {
						img = v.fireAnimation.CurrentFrame()
					}
					// Se não houver animação ou falhar, usa a imagem estática.
					if img == nil {
						img = v.hitImage
					}

					if img != nil {
						// Configura as opções de desenho (escala e translação).
						op := &ebiten.DrawImageOptions{}
						iw, ih := img.Size()
						// Ajusta a escala da imagem para caber na célula.
						op.GeoM.Scale(cellSize/float64(iw), cellSize/float64(ih))
						// Move a imagem para a posição correta.
						op.GeoM.Translate(x, y)
						screen.DrawImage(img, op)
					} else {
						// Fallback: se não houver imagem nenhuma, desenha um quadrado vermelho.
						ebitenutil.DrawRect(screen, x+cellSize*0.25, y+cellSize*0.25, cellSize*0.5, cellSize*0.5, color.RGBA{255, 0, 0, 150})
					}
				} else if cell.State == board.Miss && v.missImage != nil {
					// Lógica para desenhar um Erro (Miss) com imagem.
					op := &ebiten.DrawImageOptions{}
					iw, ih := v.missImage.Size()
					op.GeoM.Scale(cellSize/float64(iw), cellSize/float64(ih))
					op.GeoM.Translate(x, y)
					screen.DrawImage(v.missImage, op)
				} else {
					// Fallback genérico ou para Miss sem imagem: desenha um quadrado cinza (ou vermelho se for hit e caiu aqui por algum motivo).
					c := color.RGBA{200, 200, 200, 150} // Cor padrão para Miss (Cinza)
					if cell.State == board.Hit {
						c = color.RGBA{255, 0, 0, 150} // Cor para Hit (Vermelho)
					}
					// Desenha um pequeno quadrado centralizado na célula.
					ebitenutil.DrawRect(screen, x+cellSize*0.25, y+cellSize*0.25, cellSize*0.5, cellSize*0.5, c)
				}
			}
		}
	}
}
