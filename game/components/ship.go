// Package components contém widgets e helpers visuais reutilizáveis do jogo.
package components

import (
	"image/color"
	"math"

	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/game/shared/placement"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// DrawShip desenha um navio em três estados possíveis:
// - colocado no tabuleiro
// - sendo arrastado
// - parado na lista de seleção à direita
func DrawShip(screen *ebiten.Image, b *board.Board, ship *placement.ShipPlacement, active bool, orientation board.Orientation) {
	// Se o navio ou a imagem não existem, não há nada para desenhar.
	if ship == nil || ship.Image == nil {
		return
	}

	// Options de desenho usadas para escalar, rotacionar e transladar a imagem.
	op := &ebiten.DrawImageOptions{}

	// Caso 1: navio já está colocado no tabuleiro.
	if ship.Placed {
		// Cada célula do tabuleiro tem o mesmo tamanho em pixels.
		cellSize := b.Size / float64(board.Cols)

		// Converte posição em células (X,Y) para coordenadas de tela.
		x := b.X + float64(ship.X)*cellSize
		y := b.Y + float64(ship.Y)*cellSize

		// Pega largura/altura originais da imagem.
		iw, ih := ship.Image.Size()

		// Escala a imagem para ocupar o número de células do navio.
		op.GeoM.Scale(
			(cellSize*float64(ship.Size))/float64(iw),
			cellSize/float64(ih),
		)

		// Se o navio está vertical, rotaciona 90° e ajusta a origem.
		if ship.Orientation == board.Vertical {
			op.GeoM.Rotate(math.Pi / 2)
			op.GeoM.Translate(cellSize, 0)
		}

		// Move a imagem escalada para a posição calculada no tabuleiro.
		op.GeoM.Translate(x, y)

		// Se este navio é o ativo, desenha um retângulo de destaque em volta.
		if active {
			highlightColor := color.RGBA{255, 255, 0, 160}
			highlightW := cellSize * float64(ship.Size)
			highlightH := cellSize
			if ship.Orientation == board.Vertical {
				highlightW, highlightH = highlightH, highlightW
			}
			ebitenutil.DrawRect(screen, x-2, y-2, highlightW+4, highlightH+4, highlightColor)
		}

		// Finalmente desenha a imagem do navio.
		screen.DrawImage(ship.Image, op)
		return
	}

	// Caso 2: navio está sendo arrastado pelo jogador.
	if ship.Dragging {
		// Usa o mesmo tamanho de célula do tabuleiro para manter escala consistente.
		cellSize := b.Size / float64(board.Cols)
		iw, ih := ship.Image.Size()

		// Escala imagem proporcional ao tamanho do navio.
		op.GeoM.Scale(
			(cellSize*float64(ship.Size))/float64(iw),
			cellSize/float64(ih),
		)

		// Para navio sendo arrastado, usamos a orientação "corrente" do service.
		if orientation == board.Vertical {
			op.GeoM.Rotate(math.Pi / 2)
			op.GeoM.Translate(cellSize, 0)
		}

		// Posiciona a imagem na posição de drag (em pixels).
		op.GeoM.Translate(ship.DragX, ship.DragY)

		// Se está ativo, desenha highlight em volta da posição de drag.
		if active {
			highlightColor := color.RGBA{255, 255, 0, 160}
			highlightW := cellSize * float64(ship.Size)
			highlightH := cellSize
			if orientation == board.Vertical {
				highlightW, highlightH = highlightH, highlightW
			}
			ebitenutil.DrawRect(screen, ship.DragX-2, ship.DragY-2, highlightW+4, highlightH+4, highlightColor)
		}
	} else {
		// Caso 3: navio está parado na lista lateral (ainda não selecionado).
		// Apenas translada a imagem para a posição fixa da lista.
		op.GeoM.Translate(ship.ListX, ship.ListY)

		// Se é o navio ativo na lista, desenha uma borda branca em volta.
		if active {
			w, h := ship.Image.Size()
			highlightColor := colors.White
			ebitenutil.DrawRect(screen, ship.ListX-2, ship.ListY-2, float64(w)+4, float64(h)+4, highlightColor)
		}
	}
	// Em qualquer um dos casos acima, desenha a imagem do navio com as transformações calculadas.
	screen.DrawImage(ship.Image, op)
}
