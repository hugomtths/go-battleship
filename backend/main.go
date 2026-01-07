package main
import (
	"gobattleship/game"
	"fmt"
)

func main() {
	fmt.Println("Backend is running...");

	board1 := new(game.Board);

	game.PrintBoard(board1);

	

}