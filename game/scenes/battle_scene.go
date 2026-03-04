package scenes

import (
	"image/color"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/internal/service"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// BattleScene representa a tela de batalha em si.
// Ela coordena o serviço de batalha, o renderer e os botões de interface.
type BattleScene struct {
	// battleSvc orquestra turnos, ataques e persistência de resultado
	battleSvc service.BattleService
	// backButtonContainer é o container que contém a linha de botões "Recomeçar" e "Sair"
	backButtonContainer components.Widget

	assets *BattleAssets

	playerHUD *components.BattleHUD
	aiHUD     *components.BattleHUD
	inputCtrl *components.BattleInput

	boardView *components.BattleBoardView
	divider   *components.VerticalDivider
	StackHandler
}

// NewBattleScene cria a cena de batalha.
// O estado do jogo (Match) deve ser passado via Contexto.
func NewBattleScene() *BattleScene {
	return &BattleScene{}
}

// OnEnter é chamado quando a cena de batalha entra em foco.
// Aqui configuramos o fundo dos tabuleiros, inicializamos o MatchService e criamos o botão de recomeçar.
func (s *BattleScene) OnEnter(prev Scene, size basic.Size) {
	if s.ctx == nil || s.ctx.Match == nil {
		return
	}

	match := s.ctx.Match
	playerBoard := match.PlayerBoard
	aiBoard := match.EnemyBoard

	if s.ctx.BattleService != nil {
		s.battleSvc = s.ctx.BattleService
	} else if svc, err := service.NewBattleServiceFromMatch(match); err == nil {
		s.battleSvc = svc
		s.ctx.SetBattleService(svc)
	}

	// Fundo compartilhado para os dois tabuleiros
	bg, _, errBg := ebitenutil.NewImageFromFile("assets/images/Mask group.png")
	if errBg == nil {
		playerBoard.BackgroundImage = bg
		aiBoard.BackgroundImage = bg
	}

	lineX := 640.0
	lineY1 := playerBoard.Y
	lineY2 := playerBoard.Y + playerBoard.Size
	s.divider = components.NewVerticalDivider(lineX, lineY1, lineY2, colors.White)

	// Linha com botão "Recomeçar" e "Sair"
	row := components.NewRow(
		basic.Point{}, // Posição relativa 0,0 dentro do container
		20,            // Gap
		basic.Size{W: 400, H: 50},
		basic.Center,
		basic.Center,
		[]components.Widget{
			components.NewButton(
				basic.Point{},
				basic.Size{W: 150, H: 50},
				"Recomeçar",
				color.RGBA{48, 67, 103, 255},
				colors.White,
				func(b *components.Button) {
					if match.Profile != nil {
						SwitchTo(NewPlacementSceneWithProfile(match.Profile))
					} else {
						SwitchTo(NewPlacementScene())
					}
				},
			),
			components.NewButton(
				basic.Point{},
				basic.Size{W: 150, H: 50},
				"Sair",
				colors.Red,
				colors.White,
				func(b *components.Button) {
					SwitchTo(&SelectProfileScene{})
				},
			),
		},
	)

	s.backButtonContainer = components.NewContainer(
		basic.Point{X: 440, Y: 650}, // Centralizado (640 - 200) e abaixo dos stats
		basic.Size{W: 400, H: 50},
		0,                  // Radius 0 (não visível)
		colors.Transparent, // Cor transparente
		basic.Center,       // MainAlign
		basic.Center,       // CrossAlign
		row,
	)

	if s.assets == nil {
		s.assets = LoadBattleAssets()
	}

	s.boardView = components.NewBattleBoardView(s.assets.FireAnimation, s.assets.HitImage, s.assets.MissImage)

	playerBaseX := playerBoard.X
	playerBaseY := playerBoard.Y + playerBoard.Size + 40

	aiBaseX := aiBoard.X
	aiBaseY := aiBoard.Y + aiBoard.Size + 40

	playerNameLabel := components.NewText(
		basic.Point{X: float32(playerBaseX + 30), Y: float32(playerBaseY)},
		func() string {
			if match.Profile != nil && match.Profile.Username != "" {
				return match.Profile.Username
			}
			return "Jogador 1"
		}(),
		colors.White,
		20,
	)
	playerAttemptsLabel := components.NewText(
		basic.Point{X: float32(playerBaseX + 30), Y: float32(playerBaseY + 30)},
		"",
		colors.White,
		16,
	)
	playerHitsLabel := components.NewText(
		basic.Point{X: float32(playerBaseX + 30), Y: float32(playerBaseY + 55)},
		"",
		colors.White,
		16,
	)

	aiNameLabel := components.NewText(
		basic.Point{X: float32(aiBaseX + 30), Y: float32(aiBaseY)},
		"IA_MAR",
		colors.White,
		20,
	)
	aiAttemptsLabel := components.NewText(
		basic.Point{X: float32(aiBaseX + 30), Y: float32(aiBaseY + 30)},
		"",
		colors.White,
		16,
	)
	aiHitsLabel := components.NewText(
		basic.Point{X: float32(aiBaseX + 30), Y: float32(aiBaseY + 55)},
		"",
		colors.White,
		16,
	)

	playerHUD := components.NewBattleHUD(playerNameLabel, playerAttemptsLabel, playerHitsLabel, s.battleSvc, components.SidePlayer)
	aiHUD := components.NewBattleHUD(aiNameLabel, aiAttemptsLabel, aiHitsLabel, s.battleSvc, components.SideAI)

	s.playerHUD = playerHUD
	s.aiHUD = aiHUD
	s.inputCtrl = components.NewBattleInput(aiBoard)
}

// OnExit é chamado quando saímos da cena de batalha.
// Não há limpeza específica necessária aqui.
func (s *BattleScene) OnExit(next Scene) {}

// Update trata entradas do usuário e delega a lógica de turnos ao serviço.
// Se algum jogador vencer, a cena muda para a tela de Game Over.
func (s *BattleScene) Update() error {
	s.backButtonContainer.Update(basic.Point{})

	// atualiza HUDs
	if s.playerHUD != nil {
		s.playerHUD.Update(basic.Point{})
	}
	if s.aiHUD != nil {
		s.aiHUD.Update(basic.Point{})
	}

	if s.battleSvc == nil {
		return nil
	}

	if s.inputCtrl != nil {
		row, col, ok := s.inputCtrl.ClickedCell()
		if ok {
			if res, err := s.battleSvc.HandlePlayerClick(row, col); err == nil && res != nil {
				SwitchTo(NewGameOverScene(s.battleSvc.WinnerName()))
				return nil
			}
		}
	}

	if res, err := s.battleSvc.HandleEnemyTurn(); err == nil && res != nil {
		SwitchTo(NewGameOverScene(s.battleSvc.WinnerName()))
		return nil
	}

	return nil
}

// Draw desenha o estado atual da batalha e o botão de recomeçar.
func (s *BattleScene) Draw(screen *ebiten.Image) {
	if s.ctx == nil || s.ctx.Match == nil {
		return
	}
	match := s.ctx.Match
	playerBoard := match.PlayerBoard
	aiBoard := match.EnemyBoard

	playerBoard.Draw(screen)
	aiBoard.Draw(screen)

	if s.divider != nil {
		s.divider.Draw(screen)
	}

	if s.playerHUD != nil {
		s.playerHUD.Draw(screen, playerBoard)
	}
	if s.aiHUD != nil {
		s.aiHUD.Draw(screen, aiBoard)
	}

	if s.boardView != nil {
		s.boardView.DrawBoard(screen, playerBoard, match.PlayerShips)
		s.boardView.DrawBoard(screen, aiBoard, nil)
	}
	s.backButtonContainer.Draw(screen)
}
