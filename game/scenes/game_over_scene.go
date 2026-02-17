package scenes

import (
	"image/color"
	"image/gif"
	"os"
	"time"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameOverScene struct {
	winnerName  string
	danceFrames []*ebiten.Image
	danceDelays []int
	currentImg  *ebiten.Image

	congratsLabel *components.Text
	winnerLabel   *components.Text
	restartLabel  *components.Text
}

func NewGameOverScene(winnerName string) *GameOverScene {
	return &GameOverScene{
		winnerName: winnerName,
	}
}

func (s *GameOverScene) OnEnter(prev Scene, size basic.Size) {
	// Carregar GIF pirate-dance.gif
	f, err := os.Open("assets/images/pirate-dance.gif")
	if err == nil {
		defer f.Close()
		g, err := gif.DecodeAll(f)
		if err == nil {
			s.danceFrames = make([]*ebiten.Image, len(g.Image))
			s.danceDelays = make([]int, len(g.Delay))
			for i, img := range g.Image {
				s.danceFrames[i] = ebiten.NewImageFromImage(img)
				s.danceDelays[i] = g.Delay[i]
			}
			if len(s.danceFrames) > 0 {
				s.currentImg = s.danceFrames[0]
			}
		}
	}

	// Textos
	s.congratsLabel = components.NewText(
		basic.Point{X: 0, Y: 50}, // Mais no topo
		"PARABÉNS!",
		colors.White,
		48,
	)

	s.winnerLabel = components.NewText(
		basic.Point{X: 0, Y: 520}, // Abaixo do GIF
		"Vencedor: "+s.winnerName,
		colors.White,
		32,
	)

	s.restartLabel = components.NewText(
		basic.Point{X: 0, Y: 600},
		"Clique para Recomeçar",
		color.RGBA{200, 200, 200, 255},
		24,
	)

	// Centralizar textos horizontalmente (baseado em tela 1280)
	centerX := float32(1280 / 2)

	// Ajuste manual simples ou pegando tamanho
	cw := s.congratsLabel.GetSize().W
	s.congratsLabel.SetPos(basic.Point{X: centerX - float32(cw)/2, Y: 50})

	ww := s.winnerLabel.GetSize().W
	s.winnerLabel.SetPos(basic.Point{X: centerX - float32(ww)/2, Y: 520})

	rw := s.restartLabel.GetSize().W
	s.restartLabel.SetPos(basic.Point{X: centerX - float32(rw)/2, Y: 600})
}

func (s *GameOverScene) OnExit(next Scene) {}

func (s *GameOverScene) Update() error {
	// Animação do GIF
	if len(s.danceFrames) > 0 {
		totalDuration := 0
		for _, d := range s.danceDelays {
			totalDuration += d * 10
		}
		if totalDuration == 0 {
			totalDuration = 100
		}

		now := int(time.Now().UnixMilli())
		cycleTime := now % totalDuration

		currentDuration := 0
		for k, d := range s.danceDelays {
			frameDuration := d * 10
			if frameDuration == 0 {
				frameDuration = 100
			}
			if cycleTime < currentDuration+frameDuration {
				s.currentImg = s.danceFrames[k]
				break
			}
			currentDuration += frameDuration
		}
	}

	// Clique para reiniciar
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		SwitchTo(NewPlacementScene())
	}

	return nil
}

func (s *GameOverScene) Draw(screen *ebiten.Image) {
	// Fundo escuro
	screen.Fill(color.RGBA{20, 20, 40, 255})

	s.congratsLabel.Draw(screen)
	s.winnerLabel.Draw(screen)

	// Desenhar GIF centralizado
	if s.currentImg != nil {
		w, _ := s.currentImg.Size()
		op := &ebiten.DrawImageOptions{}
		
		// Escala se necessário (opcional)
		scale := 1.2
		op.GeoM.Scale(scale, scale)

		// Centralizar
		// Tela 1280x720 (assumido)
		// Posição Y ~ 130
		x := (1280.0 - float64(w)*scale) / 2
		y := 130.0

		op.GeoM.Translate(x, y)
		screen.DrawImage(s.currentImg, op)
	}

	s.restartLabel.Draw(screen)
}
