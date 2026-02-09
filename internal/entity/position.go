package entity

type Position struct {
	attacked      bool
	blocked       bool
	shipReference *Ship
}

func attack(pos *Position) {
	pos.attacked = true

	if pos.shipReference != nil {
		pos.shipReference.HitCount += 1
	}
}

func Block(pos *Position) {
	pos.blocked = true
}

func Unblock(pos *Position) {
	pos.blocked = false
}

func PlaceShip(pos *Position, ship *Ship) {
	pos.shipReference = ship
}

func RemoveShip(pos *Position) {
	pos.shipReference = nil
}

func IsAttacked(pos Position) bool {
	return pos.attacked
}

func IsBlocked(pos Position) bool {
	return pos.blocked
}

func GetShipReference(pos Position) *Ship {
	return pos.shipReference
}

func IsValidPosition(pos Position) bool {
	return !pos.attacked && !pos.blocked
}
