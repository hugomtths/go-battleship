package game

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/state"

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
		stack: scenes.NewSceneStack(windowSize, &scenes.HomeScreen{}, state.NewGameContext()), //incializa com primeira scene
	}

	scenes.SwitchTo = func(next scenes.Scene) {
		g.stack.Replace(next)
	}

	return g

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

// SetGameWindowIcon carrega um PNG e define como ícone da janela do jogo
func SetGameWindowIcon(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			fmt.Println("erro ao fechar arquivo:", cerr)
		}
	}()

	img, err := png.Decode(f)
	if err != nil {
		return err
	}

	ebiten.SetWindowIcon([]image.Image{img})
	return nil
}
