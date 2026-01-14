package game

type Fleet struct {
	Ships [5]*Ship;
}

func NewFleet() *Fleet {
	fleet := &Fleet{}

	fleet.Ships[0] = &Ship{Name: "Porta-Aviões", Size: 6, Horizontal: true}
	fleet.Ships[1] = &Ship{Name: "Navio de Guerra", Size: 4, Horizontal: true}
	fleet.Ships[2] = &Ship{Name: "Encouraçado", Size: 3, Horizontal: true}
	fleet.Ships[3] = &Ship{Name: "Encouraçado", Size: 3, Horizontal: true}
	fleet.Ships[4] = &Ship{Name: "Submarino", Size: 1, Horizontal: true}

	return fleet
}

func isFleetDestroyed(fleet *Fleet) bool {
	for i := 0; i < 5; i++ {
		if !isDestroyed(fleet.Ships[i]) { //se algum navio ainda não foi destruido
			return false //retorna falso
		}
	}

	return true
}

func getFleetShips(Ships [5]*Ship) (ships []*Ship) {
	return Ships[:]
}

func getShipByIndex(fleet *Fleet, index int) *Ship {
	return fleet.Ships[index]
}
