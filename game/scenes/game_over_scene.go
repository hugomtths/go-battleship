package scenes

import (
	"fmt"
	"image/color"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameOverScene struct {
	winnerName  string
	result      *entity.MatchResult
	layout      components.LayoutWidget
	actionLabel string
	onAction    func()
	StackHandler
}

func NewGameOverScene(winnerName string, result *entity.MatchResult, actionLabel string, onAction func()) *GameOverScene {
	return &GameOverScene{
		winnerName:  winnerName,
		result:      result,
		actionLabel: actionLabel,
		onAction:    onAction,
	}
}

func (s *GameOverScene) GetMusic() string {
	return "loss"
}

func (s *GameOverScene) OnEnter(prev Scene, size basic.Size) {
	// Layout unificado para vitória e derrota, mudando apenas Título e GIF
	s.setupGameOverLayout(size)
	_ = s.Update()
	s.stack.ctx.CanPopOrPush = true
}

func (s *GameOverScene) setupGameOverLayout(size basic.Size) {
	titleText := "VOCÊ GANHOU"
	titleColor := colors.Green
	gifPath := "assets/images/pirate-dance.gif"
	gifScale := 0.7

	// Se perdeu, muda título e gif
	if s.result != nil && !s.result.Win {
		titleText = "VOCÊ PERDEU"
		titleColor = colors.Red
		gifPath = "assets/images/pirate-funny-treasure-box.gif"
		gifScale = 2 // Aumenta o tamanho do pirata na derrota
	}

	// -------------------------
	// Título
	// -------------------------

	titleLabel := components.NewText(
		basic.Point{},
		titleText,
		titleColor,
		48,
	)

	// -------------------------
	// Nome do vencedor
	// -------------------------

	winnerLabel := components.NewText(
		basic.Point{},
		fmt.Sprintf("Vencedor: %s", s.winnerName),
		colors.White,
		32,
	)

	// -------------------------
	// Estatísticas
	// -------------------------

	var statsChildren []components.Widget

	if s.result != nil {

		statsData := []string{
			fmt.Sprintf("Disparos: %d", s.result.PlayerShots),
			fmt.Sprintf("Maior Sequência: %d", s.result.HigherHitSequence),
			fmt.Sprintf("Navios Perdidos: %d", s.result.LostShips),
			fmt.Sprintf("Acertos: %d", s.result.Hits),
			fmt.Sprintf("Duração: %s", s.result.FormattedDuration()),
		}

		for _, stat := range statsData {

			txt := components.NewText(
				basic.Point{},
				stat,
				colors.White,
				24,
			)

			statsChildren = append(statsChildren, txt)
		}

	}

	// Coluna das estatísticas

	statsColumn := components.NewColumn(
		basic.Point{},
		12,
		basic.Size{W: 600, H: 100},
		basic.End,
		basic.End,
		statsChildren,
	)

	// -------------------------
	// Linha divisória
	// -------------------------

	divider := components.NewContainer(
		basic.Point{},
		basic.Size{W: 4, H: 300}, // ALTERADO: altura ajustada
		0,
		colors.White,
		basic.Start,
		basic.Start,
		nil,
	)

	// -------------------------
	// GIF animado
	// -------------------------

	gifWidget, _ := components.NewGIFWidget(
		gifPath,
		basic.Point{},
		gifScale,
	)

	var rightWidget components.Widget

	if gifWidget != nil {
		rightWidget = gifWidget
	} else {
		rightWidget = components.NewText(
			basic.Point{},
			"",
			colors.White,
			20,
		)
	}

	// -------------------------
	// Linha central
	// -------------------------

	centerRow := components.NewRow(
		basic.Point{},
		40,
		basic.Size{W: size.W - 100, H: 320},
		basic.Center,
		basic.Center,
		[]components.Widget{
			statsColumn,
			divider,
			rightWidget,
		},
	)

	// -------------------------
	// Botão
	// -------------------------

	label := "Voltar"
	if s.actionLabel != "" {
		label = s.actionLabel
	}

	restartBtn := components.NewButton(
		basic.Point{},
		basic.Size{W: 350, H: 60},
		label,
		color.RGBA{65, 81, 100, 255},
		colors.White,
		func(b *components.Button) {
			if s.onAction != nil {
				s.onAction()
			} else {
				s.ctx.SoundService.PlaySFX("backclick", 0.8)
				SwitchTo(&ModeSelectionScene{})
			}
		},
	)

	// Espaço antes do botão

	spacer := components.NewContainer(
		basic.Point{},
		basic.Size{W: 1, H: 140},
		0,
		color.RGBA{},
		basic.Start,
		basic.Start,
		nil,
	)

	// -------------------------
	// Layout principal
	// -------------------------

	mainColumn := components.NewColumn(
		basic.Point{},
		20,
		size,
		basic.Center,
		basic.Center,
		[]components.Widget{
			titleLabel,
			winnerLabel,
			centerRow,
			spacer,
			restartBtn,
		},
	)

	s.layout = mainColumn
}

func (s *GameOverScene) OnExit(next Scene) {
	s.stack.ctx.CanPopOrPush = false
}

func (s *GameOverScene) Update() error {

	if s.layout != nil {
		s.layout.Update(basic.Point{})
	}

	return nil
}

func (s *GameOverScene) Draw(screen *ebiten.Image) {

	screen.Fill(colors.Background)

	if s.layout != nil {
		s.layout.Draw(screen)
	}
}
