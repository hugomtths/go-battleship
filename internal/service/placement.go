package service

import (
	"math"
	"math/rand"

	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/game/shared/placement"
)

// PlacementRenderer descreve qualquer tipo capaz de desenhar o tabuleiro
// de posicionamento e os navios, usado pela cena de placement.
type PlacementRenderer interface {
	Draw(b *board.Board, ships []*placement.ShipPlacement, active *placement.ShipPlacement, orientation board.Orientation)
}

// PlacementService define as operações de alto nível usadas pela cena
// de placement para manipular navios.
type PlacementService interface {
	AllShipsPlaced() bool
	RandomPlacement()
	Rotate()
	SelectOnBoard(mouseX, mouseY float64) bool
	SelectOnList(mouseX, mouseY float64) bool
	UpdateDragging(mouseX, mouseY float64)
	DropSelected() bool
	Draw(r PlacementRenderer)
	BoardRect() (x, y, size float64)
}

// placementService é a implementação concreta de PlacementService.
// Ela guarda o tabuleiro, a lista de navios e o estado de seleção/drag.
type placementService struct {
	board       *board.Board
	ships       []*placement.ShipPlacement
	selected    *placement.ShipPlacement
	activeShip  *placement.ShipPlacement
	orientation board.Orientation
}

// NewPlacementService cria um novo serviço de posicionamento com
// orientação inicial horizontal.
func NewPlacementService(b *board.Board, ships []*placement.ShipPlacement) PlacementService {
	return &placementService{
		board:       b,
		ships:       ships,
		orientation: board.Horizontal,
	}
}

// AllShipsPlaced retorna true quando todos os navios já foram
// posicionados no tabuleiro visual.
func (p *placementService) AllShipsPlaced() bool {
	for _, ship := range p.ships {
		if !ship.Placed {
			return false
		}
	}
	return true
}

// RandomPlacement limpa o tabuleiro visual e reposiciona todos os navios
// aleatoriamente, atualizando também os metadados de cada ShipPlacement.
func (p *placementService) RandomPlacement() {
	p.board.Clear()

	for _, ship := range p.ships {
		ship.Placed = false
	}

	var lastPlaced *placement.ShipPlacement

	for _, ship := range p.ships {
		for {
			row := rand.Intn(board.Rows)
			col := rand.Intn(board.Cols)
			or := board.Orientation(rand.Intn(2))

			if p.board.CanPlace(ship.Size, row, col, or) {
				p.board.PlaceShip(ship.Size, row, col, or)
				ship.Placed = true
				ship.X = col
				ship.Y = row
				ship.Orientation = or
				lastPlaced = ship
				break
			}
		}
	}

	p.activeShip = lastPlaced
}

// DropSelected tenta soltar o navio selecionado no tabuleiro,
// convertendo a posição do drag em linha/coluna. Retorna true em caso de sucesso.
func (p *placementService) DropSelected() bool {
	if p.selected == nil {
		return false
	}

	ship := p.selected
	ship.Dragging = false

	cellSize := p.board.Size / float64(board.Cols)

	targetX := ship.DragX
	targetY := ship.DragY

	// Verifica se a área do navio arrastado intersecta o tabuleiro
	if targetX+cellSize*float64(ship.Size) > p.board.X && targetX < p.board.X+p.board.Size &&
		targetY+cellSize*float64(ship.Size) > p.board.Y && targetY < p.board.Y+p.board.Size {

		// Converte a posição de drag para índices de célula
		col := int(math.Round((targetX - p.board.X) / cellSize))
		row := int(math.Round((targetY - p.board.Y) / cellSize))

		if p.board.CanPlace(ship.Size, row, col, p.orientation) {
			p.board.PlaceShip(ship.Size, row, col, p.orientation)
			ship.Placed = true
			ship.X = col
			ship.Y = row
			ship.Orientation = p.orientation
			p.activeShip = ship
			return true
		}
	}

	// Se não conseguir colocar, o navio volta para a lista
	ship.Placed = false
	return false
}

// removeShipFromBoard limpa o tabuleiro e recoloca todos os navios
// exceto o alvo, usado principalmente durante a rotação.
func (p *placementService) removeShipFromBoard(target *placement.ShipPlacement) {
	if target == nil {
		return
	}

	target.Placed = false
	p.board.Clear()

	for _, ship := range p.ships {
		if ship != target && ship.Placed {
			p.board.PlaceShip(ship.Size, ship.Y, ship.X, ship.Orientation)
		}
	}
}

// Rotate rotaciona o navio ativo no tabuleiro, se existir,
// ou apenas alterna a orientação padrão dos próximos navios.
func (p *placementService) Rotate() {
	if p.activeShip != nil && p.activeShip.Placed {
		p.rotateActive()
	} else {
		p.toggleOrientation()
	}
}

// rotateActive tenta rotacionar o navio ativo mantendo a posição de origem.
// Caso não seja possível, restaura a orientação anterior.
func (p *placementService) rotateActive() bool {
	ship := p.activeShip
	if ship == nil || !ship.Placed {
		return false
	}

	newOri := board.Horizontal
	if ship.Orientation == board.Horizontal {
		newOri = board.Vertical
	}

	if ship.Orientation == newOri {
		return true
	}

	p.removeShipFromBoard(ship)

	if p.board.CanPlace(ship.Size, ship.Y, ship.X, newOri) {
		p.board.PlaceShip(ship.Size, ship.Y, ship.X, newOri)
		ship.Orientation = newOri
		ship.Placed = true
		p.orientation = newOri
		return true
	}

	p.board.PlaceShip(ship.Size, ship.Y, ship.X, ship.Orientation)
	ship.Placed = true
	return false
}

// toggleOrientation alterna a orientação padrão usada para novos placements.
func (p *placementService) toggleOrientation() {
	if p.orientation == board.Horizontal {
		p.orientation = board.Vertical
	} else {
		p.orientation = board.Horizontal
	}
}

// SelectOnBoard tenta selecionar um navio já colocado no tabuleiro
// para iniciar um drag a partir da posição real dele.
func (p *placementService) SelectOnBoard(mouseX, mouseY float64) bool {
	for _, ship := range p.ships {
		if !ship.Placed {
			continue
		}

		cellSize := p.board.Size / float64(board.Cols)
		x := p.board.X + float64(ship.X)*cellSize
		y := p.board.Y + float64(ship.Y)*cellSize

		w := cellSize * float64(ship.Size)
		h := cellSize
		if ship.Orientation == board.Vertical {
			w, h = h, w
		}

		// Verifica se o clique caiu dentro da área ocupada pelo navio
		if mouseX >= x && mouseX <= x+w && mouseY >= y && mouseY <= y+h {
			ship.Dragging = true
			ship.DragX = x
			ship.DragY = y
			ship.OffsetX = mouseX - x
			ship.OffsetY = mouseY - y
			p.selected = ship
			p.activeShip = ship
			p.orientation = ship.Orientation
			return true
		}
	}

	return false
}

// SelectOnList tenta selecionar um navio ainda não colocado na lista
// lateral, iniciando um drag a partir da posição fixa da lista.
func (p *placementService) SelectOnList(mouseX, mouseY float64) bool {
	for _, ship := range p.ships {
		if ship.Placed {
			continue
		}

		w, h := ship.Image.Size()

		if mouseX >= ship.ListX && mouseX <= ship.ListX+float64(w) &&
			mouseY >= ship.ListY && mouseY <= ship.ListY+float64(h) {

			ship.Dragging = true
			ship.DragX = ship.ListX
			ship.DragY = ship.ListY
			ship.OffsetX = mouseX - ship.ListX
			ship.OffsetY = mouseY - ship.ListY
			p.selected = ship
			p.activeShip = ship
			return true
		}
	}

	return false
}

// UpdateDragging atualiza a posição de drag do navio selecionado
// com base na posição atual do mouse.
func (p *placementService) UpdateDragging(mouseX, mouseY float64) {
	ship := p.selected
	if ship == nil || !ship.Dragging {
		return
	}

	ship.DragX = mouseX - ship.OffsetX
	ship.DragY = mouseY - ship.OffsetY
}

// Draw delega o desenho do tabuleiro e navios para o renderer injetado.
func (p *placementService) Draw(r PlacementRenderer) {
	if r == nil {
		return
	}
	r.Draw(p.board, p.ships, p.activeShip, p.orientation)
}

// BoardRect retorna a posição e o tamanho do tabuleiro,
// útil para alinhar elementos visuais na cena.
func (p *placementService) BoardRect() (x, y, size float64) {
	return p.board.X, p.board.Y, p.board.Size
}
