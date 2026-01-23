package main
import (
	"github.com/allanjose001/go-battleship/UI"
	"github.com/allanjose001/go-battleship/game"
	"fmt"
)

func main() {
	fmt.Println("Backend is running...");
	
	ui.Run();

	board1 := new(game.Board);

	game.PrintBoard(board1);

	barco1 := new(game.Ship);
	barco1.Size = 3;
	barco1.Horizontal = true;

	barco2 := new(game.Ship);
	barco2.Size = 3;

	fmt.Println("");

	game.PlaceShip(board1, barco1, 1, 1);

	game.PrintBoard(board1);

	fmt.Println("");

	fmt.Println("hit count barco1:", barco1.HitCount);

	game.AttackPosition(board1, 1, 1);

	fmt.Println("hit count barco1:", barco1.HitCount);

}