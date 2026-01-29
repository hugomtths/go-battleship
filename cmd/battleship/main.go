package main

import (
	//"gobattleship/UI"
	"gobattleship/game"
	"gobattleship/internal/service"
	"fmt"
	"log"
)

func main() {
	//fmt.Println("Backend is running...");
	
	//UI.Run();

	board1 := new(game.Board);

	//game.PrintBoard(board1);

	barco1 := new(game.Ship);
	barco1.Size = 3;
	barco1.Horizontal = true;

	barco2 := new(game.Ship);
	barco2.Size = 3;

	fmt.Println("");

	game.PlaceShip(board1, barco1, 1, 1);

	//game.PrintBoard(board1);

	fmt.Println("");

	//fmt.Println("hit count barco1:", barco1.HitCount);

	game.AttackPosition(board1, 1, 1);

	//fmt.Println("hit count barco1:", barco1.HitCount);

	//========== teste de profile ===========

	profile1 := new(service.Profile);
	profile1.Username = "Player1";
	profile1.TotalScore = 200
	profile1.HighestScore = 50
	profile1.GamesPlayed = 5
	profile1.MedalsEarned = 2

	service.SaveProfile(*profile1);
	err := service.SaveProfile(*profile1)
	if err != nil {
		log.Fatal(err)
	}

	profile2, err := service.FindProfile("Player2");
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("perfil encontrado: %+v\n", profile2);

	//service.RemoveProfile("Player1");

}