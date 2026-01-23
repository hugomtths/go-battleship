package game

type Position struct {
	attacked bool;
	blocked bool;
	shipReference *Ship;
}

func attack(pos *Position) {
	pos.attacked = true;
	
	if (pos.shipReference != nil) {
		pos.shipReference.HitCount += 1;
	}
}

func block(pos *Position) {
	pos.blocked = true;
}

func unblock(pos *Position) {
	pos.blocked = false;
}

func placeShip(pos *Position, ship *Ship) {
	pos.shipReference = ship;
}

func removeShip(pos *Position) {
	pos.shipReference = nil;
}

func isAttacked(pos Position) bool {
	return pos.attacked;
}

func isBlocked(pos Position) bool {
	return pos.blocked;
}

func getShipReference(pos Position) *Ship {
	return pos.shipReference;
}

func isValidPosition(pos Position) bool {
	return !pos.attacked && !pos.blocked;
}