package game

import (
	"log"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/game/scenes"
	"github.com/hajimehoshi/ebiten/v2"
)

var windowSize = basic.Size{W: 1280, H: 720}

type Game struct {
	// stack que gerencia as rotas das telas do jogo - é como um singleton (única para tod0 o jogo)
	stack *scenes.SceneStack
}

func NewGame() *Game {
	//inicializa fonte ao inicializar game
	components.InitFonts()
	g := &Game{
		stack: scenes.NewSceneStack(windowSize, &scenes.HomeScreen{}), //incializa com primeira scene
	}

	return g
}
func (g *Game) Update() error {

	if g.stack.IsEmpty() {
		return ebiten.Termination
	}
	err := g.stack.Update()
	if err != nil {
		log.Fatal("Erro em stack.Update(): ", err)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//pinta background
	screen.Fill(colors.Background)

	if !g.stack.IsEmpty() {
		g.stack.Draw(screen)
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return int(windowSize.W), int(windowSize.H)
}
