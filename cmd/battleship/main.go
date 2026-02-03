package main

import (
	"github.com/allanjose001/go-battleship/game"
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
<<<<<<< Updated upstream
=======

	/*//========= teste inserção no tabuleiro ============
	board1 := new(entity.Board)

	barco1 := new(entity.Ship)
	barco1.Size = 3
	barco1.Horizontal = true

	barco2 := new(entity.Ship)
	barco2.Size = 3
	barco2.Horizontal = false

	entity.PlaceShip(board1, barco1, 1, 1)

	entity.PlaceShip(board1, barco2, 5, 5)

	entity.PrintBoard(board1)

	//fmt.Println("");

	//entity.AttackPosition(board1, 1, 1);

	//========= teste AI ===========

	aiPlayer := ai.NewEasyAIPlayer()

	aiPlayer.Attack(board1)

	entity.PrintBoard(board1)

	//========== teste de profile ===========

	profile1 := new(service.Profile)
	profile1.Username = "Player1"
	profile1.TotalScore = 200
	profile1.HighestScore = 50
	profile1.GamesPlayed = 5
	profile1.MedalsEarned = 2

	service.SaveProfile(*profile1)
	err := service.SaveProfile(*profile1)

	//profile2, err := service.FindProfile("Player2");
	//if err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Printf("perfil encontrado: %+v\n", profile2);

	//service.RemoveProfile("Player1");
	*/
	//============= teste do front ========================

>>>>>>> Stashed changes
	components.InitFonts() //carrega a fonte apenas uma vez
	g := game.NewGame()
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}

}
