package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/allanjose001/go-battleship/internal/ai"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/allanjose001/go-battleship/internal/service"
)

func totalShipCells(f *entity.Fleet) int {
	sum := 0
	for _, s := range f.GetFleetShips() {
		if s != nil {
			sum += s.Size
		}
	}
	return sum
}

func findShipTopLeft(b *entity.Board, ship *entity.Ship) (int, int, bool) {
	found := false
	minR, minC := 999, 999
	for r := 0; r < entity.BoardSize; r++ {
		for c := 0; c < entity.BoardSize; c++ {
			if entity.GetShipReference(b.Positions[r][c]) == ship {
				found = true
				if r < minR {
					minR = r
				}
				if c < minC {
					minC = c
				}
			}
		}
	}
	if !found {
		return 0, 0, false
	}
	return minR, minC, true
}

func printHelp() {
	fmt.Println("Comandos:")
	fmt.Println(" show                - mostra tabuleiro do oponente")
	fmt.Println(" showplayer          - mostra tabuleiro do jogador")
	fmt.Println(" attack R C          - ataca posição (R,C) no tabuleiro do oponente")
	fmt.Println(" move IDX DIR        - move navio IDX (up/down/left/right) - consome turno")
	fmt.Println(" turn                - mostra quem tem o turno")
	fmt.Println(" list                - lista navios do jogador")
	fmt.Println(" quit                - sai")
}

// runEnemyTurns chama a IA diretamente no entity.Board do jogador,
// sem passar pelo AttackService/MatchService que exigem o board visual.
func runEnemyTurns(match *entity.Match, aiPlayer *ai.AIPlayer, playerEntityBoard *entity.Board, playerFleet *entity.Fleet) {
	for match.Turn == entity.TurnEnemy && !match.IsFinished() {

		// snapshot das células já atacadas ANTES do ataque
		type cell struct{ r, c int }
		var attackedBefore []cell
		for r := 0; r < entity.BoardSize; r++ {
			for c := 0; c < entity.BoardSize; c++ {
				if entity.IsAttacked(playerEntityBoard.Positions[r][c]) {
					attackedBefore = append(attackedBefore, cell{r, c})
				}
			}
		}

		aiPlayer.Attack(playerEntityBoard)

		// descobre a célula recém-atacada e se foi hit ou miss
		wasAttacked := func(r, c int) bool {
			for _, p := range attackedBefore {
				if p.r == r && p.c == c {
					return true
				}
			}
			return false
		}

		hit := false
		for r := 0; r < entity.BoardSize; r++ {
			for c := 0; c < entity.BoardSize; c++ {
				if entity.IsAttacked(playerEntityBoard.Positions[r][c]) && !wasAttacked(r, c) {
					ship := entity.GetShipReference(playerEntityBoard.Positions[r][c])
					aiPlayer.AdjustStrategy(playerEntityBoard, r, c, ship)
					if ship != nil {
						hit = true
						fmt.Printf("[IA] Acertou em (%d,%d) — %s!\n", r, c, ship.Name)
						if ship.IsDestroyed() {
							fmt.Printf("[IA] %s destruído!\n", ship.Name)
						}
					} else {
						fmt.Printf("[IA] Errou em (%d,%d).\n", r, c)
					}
					goto done
				}
			}
		}
	done:
		fmt.Println("[IA] Estado do seu tabuleiro:")
		entity.PrintBoard(playerEntityBoard)

		if playerFleet.IsFleetDestroyed() {
			fmt.Println("*** IA venceu! Game Over. ***")
			match.Status = entity.MatchStatusFinished
			match.Winner = entity.TurnEnemy
			return
		}

		// hit: IA continua — miss: devolve turno ao jogador
		if !hit {
			match.Turn = entity.TurnPlayer
			match.ClearNextAction()
		}
	}
}

func main() {
	dynSvc := service.NewDynamicMatchService(nil, time.Millisecond)

	// boards lógicos — 1 por jogador, sem board visual
	playerEntityBoard := &entity.Board{}
	enemyEntityBoard := &entity.Board{} // <- tabuleiro DA IA (onde ela posiciona navios)

	playerFleet := entity.NewFleet()
	enemyFleet := entity.NewFleet()

	service.PositionShipsRandomly(enemyEntityBoard, enemyFleet)
	service.PositionShipsRandomly(playerEntityBoard, playerFleet)

	// ownBoard = enemyEntityBoard (tabuleiro onde A IA tem seus navios)
	// enemyFleet no contexto da IA = playerFleet (frota que ELA ataca)
	aiPlayer := ai.NewDynamicAIPlayer(playerFleet, enemyEntityBoard)

	// Create apenas cria o match em memória
	match := dynSvc.Create("test-match", "easy")

	// injeta manualmente as referências runtime que o DynamicMatchService precisa,
	// sem passar boards visuais (nil é seguro pois runEnemyTurns não usa EnemyAttackStep)
	match.PlayerEntityBoard = playerEntityBoard
	match.PlayerFleet = playerFleet
	match.TotalEnemyShipCells = totalShipCells(enemyFleet)
	match.TotalPlayerShipCells = totalShipCells(playerFleet)
	match.Status = entity.MatchStatusInProgress
	match.Turn = entity.TurnPlayer
	match.StartedAt = time.Now()

	fmt.Println("Dynamic match iniciado. Turn inicial:", match.Turn)
	fmt.Printf("Células inimigas: %d | Células jogador: %d\n",
		match.TotalEnemyShipCells, match.TotalPlayerShipCells)
	printHelp()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		if match.IsFinished() {
			fmt.Println("Partida encerrada. Vencedor:", match.Winner)
			return
		}

		fmt.Printf("[Turn: %s] > ", match.Turn)
		if !scanner.Scan() {
			return
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "help":
			printHelp()

		case "show":
			fmt.Println("Tabuleiro do oponente (IA):")
			entity.PrintBoard(enemyEntityBoard)

		case "showplayer":
			fmt.Println("Seu tabuleiro:")
			entity.PrintBoard(playerEntityBoard)

		case "list":
			for i, s := range match.PlayerFleet.GetFleetShips() {
				if s == nil {
					fmt.Printf("%d: <nil>\n", i)
					continue
				}
				fmt.Printf("%d: %s size=%d horiz=%v hit=%d destroyed=%v\n",
					i, s.Name, s.Size, s.Horizontal, s.HitCount, s.IsDestroyed())
			}

		case "attack":
			if match.Turn != entity.TurnPlayer {
				fmt.Println("não é o turno do jogador")
				continue
			}
			if len(parts) < 3 {
				fmt.Println("uso: attack R C")
				continue
			}
			r, e1 := strconv.Atoi(parts[1])
			c, e2 := strconv.Atoi(parts[2])
			if e1 != nil || e2 != nil {
				fmt.Println("coordenadas inválidas")
				continue
			}

			if !enemyEntityBoard.CheckPosition(r, c) {
				fmt.Println("posição inválida ou já atacada")
				continue
			}

			ship := enemyEntityBoard.AttackPositionB(r, c)
			if ship != nil {
				fmt.Printf("[Jogador] Hit em %s em (%d,%d)!\n", ship.Name, r, c)
				if ship.IsDestroyed() {
					fmt.Printf("[Jogador] %s destruído!\n", ship.Name)
				}
				entity.PrintBoard(enemyEntityBoard)

				// Notifica a IA que ela foi atingida nesta posição
				aiPlayer.RegisterIncomingHit(r, c)

				if enemyFleet.IsFleetDestroyed() {
					fmt.Println("*** Jogador venceu! Game Over. ***")
					match.Status = entity.MatchStatusFinished
					match.Winner = entity.TurnPlayer
				}
				// hit: jogador continua, turno não muda
			} else {
				fmt.Printf("[Jogador] Miss em (%d,%d).\n", r, c)
				entity.PrintBoard(enemyEntityBoard)

				match.Turn = entity.TurnEnemy
				match.NextAction = entity.NextActionEnemyAttack

				runEnemyTurns(match, aiPlayer, playerEntityBoard, playerFleet)
			}

		case "move":
			if match.Turn != entity.TurnPlayer {
				fmt.Println("não é o turno do jogador")
				continue
			}
			if len(parts) < 3 {
				fmt.Println("uso: move IDX DIR")
				continue
			}
			idx, err := strconv.Atoi(parts[1])
			if err != nil || idx < 0 || idx >= len(match.PlayerFleet.Ships) {
				fmt.Println("índice inválido")
				continue
			}
			ship := match.PlayerFleet.GetShipByIndex(idx)
			if ship == nil {
				fmt.Println("navio inexistente")
				continue
			}
			r, c, ok := findShipTopLeft(playerEntityBoard, ship)
			if !ok {
				fmt.Println("não foi possível localizar o navio")
				continue
			}
			dr, dc := 0, 0
			switch parts[2] {
			case "up":
				dr = -1
			case "down":
				dr = 1
			case "left":
				dc = -1
			case "right":
				dc = 1
			default:
				fmt.Println("direção inválida (use up/down/left/right)")
				continue
			}

			// MoveShip diretamente no entity.Board (sem passar pelo DynamicMatchService
			// que exige m.PlayerBoard != nil)
			if err := playerEntityBoard.MoveShip(ship, r+dr, c+dc); err != nil {
				fmt.Println("falha ao mover:", err)
				continue
			}
			fmt.Println("movimento aplicado.")
			entity.PrintBoard(playerEntityBoard)

			// move consome turno: passa para IA
			match.Turn = entity.TurnEnemy
			match.NextAction = entity.NextActionEnemyAttack

			runEnemyTurns(match, aiPlayer, playerEntityBoard, playerFleet)

		case "turn":
			fmt.Println("Turn atual:", match.Turn)

		case "quit":
			return

		default:
			fmt.Println("comando desconhecido. use 'help'")
		}
	}
}
