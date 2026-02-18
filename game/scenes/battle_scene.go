package scenes

import (
	"fmt"
	"image/color"
	"time"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/game/shared/placement"
	"github.com/allanjose001/go-battleship/game/state"
	"github.com/allanjose001/go-battleship/internal/assets"
	"github.com/allanjose001/go-battleship/internal/service"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// BattleScene representa a tela de batalha em si.
// Ela coordena o serviço de batalha, o renderer e os botões de interface.
type BattleScene struct {
	// svc contém toda a lógica de turnos, ataques e vitória
	svc *service.BattleService
	// backButtonRow é a linha que contém o botão de "Recomeçar"
	backButtonRow *components.Row
	// playerShips guarda os navios posicionados na fase de placement
	playerShips []*placement.ShipPlacement

	hitImage   *ebiten.Image
	missImage  *ebiten.Image
	fireFrames []*ebiten.Image
	fireDelays []int

	playerNameLabel     *components.Text
	playerAttemptsLabel *components.Text
	playerHitsLabel     *components.Text

	aiNameLabel     *components.Text
	aiAttemptsLabel *components.Text
	aiHitsLabel     *components.Text
}

// NewBattleScene recebe o estado de jogo gerado na fase de placement
// e cria um BattleService novo para controlar a batalha.
func NewBattleScene(gs *state.GameState) *BattleScene {
	ships, _ := gs.PlayerShips.([]*placement.ShipPlacement)
	gameSvc := service.NewGameService()
	svc := service.NewBattleService(gs, gameSvc, ships)

	return &BattleScene{
		svc:         svc,
		playerShips: ships,
	}
}

// OnEnter é chamado quando a cena de batalha entra em foco.
// Aqui configuramos o fundo dos tabuleiros e criamos o botão de recomeçar.
func (s *BattleScene) OnEnter(prev Scene, size basic.Size) {
	playerBoard := s.svc.PlayerBoard()
	aiBoard := s.svc.AIBoard()

	// Fundo compartilhado para os dois tabuleiros
	bg, _, errBg := ebitenutil.NewImageFromFile("assets/images/Mask group.png")
	if errBg == nil {
		playerBoard.BackgroundImage = bg
		aiBoard.BackgroundImage = bg
	}

	// Linha com botão "Recomeçar", que volta para a fase de placement
	s.backButtonRow = components.NewRow(
		basic.Point{X: 540, Y: 650},
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

	frames, delays, _ := assets.LoadFireAnimation()
	hit, _ := assets.LoadHitImage()
	miss, _ := assets.LoadMissImage()
	if hit == nil && len(frames) > 0 {
		hit = frames[0]
	}
	s.fireFrames = frames
	s.fireDelays = delays
	s.hitImage = hit
	s.missImage = miss

	playerBaseX := playerBoard.X
	playerBaseY := playerBoard.Y + playerBoard.Size + 20

	aiBaseX := aiBoard.X
	aiBaseY := aiBoard.Y + aiBoard.Size + 20

	s.playerNameLabel = components.NewText(
		basic.Point{X: float32(playerBaseX + 30), Y: float32(playerBaseY)},
		"Jogador 1",
		colors.White,
		20,
	)
	s.playerAttemptsLabel = components.NewText(
		basic.Point{X: float32(playerBaseX + 30), Y: float32(playerBaseY + 30)},
		"",
		colors.White,
		16,
	)
	s.playerHitsLabel = components.NewText(
		basic.Point{X: float32(playerBaseX + 30), Y: float32(playerBaseY + 55)},
		"",
		colors.White,
		16,
	)

	s.aiNameLabel = components.NewText(
		basic.Point{X: float32(aiBaseX + 30), Y: float32(aiBaseY)},
		"Jogador 2",
		colors.White,
		20,
	)
	s.aiAttemptsLabel = components.NewText(
		basic.Point{X: float32(aiBaseX + 30), Y: float32(aiBaseY + 30)},
		"",
		colors.White,
		16,
	)
	s.aiHitsLabel = components.NewText(
		basic.Point{X: float32(aiBaseX + 30), Y: float32(aiBaseY + 55)},
		"",
		colors.White,
		16,
	)
}

// OnExit é chamado quando saímos da cena de batalha.
// Não há limpeza específica necessária aqui.
func (s *BattleScene) OnExit(next Scene) {}

// Update trata entradas do usuário e delega a lógica de turnos ao serviço.
// Se algum jogador vencer, a cena muda para a tela de Game Over.
func (s *BattleScene) Update() error {
	s.backButtonRow.Update(basic.Point{})

	// Clique do mouse no tabuleiro do inimigo (AI)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		mouseX, mouseY := float64(mx), float64(my)

		// O serviço decide se o clique resultou em vitória imediata
		winner := s.svc.HandlePlayerClick(mouseX, mouseY)
		if winner != "" {
			SwitchTo(NewGameOverScene(winner))
			return nil
		}
	}

	// Atualiza o turno da IA e verifica novamente vitória
	winner := s.svc.Update()
	if winner != "" {
		SwitchTo(NewGameOverScene(winner))
		return nil
	}

	return nil
}

// Draw desenha o estado atual da batalha e o botão de recomeçar.
func (s *BattleScene) Draw(screen *ebiten.Image) {
	playerBoard := s.svc.PlayerBoard()
	aiBoard := s.svc.AIBoard()

	playerBoard.Draw(screen)
	aiBoard.Draw(screen)

	lineX := 640.0
	lineY1 := playerBoard.Y
	lineY2 := playerBoard.Y + playerBoard.Size
	ebitenutil.DrawLine(screen, lineX, lineY1, lineX, lineY2, colors.White)

	pAttempts, pHits, aiAttempts, aiHits, isPlayerTurn := s.svc.Stats()

	s.drawPlayerInfo(screen, playerBoard, s.playerNameLabel, s.playerAttemptsLabel, s.playerHitsLabel, pAttempts, pHits, isPlayerTurn)
	s.drawPlayerInfo(screen, aiBoard, s.aiNameLabel, s.aiAttemptsLabel, s.aiHitsLabel, aiAttempts, aiHits, !isPlayerTurn)

	for _, ship := range s.svc.PlayerShips() {
		if ship.Placed {
			components.DrawShip(screen, playerBoard, ship, false, ship.Orientation)
		}
	}

	s.drawMarkers(screen, playerBoard)
	s.drawMarkers(screen, aiBoard)
	s.backButtonRow.Draw(screen)
}

func (s *BattleScene) drawPlayerInfo(
	screen *ebiten.Image,
	b *board.Board,
	nameLabel *components.Text,
	attemptsLabel *components.Text,
	hitsLabel *components.Text,
	attempts int,
	hits int,
	isTurn bool,
) {
	baseX := b.X
	baseY := b.Y + b.Size + 20

	indicatorColor := color.RGBA{255, 0, 0, 255}
	if isTurn {
		indicatorColor = color.RGBA{0, 255, 0, 255}
	}

	ebitenutil.DrawRect(screen, baseX, baseY, 20, 20, indicatorColor)

	nameLabel.Draw(screen)

	attemptsLabel.Text = fmt.Sprintf("Tentativa: %d", attempts)
	attemptsLabel.Draw(screen)

	hitsLabel.Text = fmt.Sprintf("Acertos: %d", hits)
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
					if len(s.fireFrames) > 0 && len(s.fireDelays) == len(s.fireFrames) {
						totalDuration := 0
						for _, d := range s.fireDelays {
							totalDuration += d * 10
						}
						if totalDuration == 0 {
							totalDuration = 100
						}

						now := int(time.Now().UnixMilli())
						cycleTime := now % totalDuration

						currentDuration := 0
						for k, d := range s.fireDelays {
							frameDuration := d * 10
							if frameDuration == 0 {
								frameDuration = 100
							}
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
						ebitenutil.DrawRect(screen, x+cellSize*0.25, y+cellSize*0.25, cellSize*0.5, cellSize*0.5, color.RGBA{255, 0, 0, 150})
					}
				} else if cell.State == board.Miss && s.missImage != nil {
					op := &ebiten.DrawImageOptions{}
					iw, ih := s.missImage.Size()
					op.GeoM.Scale(cellSize/float64(iw), cellSize/float64(ih))
					op.GeoM.Translate(x, y)
					screen.DrawImage(s.missImage, op)
				} else {
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
