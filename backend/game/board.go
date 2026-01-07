package game

type Board struct {
	positions [10][10]Position;
}

func PrintBoard(b *Board) {
	for i :=0; i<10; i++ { // itera pelas linhas
		for j:=0; j<10; j++ { // itera pelas colunas
			if (isAttacked(b.positions[i][j])) { // se a posição foi atacada
				if (getShipReference(b.positions[i][j]) != nil) {
					print("x "); // posição atacada com navio
					continue;
				}

				print("o "); // posição atacada sem navio
				continue;
			} else if (getShipReference(b.positions[i][j]) != nil) {
				print("B "); // marca como bloqueada.
				continue;
			}

			//posição valida e não atacada.
			print("- ");
		}
		print("\n"); // nova linha apos cada linha do tabuleiro

	}
}