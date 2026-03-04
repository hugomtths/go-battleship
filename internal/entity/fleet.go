package entity

type Fleet struct {
	Ships [6]*Ship
}

func NewFleet() *Fleet {
	fleet := &Fleet{}

	fleet.Ships[0] = &Ship{Name: "Porta-Aviões", Size: 6, Horizontal: true}
	fleet.Ships[1] = &Ship{Name: "Porta-Aviões", Size: 6, Horizontal: true}
	fleet.Ships[2] = &Ship{Name: "Navio de Guerra", Size: 4, Horizontal: true}
	fleet.Ships[3] = &Ship{Name: "Navio de Guerra", Size: 4, Horizontal: true}
	fleet.Ships[4] = &Ship{Name: "Encouraçado", Size: 3, Horizontal: true}
	fleet.Ships[5] = &Ship{Name: "Submarino", Size: 1, Horizontal: true}

	return fleet
}

func (fleet *Fleet) IsFleetDestroyed() bool {
	for i := 0; i < len(fleet.Ships); i++ {
		if fleet.Ships[i] != nil && !fleet.Ships[i].IsDestroyed() {
			return false
		}
	}
	return true
}

func (fleet *Fleet) GetFleetShips() (ships []*Ship) {
	return fleet.Ships[:]
}

func (fleet *Fleet) GetShipByIndex(index int) *Ship {
	return fleet.Ships[index]
}
