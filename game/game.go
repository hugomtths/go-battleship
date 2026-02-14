package game

import (
	"image/color"

	"github.com/allanjose001/go-battleship/game/state"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/scenes"
	"github.com/hajimehoshi/ebiten/v2"
)

var windowSize = basic.Size{W: 1280, H: 720}

type Game struct {
	scene scenes.Scene
}

func NewGame() *Game {
	// 1. Inicializa o estado global do jogo (onde ficam os dados de tabuleiros, etc)
    state := &state.GameState{} 

    // 2. Cria a cena de perfil passando o estado
    g := &Game{
        scene: scenes.NewProfileScene(state),
    }

    // 3. Notifica a cena que ela entrou em foco
    g.scene.OnEnter(nil, windowSize) 
    
    return g
}

func (g *Game) Update() error {
	err := g.scene.Update()
	return err
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 30, G: 30, B: 30, A: 255})
	g.scene.Draw(screen)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return int(windowSize.W), int(windowSize.H)
}
