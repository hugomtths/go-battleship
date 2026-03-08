package service

import (
	"math/rand"

	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/game/shared/placement"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

// DynamicBattleService define a interface para lógica específica do modo dinâmico (movimento de navios).
// Esta interface encapsula todas as regras de negócio relacionadas à movimentação das embarcações
// durante a fase de batalha dinâmica.
type DynamicBattleService interface {
	// HandlePlayerMove processa o movimento do navio do jogador.
	// Recebe a posição (linha, coluna) de uma parte do navio e a direção desejada.
	// Retorna erro se o movimento for inválido (fora do mapa, colisão, etc).
	HandlePlayerMove(row, col int, dir entity.Direction) error

	// HandleEnemyMove executa a lógica de movimento dos navios do inimigo (IA).
	// A IA tenta mover seus navios aleatoriamente para dificultar o jogo.
	HandleEnemyMove() error
}

// dynamicBattleService é a implementação concreta da interface DynamicBattleService.
// Mantém uma referência à partida atual (Match) para manipular os tabuleiros.
type dynamicBattleService struct {
	match *entity.Match
}

// NewDynamicBattleService cria uma nova instância do serviço de batalha dinâmica.
// Recebe o ponteiro para a struct Match que contém o estado atual do jogo.
func NewDynamicBattleService(match *entity.Match) DynamicBattleService {
	return &dynamicBattleService{
		match: match,
	}
}

// HandlePlayerMove implementa a lógica de movimento do jogador.
// 1. Valida se a partida está pronta.
// 2. Tenta mover o navio na matriz lógica (EntityBoard).
// 3. Sincroniza o estado visual (células do tabuleiro) com o estado lógico.
// 4. Reconstrói a lista de navios (PlayerShips) para atualização visual correta.
func (s *dynamicBattleService) HandlePlayerMove(row, col int, dir entity.Direction) error {
	// Verificação de segurança: garante que os dados necessários existem.
	if s.match == nil || s.match.PlayerEntityBoard == nil || s.match.PlayerBoard == nil {
		return entity.ErrMatchNotReady
	}

	// Tenta realizar o movimento no tabuleiro lógico (onde as entidades/navios realmente existem).
	// O método MoveShipSegment cuida da validação de limites e colisões.
	_, err := s.match.PlayerEntityBoard.MoveShipSegment(row, col, dir)
	if err != nil {
		return err // Retorna erro se o movimento for inválido.
	}

	// Sincronização: Atualiza o estado visual de cada célula do tabuleiro do jogador
	// baseando-se na nova posição dos navios no tabuleiro lógico.
	for r := 0; r < board.Rows; r++ {
		for c := 0; c < board.Cols; c++ {
			// Obtém a referência do navio na posição (r, c) do tabuleiro lógico.
			entShip := entity.GetShipReference(s.match.PlayerEntityBoard.Positions[r][c])
			// Obtém a célula visual correspondente.
			cell := &s.match.PlayerBoard.Cells[r][c]

			// Apenas atualiza células que NÃO foram atingidas (Hit ou Miss).
			// Isso preserva o histórico de tiros no tabuleiro.
			if cell.State != board.Hit && cell.State != board.Miss {
				if entShip != nil {
					// Se há um navio aqui logicamente, define visualmente como Ship.
					cell.State = board.Ship
				} else {
					// Se não há navio, define como Empty (água).
					cell.State = board.Empty
				}
			}
		}
	}

	// Reconstrói a lista de objetos ShipPlacement usada para desenhar os sprites dos navios.
	// Isso é necessário porque a posição e orientação dos navios mudaram.
	s.rebuildPlayerShips()

	return nil
}

// rebuildPlayerShips reconstrói a lista match.PlayerShips baseada no estado atual do tabuleiro lógico.
// É essencial para que a renderização visual (sprites) acompanhe a lógica do jogo após um movimento.
func (s *dynamicBattleService) rebuildPlayerShips() {
	// Passo 1: Identificar todos os navios únicos presentes no tabuleiro.
	uniqueShips := make(map[*entity.Ship]bool)
	for r := 0; r < board.Rows; r++ {
		for c := 0; c < board.Cols; c++ {
			ship := entity.GetShipReference(s.match.PlayerEntityBoard.Positions[r][c])
			if ship != nil {
				uniqueShips[ship] = true
			}
		}
	}

	// Passo 2: Mapear tamanhos de navios para suas imagens originais.
	// Isso preserva a arte original (sprites) dos navios ao invés de usar placeholders.
	sizeToImage := make(map[int]*ebiten.Image)
	var availableImages []*ebiten.Image
	seenImages := make(map[*ebiten.Image]bool)

	if s.match.PlayerShips != nil {
		for _, sp := range s.match.PlayerShips {
			if sp.Image != nil {
				sizeToImage[sp.Size] = sp.Image
				// Guarda imagens únicas disponíveis para uso como fallback.
				if !seenImages[sp.Image] {
					availableImages = append(availableImages, sp.Image)
					seenImages[sp.Image] = true
				}
			}
		}
	}

	// Passo 3: Criar novos ShipPlacements para cada navio encontrado.
	var newPlacements []*placement.ShipPlacement
	for ship := range uniqueShips {
		// Encontrar a coordenada superior esquerda (top-left) do navio no grid.
		minR, minC := board.Rows, board.Cols
		found := false
		for r := 0; r < board.Rows; r++ {
			for c := 0; c < board.Cols; c++ {
				if entity.GetShipReference(s.match.PlayerEntityBoard.Positions[r][c]) == ship {
					if r < minR {
						minR = r
					}
					if c < minC {
						minC = c
					}
					found = true
				}
			}
		}

		if !found {
			continue // Navio não encontrado no grid (improvável se veio de uniqueShips).
		}

		// Determina a orientação visual baseada na propriedade Horizontal do navio.
		orientation := board.Vertical
		if ship.Horizontal {
			orientation = board.Horizontal
		}

		// Determina qual imagem usar para este navio.
		img := sizeToImage[ship.Size]
		if img == nil {
			// Fallback: Se não houver imagem específica para o tamanho, usa a maior disponível
			// para garantir qualidade ao redimensionar (se necessário).
			if len(availableImages) > 0 {
				img = availableImages[0]
				maxW, _ := img.Size()
				for _, cand := range availableImages {
					w, _ := cand.Size()
					if w > maxW {
						img = cand
						maxW = w
					}
				}
			}
		}

		// Adiciona o novo posicionamento à lista.
		newPlacements = append(newPlacements, &placement.ShipPlacement{
			Placed:      true,
			X:           minC,
			Y:           minR,
			Size:        ship.Size,
			Orientation: orientation,
			Image:       img,
		})
	}
	// Atualiza a lista oficial de navios do jogador na partida.
	s.match.PlayerShips = newPlacements
}

// HandleEnemyMove implementa a lógica de movimento da IA (Inimigo).
// Como a IA não mantém um EntityBoard persistente da mesma forma que o jogador (para simplificação),
// este método precisa reconstruir o estado lógico a partir do tabuleiro visual antes de mover.
func (s *dynamicBattleService) HandleEnemyMove() error {
	// 1. Reconstrói o tabuleiro lógico (EntityBoard) a partir do estado visual atual do inimigo.
	// Isso agrupa células adjacentes de 'Ship' ou 'Hit' em objetos navio.
	entityBoard := s.reconstructEntityBoard(s.match.EnemyBoard)

	// 2. Identifica todos os navios presentes no tabuleiro reconstruído.
	type shipRef struct {
		ship *entity.Ship
		r, c int
	}
	var ships []shipRef
	visited := make(map[*entity.Ship]bool)

	for r := 0; r < 10; r++ {
		for c := 0; c < 10; c++ {
			ship := entity.GetShipReference(entityBoard.Positions[r][c])
			if ship != nil && !visited[ship] {
				visited[ship] = true
				// Guarda uma referência (coordenada) para tentar mover este navio depois.
				ships = append(ships, shipRef{ship, r, c})
			}
		}
	}

	if len(ships) == 0 {
		return nil // Nenhum navio para mover.
	}

	// 3. Embaralha a ordem dos navios para que a IA não mova sempre o mesmo navio primeiro.
	rand.Shuffle(len(ships), func(i, j int) { ships[i], ships[j] = ships[j], ships[i] })

	// Tenta mover cada navio até conseguir mover um com sucesso.
	for _, ref := range ships {
		// Tenta todas as 4 direções em ordem aleatória.
		dirs := []entity.Direction{entity.Up, entity.Down, entity.Left, entity.Right}
		rand.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })

		for _, dir := range dirs {
			// Tenta mover o navio na direção escolhida.
			_, err := entityBoard.MoveShipSegment(ref.r, ref.c, dir)
			if err == nil {
				// Sucesso! O movimento foi válido.
				// Sincroniza o tabuleiro visual do inimigo com o novo estado lógico.
				s.syncVisualBoard(s.match.EnemyBoard, entityBoard)
				return nil // A IA move apenas um navio por turno, então retorna após o primeiro sucesso.
			}
		}
	}

	return nil // Se não conseguiu mover nenhum navio, simplesmente termina o turno de movimento.
}

// reconstructEntityBoard cria um EntityBoard temporário analisando o tabuleiro visual.
// Ele usa um algoritmo de busca (BFS) para identificar componentes conectados (navios).
func (s *dynamicBattleService) reconstructEntityBoard(visualBoard *board.Board) *entity.Board {
	eb := &entity.Board{}

	// Estrutura auxiliar para coordenadas.
	type point struct{ r, c int }
	visited := make(map[point]bool)

	// Varre todo o tabuleiro procurando células que fazem parte de navios (Ship ou Hit).
	for r := 0; r < 10; r++ {
		for c := 0; c < 10; c++ {
			state := visualBoard.Cells[r][c].State
			// Se encontrou uma parte de navio não visitada, inicia a descoberta do navio completo.
			if (state == board.Ship || state == board.Hit) && !visited[point{r, c}] {
				// BFS (Busca em Largura) para encontrar todas as células conectadas deste navio.
				var component []point
				queue := []point{{r, c}}
				visited[point{r, c}] = true

				for len(queue) > 0 {
					curr := queue[0]
					queue = queue[1:]
					component = append(component, curr)

					// Verifica vizinhos (cima, baixo, esquerda, direita).
					neighbors := []point{
						{curr.r - 1, curr.c}, {curr.r + 1, curr.c},
						{curr.r, curr.c - 1}, {curr.r, curr.c + 1},
					}
					for _, n := range neighbors {
						if n.r >= 0 && n.r < 10 && n.c >= 0 && n.c < 10 {
							nState := visualBoard.Cells[n.r][n.c].State
							if (nState == board.Ship || nState == board.Hit) && !visited[n] {
								visited[n] = true
								queue = append(queue, n)
							}
						}
					}
				}

				// Determina os limites (bounding box) do navio encontrado para calcular tamanho e orientação.
				minR, maxR, minC, maxC := 10, -1, 10, -1
				for _, p := range component {
					if p.r < minR {
						minR = p.r
					}
					if p.r > maxR {
						maxR = p.r
					}
					if p.c < minC {
						minC = p.c
					}
					if p.c > maxC {
						maxC = p.c
					}
				}

				width := maxC - minC + 1
				height := maxR - minR + 1
				size := len(component)
				horizontal := width >= height // Se largura >= altura, é horizontal.

				// Cria o objeto Ship lógico.
				ship := &entity.Ship{
					Name:       "EnemyShip",
					Size:       size,
					HitCount:   0,
					Horizontal: horizontal,
				}

				// Posiciona o navio no EntityBoard temporário.
				for _, p := range component {
					entity.PlaceShip(&eb.Positions[p.r][p.c], ship)
				}
			}
		}
	}

	// Aplica os estados de ataque (Hit/Miss) no EntityBoard reconstruído.
	for r := 0; r < 10; r++ {
		for c := 0; c < 10; c++ {
			state := visualBoard.Cells[r][c].State
			if state == board.Hit {
				eb.AttackPositionA(r, c)
			} else if state == board.Miss {
				eb.AttackPositionA(r, c)
			}
		}
	}

	return eb
}

// syncVisualBoard atualiza o tabuleiro visual (células) para refletir o estado do EntityBoard.
// Usado após a IA mover seus navios logicamente.
func (s *dynamicBattleService) syncVisualBoard(visualBoard *board.Board, entityBoard *entity.Board) {
	for r := 0; r < 10; r++ {
		for c := 0; c < 10; c++ {
			ship := entity.GetShipReference(entityBoard.Positions[r][c])
			isAttacked := entity.IsAttacked(entityBoard.Positions[r][c])

			// Atualiza o estado visual baseado na presença de navio e se foi atacado.
			if ship != nil {
				if isAttacked {
					visualBoard.Cells[r][c].State = board.Hit
				} else {
					visualBoard.Cells[r][c].State = board.Ship
				}
			} else {
				if isAttacked {
					visualBoard.Cells[r][c].State = board.Miss
				} else {
					visualBoard.Cells[r][c].State = board.Empty
				}
			}
		}
	}
}
