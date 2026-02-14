package entity

type Ship struct {
	Name       string
	Size       int
	HitCount   int
	Horizontal bool
}

func (s *Ship) Rotate() {
	s.Horizontal = !s.Horizontal
}

func (s *Ship) IsDestroyed() bool {
	return s.HitCount >= s.Size
}

func (s *Ship) IsHorizontal() bool {
	return s.Horizontal
}
