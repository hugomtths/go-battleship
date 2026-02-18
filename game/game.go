package game

import (
	"log"

	"github.com/allanjose001/go-battleship/game/components"
	//"github.com/allanjose001/go-battleship/game/state"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/game/scenes"
	"github.com/hajimehoshi/ebiten/v2"
)

var windowSize = basic.Size{W: 1280, H: 800}

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

	scenes.SwitchTo = func(next scenes.Scene) {
		g.stack.Replace(next)
	}

	return g

	// 1. Inicializa o estado global do jogo (onde ficam os dados de tabuleiros, etc)
    //state := &state.GameState{} 

    // 2. Cria a cena de perfil passando o estado
    //g := &Game{
        //scene: scenes.NewProfileScene(state),
    //}

    // 3. Notifica a cena que ela entrou em foco
    //g.scene.OnEnter(nil, windowSize) 
    
    //return g

}

func (g *Game) Update() error {

	if g.stack.IsEmpty() {
		return ebiten.Termination
	}
	err := g.stack.Update()
	if err != nil {
		log.Fatal("Erro em stack.Update() em game.go: ", err)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colors.Background)

	if !g.stack.IsEmpty() {
		g.stack.Draw(screen)
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return int(windowSize.W), int(windowSize.H)
}
