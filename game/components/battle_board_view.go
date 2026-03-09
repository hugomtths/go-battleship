package components

import (
	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/game/shared/placement"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
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
func (v *BattleBoardView) DrawBoard(screen *ebiten.Image, b *board.Board, ships []*placement.ShipPlacement, fleet *entity.Fleet, entityBoard *entity.Board, hideUndestroyed bool) {
	// Se o tabuleiro for nil, não há nada para desenhar.
	if b == nil {
		return
	}

	// 1. Desenha navios, se houver uma lista de navios fornecida.
	// Isso geralmente é usado apenas para o tabuleiro do próprio jogador.
	if ships != nil && fleet != nil {
		for i, ship := range ships {
			// Só desenha se o navio estiver marcado como posicionado.
			if ship.Placed {
				// Verifica se o navio correspondente na frota lógica está afundado.
				isSunk := false

				// 1. Tenta via mapeamento espacial (entityBoard)
				if entityBoard != nil {
					// Usa a posição do navio para encontrar a referência lógica correta
					if ship.Y >= 0 && ship.Y < len(entityBoard.Positions) && ship.X >= 0 && ship.X < len(entityBoard.Positions[0]) {
						pos := entityBoard.Positions[ship.Y][ship.X]
						entShip := entity.GetShipReference(pos)
						if entShip != nil && entShip.IsDestroyed() {
							isSunk = true
						}
					}
				}

				// 2. Fallback: Se não encontrou via espaço (ou entityBoard nulo), tenta via índice na frota.
				// Isso garante que se o mapeamento espacial falhar, ainda temos a correspondência por ordem/tamanho.
				if !isSunk && fleet != nil && i < len(fleet.Ships) {
					logicalShip := fleet.Ships[i]
					// Só considera se o tamanho bater (proteção contra desincronia de índices)
					if logicalShip.Size == ship.Size && logicalShip.IsDestroyed() {
						isSunk = true
					}
				}

				if isSunk {
					// Se estiver afundado, usa a imagem de navio afundado.
					originalImage := ship.Image
					ship.Image = ship.SunkImage
					DrawShip(screen, b, ship, false, ship.Orientation)
					ship.Image = originalImage // Restaura para não afetar outros estados
				} else if !hideUndestroyed {
					// Caso contrário, usa a imagem normal.
					DrawShip(screen, b, ship, false, ship.Orientation)
				}
			}
		}
	}

	// 2. Desenha marcadores (hits/misses) sobre o tabuleiro.
	// Isso acontece para ambos os jogadores (mostra onde já atiraram).
	v.DrawMarkers(screen, b, entityBoard)
}

// DrawMarkers itera sobre todas as células do tabuleiro e desenha os indicadores de tiro.
func (v *BattleBoardView) DrawMarkers(screen *ebiten.Image, b *board.Board, entityBoard *entity.Board) {
	if b == nil {
		return
	}

	cellSize := b.Size / float64(board.Cols)
	// reutiliza uma única instância de DrawImageOptions — zero alocações por frame
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterNearest

	for i := 0; i < board.Rows; i++ {
		for j := 0; j < board.Cols; j++ {
			cell := b.Cells[i][j]

			if cell.State != board.Hit && cell.State != board.Miss {
				continue
			}

			// Se a célula for um Hit, verifica se o navio naquela posição está afundado.
			// Se estiver afundado, não desenhamos o fogo (fireAnimation/hitImage).
			if cell.State == board.Hit && entityBoard != nil {
				pos := entityBoard.Positions[i][j]
				ship := entity.GetShipReference(pos)
				if ship != nil && ship.IsDestroyed() {
					continue // Pula o desenho do fogo se o navio estiver afundado
				}
			}

			x := b.X + float64(j)*cellSize
			y := b.Y + float64(i)*cellSize

			if cell.State == board.Hit {
				var img *ebiten.Image
				if v.fireAnimation != nil {
					img = v.fireAnimation.CurrentFrame()
				}
				if img == nil {
					img = v.hitImage
				}
				if img != nil {
					op.GeoM.Reset()
					iw, ih := img.Size()
					op.GeoM.Scale(cellSize/float64(iw), cellSize/float64(ih))
					op.GeoM.Translate(x, y)
					screen.DrawImage(img, op)
				}
			} else if cell.State == board.Miss {
				if v.missImage != nil {
					op.GeoM.Reset()
					iw, ih := v.missImage.Size()
					op.GeoM.Scale(cellSize/float64(iw), cellSize/float64(ih))
					op.GeoM.Translate(x, y)
					screen.DrawImage(v.missImage, op)
				}
			}
		}
	}
}
