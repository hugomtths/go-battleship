package game

type Ship struct {
	name string;
	size int;
	hitCount int;
	horizontal bool;
}

func isDestroyed(s *Ship) bool {
	return s.hitCount >= s.size;
}

func isHorizontal(s *Ship) bool {
	return s.horizontal;
}

