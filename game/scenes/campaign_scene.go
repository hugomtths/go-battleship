package scenes

import (
	"fmt"
	"image/color"
	"time"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/allanjose001/go-battleship/internal/service"
	"github.com/hajimehoshi/ebiten/v2"
)

// CampaignScene gerencia a seleção de fases do modo campanha
type CampaignScene struct {
	root components.Widget
	StackHandler
}

func (c *CampaignScene) GetMusic() string {
	return "menus"
}

func (c *CampaignScene) OnEnter(prev Scene, size basic.Size) {
	c.refreshUI(size)
	_ = c.Update()
	c.stack.ctx.CanPopOrPush = true
}

func (c *CampaignScene) refreshUI(size basic.Size) {
	profile := c.ctx.Profile
	if profile == nil {
		// Se não houver perfil, volta para evitar erro
		c.stack.Pop()
		return
	}

	// 1. Calcular Pontuação e Estado das Fases
	totalScore := 0

	for _, oldCamp := range profile.Campaigns {
		for _, res := range oldCamp.DifficultyStep {
			if res.Win {
				totalScore += res.Score
			}
		}
	}

	// Estados possíveis: "locked", "current", "done"
	states := map[string]string{
		"easy":   "current", // Padrão: começa no easy
		"medium": "locked",
		"hard":   "locked",
	}
	results := map[string]*entity.MatchResult{
		"easy":   nil,
		"medium": nil,
		"hard":   nil,
	}

	if profile.CurrentCampaign != nil {
		// Verifica Easy
		if res, ok := profile.CurrentCampaign.DifficultyStep["easy"]; ok && res.Win {
			states["easy"] = "done"
			states["medium"] = "current"
			results["easy"] = &res
			totalScore += res.Score
		}
		// Verifica Medium (só se easy estiver done, verificado acima implicitamente pela lógica de progressão)
		if res, ok := profile.CurrentCampaign.DifficultyStep["medium"]; ok && res.Win {
			states["medium"] = "done"
			states["hard"] = "current"
			results["medium"] = &res
			totalScore += res.Score
		}
		// Verifica Hard
		if res, ok := profile.CurrentCampaign.DifficultyStep["hard"]; ok && res.Win {
			states["hard"] = "done"
			results["hard"] = &res
			totalScore += res.Score
		}
	}

	// 2. Construir UI
	title := components.NewText(basic.Point{}, "Modo Campanha", colors.White, 35)
	scoreText := components.NewText(basic.Point{}, fmt.Sprintf("Pontuação Acumulada: %d", totalScore), colors.White, 24)

	stages := []struct {
		id, title string
	}{
		{"easy", "Recruta"},
		{"medium", "Imediato"},
		{"hard", "Almirante"},
	}

	var cards []components.Widget
	for _, s := range stages {
		cards = append(cards, c.createStageCard(s.title, s.id, states[s.id], results[s.id]))
	}

	stagesRow := components.NewContainer(
		basic.Point{}, basic.Size{W: 220, H: 280}, 0,
		colors.Transparent, basic.Center, basic.Center,
		components.NewRow(
			basic.Point{},
			20,
			basic.Size{W: 220, H: 280},
			basic.Center,
			basic.Center,
			cards,
		),
	)

	backBtn := components.NewButton(
		basic.Point{},
		basic.Size{W: 220, H: 50},
		"Voltar",
		colors.Dark,
		colors.White,
		func(b *components.Button) {
			c.ctx.SoundService.PlaySFX("backclick", 0.8)
			c.stack.Pop()
		},
	)

	spacer := components.NewContainer(
		basic.Point{}, basic.Size{W: 1, H: 1}, 0,
		colors.Transparent, basic.Center, basic.Center,
		nil,
	)
	spacer2 := components.NewContainer(
		basic.Point{}, basic.Size{W: 1, H: 20}, 0,
		colors.Transparent, basic.Center, basic.Center,
		nil,
	)

	spacer3 := components.NewContainer(
		basic.Point{}, basic.Size{W: 1, H: 20}, 0,
		colors.Transparent, basic.Center, basic.Center,
		nil,
	)

	c.root = components.NewColumn(
		basic.Point{},
		40,
		size,
		basic.Start,
		basic.Center,
		[]components.Widget{
			spacer,
			title,
			spacer2,
			scoreText,
			stagesRow,
			spacer3,
			backBtn,
		},
	)
}

func (c *CampaignScene) createStageCard(title, diff, state string, res *entity.MatchResult) components.Widget {
	cardSize := basic.Size{W: 220, H: 280}
	bgColor := colors.NightBlue
	textColor := colors.White

	var content []components.Widget

	// Título da fase
	content = append(content, components.NewText(basic.Point{}, title, textColor, 28))

	// Conteúdo variável baseada no estado
	switch state {
	case "done":
		bgColor = color.RGBA{40, 100, 40, 255} // Verde escuro para concluído
		content = append(content, components.NewText(basic.Point{}, "CONCLUÍDO", colors.White, 18))

		// Botão Batalhar (Rejogar)
		btn := components.NewButton(
			basic.Point{},
			basic.Size{W: 160, H: 40},
			"Batalhar",
			colors.Dark,
			colors.White,
			func(b *components.Button) {
				// Inicializa campanha se for a primeira vez
				if c.ctx.Profile.CurrentCampaign == nil {
					c.ctx.Profile.CurrentCampaign = &entity.Campaign{
						ID:             fmt.Sprintf("camp_%d", time.Now().Unix()),
						DifficultyStep: make(map[string]entity.MatchResult),
						IsActive:       true,
					}
					_ = service.UpdateProfile(*c.ctx.Profile)
				}
				// Configura dificuldade no contexto e vai para posicionamento
				c.ctx.SetDifficulty(diff)
				c.ctx.IsCampaign = true

				// Inicia a série de 3 partidas (Partida 1, Placar 0-0)
				ps := NewPlacementSceneWithProfile(c.ctx.Profile)
				ps.SetSeriesState(1, 0, 0)
				c.ctx.SoundService.PlaySFX("click", 0.8)
				c.stack.Push(ps)
			},
		)
		content = append(content, btn)

		// Botão Histórico
		histBtn := components.NewButton(
			basic.Point{},
			basic.Size{W: 160, H: 40},
			"Histórico",
			colors.Dark,
			colors.White,
			func(b *components.Button) {
				c.ctx.SoundService.PlaySFX("click", 0.8)
				c.stack.Push(NewCampaignHistoryScene(diff, title))
			},
		)
		content = append(content, histBtn)
	case "current":
		bgColor = colors.Blue // Azul destaque para atual
		content = append(content, components.NewText(basic.Point{}, "ATUAL", colors.White, 20))
		btn := components.NewButton(
			basic.Point{},
			basic.Size{W: 160, H: 40},
			"Batalhar",
			colors.Dark,
			colors.White,
			func(b *components.Button) {
				// Inicializa campanha se for a primeira vez
				if c.ctx.Profile.CurrentCampaign == nil {
					c.ctx.Profile.CurrentCampaign = &entity.Campaign{
						ID:             fmt.Sprintf("camp_%d", time.Now().Unix()),
						DifficultyStep: make(map[string]entity.MatchResult),
						IsActive:       true,
					}
					_ = service.UpdateProfile(*c.ctx.Profile)
				}
				// Configura dificuldade no contexto e vai para posicionamento
				c.ctx.SetDifficulty(diff)
				c.ctx.IsCampaign = true

				// Inicia a série de 3 partidas (Partida 1, Placar 0-0)
				ps := NewPlacementSceneWithProfile(c.ctx.Profile)
				ps.SetSeriesState(1, 0, 0)
				c.ctx.SoundService.PlaySFX("click", 0.8)
				c.stack.Push(ps)
			},
		)
		content = append(content, btn)
	case "locked":
		bgColor = color.RGBA{50, 50, 50, 255} // Cinza para bloqueado
		textColor = color.RGBA{150, 150, 150, 255}
		content = append(content, components.NewText(basic.Point{}, "BLOQUEADO", textColor, 20))
	}

	return components.NewContainer(
		basic.Point{},
		cardSize,
		10,
		bgColor,
		basic.Center,
		basic.Center,
		components.NewColumn(
			basic.Point{},
			20,
			cardSize,
			basic.Center,
			basic.Center,
			content,
		),
	)
}

func (c *CampaignScene) OnExit(next Scene) {
	c.stack.ctx.CanPopOrPush = false
}

func (c *CampaignScene) Update() error {
	if c.root != nil {
		c.root.Update(basic.Point{})
	}
	return nil
}

func (c *CampaignScene) Draw(screen *ebiten.Image) {
	if c.root != nil {
		c.root.Draw(screen)
	}
}
