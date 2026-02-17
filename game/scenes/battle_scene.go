package scenes

import (
	"fmt"
	"image/color"
	"image/gif"
	"math"
	"os"
	"time"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/game/shared/placement"
	"github.com/allanjose001/go-battleship/game/state"
	"github.com/allanjose001/go-battleship/internal/ai"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type BattleScene struct {
	state       *state.GameState
	hitImage    *ebiten.Image
	fireFrames  []*ebiten.Image
	fireDelays  []int
	missImage   *ebiten.Image
	playerShips []*placement.ShipPlacement

	// Stats
	playerAttempts int
	playerHits     int
	aiAttempts     int
	aiHits         int
	isPlayerTurn   bool

	aiPlayer    *ai.AIPlayer
	entityBoard *entity.Board
	entityFleet *entity.Fleet

	backButtonRow *components.Row
}

func NewBattleScene(gs *state.GameState) *BattleScene {
	ships, _ := gs.PlayerShips.([]*placement.ShipPlacement)
	return &BattleScene{
		state:        gs,
		playerShips:  ships,
		isPlayerTurn: true, // Player starts
	}
}

func (s *BattleScene) OnEnter(prev Scene, size basic.Size) {
	var err error

	// Carregar imagem de fundo
	bg, _, errBg := ebitenutil.NewImageFromFile("assets/images/Mask group.png")
	if errBg == nil {
		s.state.PlayerBoard.BackgroundImage = bg
		s.state.AIBoard.BackgroundImage = bg
	}

	// Carregar Fire.gif como animação
	f, err := os.Open("assets/images/Fire.gif")
	if err == nil {
		defer f.Close()
		g, err := gif.DecodeAll(f)
		if err == nil {
			s.fireFrames = make([]*ebiten.Image, len(g.Image))
			s.fireDelays = make([]int, len(g.Delay))
			for i, img := range g.Image {
				s.fireFrames[i] = ebiten.NewImageFromImage(img)
				s.fireDelays[i] = g.Delay[i]
			}
			if len(s.fireFrames) > 0 {
				s.hitImage = s.fireFrames[0]
			}
		}
	}

	if s.hitImage == nil {
		// Fallback se falhar ao carregar gif
		s.hitImage, _, err = ebitenutil.NewImageFromFile("assets/images/Fire.gif")
		if err != nil {
			// Log erro ou lidar
		}
	}

	s.missImage, _, err = ebitenutil.NewImageFromFile("assets/images/Ponto que já foi atingido 1.png")
	if err != nil {
		// Log erro ou lidar
	}

	// Botão Recomeçar usando Row (mesma lógica de PlacementScene)
	s.backButtonRow = components.NewRow(
		basic.Point{X: 540, Y: 650}, // Posição centralizada manualmente para tela 1280
		0,
		basic.Size{W: 200, H: 50},
		basic.Center,
		basic.Center,
		[]components.Widget{
			components.NewButton(
				basic.Point{},
				basic.Size{W: 200, H: 50},
				"Recomeçar",
				color.RGBA{48, 67, 103, 255},
				colors.White,
				func(b *components.Button) {
					SwitchTo(NewPlacementScene())
				},
			),
		},
	)

	// Inicializar AI
	s.initAI()
}

func (s *BattleScene) OnExit(next Scene) {}

func (s *BattleScene) Update() error {
	s.backButtonRow.Update(basic.Point{})
	if s.isPlayerTurn && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		mouseX, mouseY := float64(mx), float64(my)

		// Ataque no tabuleiro da IA
		if mouseX >= s.state.AIBoard.X && mouseX <= s.state.AIBoard.X+s.state.AIBoard.Size &&
			mouseY >= s.state.AIBoard.Y && mouseY <= s.state.AIBoard.Y+s.state.AIBoard.Size {

			cellSize := s.state.AIBoard.Size / float64(board.Cols)
			col := int((mouseX - s.state.AIBoard.X) / cellSize)
			row := int((mouseY - s.state.AIBoard.Y) / cellSize)

			if col >= 0 && col < board.Cols && row >= 0 && row < board.Rows {
				cell := &s.state.AIBoard.Cells[row][col]
				if cell.State == board.Hit || cell.State == board.Miss {
					return nil // Já foi atacado
				}

				s.playerAttempts++

				if cell.State == board.Ship {
					cell.State = board.Hit
					s.playerHits++

					if s.playerHits >= 17 {
						SwitchTo(NewGameOverScene("Jogador 1"))
						return nil
					}

					// Se acertou, pode continuar jogando? Regra padrão: continua.
					// Se quiser simplificar, troca turno. Vamos trocar turno por enquanto.
				} else if cell.State == board.Empty {
					cell.State = board.Miss
				}

				// Passa a vez para a IA
				s.isPlayerTurn = false

				// Pequeno delay para a IA jogar (opcional, aqui chamamos direto)
				go func() {
					time.Sleep(500 * time.Millisecond)
					s.aiTurn()
					s.isPlayerTurn = true
				}()
			}
		}
	}
	return nil
}

func (s *BattleScene) aiTurn() {
	if s.aiPlayer == nil {
		return
	}

	s.aiAttempts++
	// AI ataca o board "entity"
	// O método Attack da AI modifica o board e o estado interno da AI
	s.aiPlayer.Attack(s.entityBoard)

	// Sincronizar de volta para o board visual
	for r := 0; r < board.Rows; r++ {
		for c := 0; c < board.Cols; c++ {
			entPos := s.entityBoard.Positions[r][c]
			cell := &s.state.PlayerBoard.Cells[r][c]

			// Se a posição foi atacada na AI mas ainda não no visual, atualize
			if entity.IsAttacked(entPos) && cell.State != board.Hit && cell.State != board.Miss {
				if cell.State == board.Ship {
					cell.State = board.Hit
					s.aiHits++
					if s.aiHits >= 17 {
						SwitchTo(NewGameOverScene("Jogador 2"))
						return
					}
				} else {
					cell.State = board.Miss
				}
			}
		}
	}
}

func (s *BattleScene) initAI() {
	// 1. Criar Fleet da entity baseado nos playerShips
	s.entityFleet = entity.NewFleet()
	// Nota: entity.NewFleet cria navios padrão. Precisamos garantir que correspondam.
	// O NewFleet cria: Porta-Aviões(6), Navio de Guerra(4), Encouraçado(3), Encouraçado(3), Submarino(1)
	// Meus shipSizes em setup_utils: 6, 4, 3, 3, 1. Bateu!

	// 2. Criar Board da entity e posicionar navios
	s.entityBoard = &entity.Board{}

	// Mapear ships do jogo para ships da entity
	// Assumindo ordem fixa: 0=Size 6, 1=Size 4, 2=Size 3, 3=Size 3, 4=Size 1
	// Vamos iterar e encontrar correspondentes

	// Para simplificar, vamos limpar o board e reposicionar baseado no s.playerShips
	// Precisamos saber qual navio é qual.
	// s.playerShips tem Size e Orientation e Pos.

	usedShips := make(map[int]bool)

	for _, ps := range s.playerShips {
		if !ps.Placed {
			continue
		}

		// Encontrar um navio disponível na fleet com mesmo tamanho
		var entShip *entity.Ship
		for i, s := range s.entityFleet.Ships {
			if !usedShips[i] && s.Size == ps.Size {
				entShip = s
				usedShips[i] = true
				break
			}
		}

		if entShip != nil {
			// Configurar orientação
			entShip.Horizontal = (ps.Orientation == board.Horizontal)

			// Posicionar no board
			// PlaceShip da entity espera ponteiro de Ship e atualiza as posições
			// ATENÇÃO: PlaceShip espera (ship, row, col). ps.Y é Row, ps.X é Col.
			s.entityBoard.PlaceShip(entShip, ps.Y, ps.X)
		}
	}

	// 3. Inicializar AI Player (Hard para usar estratégias melhores)
	s.aiPlayer = ai.NewHardAIPlayer(s.entityFleet)
}

func (s *BattleScene) Draw(screen *ebiten.Image) {
	// Desenha os tabuleiros
	s.state.PlayerBoard.Draw(screen)
	s.state.AIBoard.Draw(screen)

	// Linha vertical separando os tabuleiros
	lineX := 640.0 // Centro da tela
	lineY1 := s.state.PlayerBoard.Y
	lineY2 := s.state.PlayerBoard.Y + s.state.PlayerBoard.Size
	ebitenutil.DrawLine(screen, lineX, lineY1, lineX, lineY2, colors.White)

	// Desenha infos completas
	s.drawPlayerInfo(screen, s.state.PlayerBoard, "Jogador 1", s.playerAttempts, s.playerHits, s.isPlayerTurn)
	s.drawPlayerInfo(screen, s.state.AIBoard, "Jogador 2", s.aiAttempts, s.aiHits, !s.isPlayerTurn)

	// Desenha o botão voltar
	s.backButtonRow.Draw(screen)

	// Desenha os barcos do jogador
	cellSize := s.state.PlayerBoard.Size / float64(board.Cols)
	for _, ship := range s.playerShips {
		if ship.Placed {
			x := s.state.PlayerBoard.X + float64(ship.X)*cellSize
			y := s.state.PlayerBoard.Y + float64(ship.Y)*cellSize

			op := &ebiten.DrawImageOptions{}
			iw, ih := ship.Image.Size()
			op.GeoM.Scale(
				(cellSize*float64(ship.Size))/float64(iw),
				cellSize/float64(ih),
			)

			if ship.Orientation == board.Vertical {
				op.GeoM.Rotate(math.Pi / 2)
				op.GeoM.Translate(cellSize, 0)
			}

			op.GeoM.Translate(x, y)
			screen.DrawImage(ship.Image, op)
		}
	}

	// Desenha os acertos/erros no tabuleiro do jogador
	s.drawMarkers(screen, s.state.PlayerBoard)
	// Desenha os acertos/erros no tabuleiro da IA
	s.drawMarkers(screen, s.state.AIBoard)
}

func (s *BattleScene) drawPlayerInfo(screen *ebiten.Image, b *board.Board, name string, attempts, hits int, isTurn bool) {
	// Posição base abaixo do tabuleiro
	baseX := b.X
	baseY := b.Y + b.Size + 20

	// Indicador de turno (círculo/quadrado)
	indicatorColor := color.RGBA{255, 0, 0, 255} // Vermelho
	if isTurn {
		indicatorColor = color.RGBA{0, 255, 0, 255} // Verde
	}
	// Desenha quadrado indicador (alinhado à esquerda)
	ebitenutil.DrawRect(screen, baseX, baseY, 20, 20, indicatorColor)

	// Nome do Jogador (ao lado do indicador)
	nameLabel := components.NewText(
		basic.Point{X: float32(baseX + 30), Y: float32(baseY)},
		name,
		colors.White,
		20,
	)
	nameLabel.Draw(screen)

	// Tentativas (alinhado com o indicador/início do tabuleiro)
	attemptsText := fmt.Sprintf("Tentativa: %d", attempts)
	attemptsLabel := components.NewText(
		basic.Point{X: float32(baseX + 30), Y: float32(baseY + 30)},
		attemptsText,
		colors.White,
		16,
	)
	attemptsLabel.Draw(screen)

	// Acertos (alinhado com o indicador/início do tabuleiro)
	hitsText := fmt.Sprintf("Acertos: %d", hits)
	hitsLabel := components.NewText(
		basic.Point{X: float32(baseX + 30), Y: float32(baseY + 55)},
		hitsText,
		colors.White,
		16,
	)
	hitsLabel.Draw(screen)
}

func (s *BattleScene) drawMarkers(screen *ebiten.Image, b *board.Board) {
	cellSize := b.Size / float64(board.Cols)
	for i := 0; i < board.Rows; i++ {
		for j := 0; j < board.Cols; j++ {
			cell := b.Cells[i][j]
			if cell.State == board.Hit || cell.State == board.Miss {
				x := b.X + float64(j)*cellSize
				y := b.Y + float64(i)*cellSize

				if cell.State == board.Hit {
					var img *ebiten.Image
					if len(s.fireFrames) > 0 {
						// Calcular frame baseado no tempo
						// Assumindo 10ms por unidade de delay (padrão GIF)
						totalDuration := 0
						for _, d := range s.fireDelays {
							totalDuration += d * 10 // delay em centésimos de segundo -> ms
						}
						if totalDuration == 0 {
							totalDuration = 100 // fallback
						}

						now := int(time.Now().UnixMilli())
						cycleTime := now % totalDuration

						currentDuration := 0
						for k, d := range s.fireDelays {
							frameDuration := d * 10
							if frameDuration == 0 {
								frameDuration = 100
							} // fallback para delay 0
							if cycleTime < currentDuration+frameDuration {
								img = s.fireFrames[k]
								break
							}
							currentDuration += frameDuration
						}
						if img == nil {
							img = s.fireFrames[0]
						}
					} else {
						img = s.hitImage
					}

					if img != nil {
						op := &ebiten.DrawImageOptions{}
						iw, ih := img.Size()
						op.GeoM.Scale(cellSize/float64(iw), cellSize/float64(ih))
						op.GeoM.Translate(x, y)
						screen.DrawImage(img, op)
					} else {
						// Fallback cor vermelha
						ebitenutil.DrawRect(screen, x+cellSize*0.25, y+cellSize*0.25, cellSize*0.5, cellSize*0.5, color.RGBA{255, 0, 0, 150})
					}
				} else if cell.State == board.Miss && s.missImage != nil {
					op := &ebiten.DrawImageOptions{}
					iw, ih := s.missImage.Size()
					op.GeoM.Scale(cellSize/float64(iw), cellSize/float64(ih))
					op.GeoM.Translate(x, y)
					screen.DrawImage(s.missImage, op)
				} else {
					// Fallback
					c := color.RGBA{200, 200, 200, 150}
					if cell.State == board.Hit {
						c = color.RGBA{255, 0, 0, 150}
					}
					ebitenutil.DrawRect(screen, x+cellSize*0.25, y+cellSize*0.25, cellSize*0.5, cellSize*0.5, c)
				}
			}
		}
	}
}
