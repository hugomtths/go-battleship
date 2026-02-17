package main

import (
	"log"

	"github.com/allanjose001/go-battleship/game"
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/internal/service"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// 1. Setup de dados do Perfil
	profile := &service.Profile{
		Username:     "malub",
		GamesPlayed:  10,
		MedalsEarned: 4, 
	}
	
	_ = service.SaveProfile(*profile)

	// 2. Recursos Visuais
	components.InitFonts() // Carrega Goldman.ttf para o NewText funcionar

	// 3. Inicializa o Jogo
	g := game.NewGame()

	// 4. Configuração de Janela e Resolução Lógica
	ebiten.SetWindowSize(1280, 720) 
	
	// SetWindowTitle ajuda a identificar a tela de teste
	ebiten.SetWindowTitle("Batalha Naval - Profile Scene Debug")
	
	// Habilitar redimensionamento sem quebrar o layout (usa o Layout() do game.go)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

<<<<<<< develop
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
=======
	// 5. Execução do Loop principal
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
>>>>>>> develop
