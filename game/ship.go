package game

type Ship struct {
	Name string;
	Size int;
	HitCount int;
	Horizontal bool;
}

func isDestroyed(s *Ship) bool {
	return s.HitCount >= s.Size;
}

func isHorizontal(s *Ship) bool {
	return s.Horizontal;
}

