package basic

type Size struct {
	W float32
	H float32
}

func NewSize(w, h float32) Size {
	return Size{W: w, H: h}
}

func (s Size) Scale(factor float32) Size {
	return Size{
		W: s.W * factor,
		H: s.H * factor,
	}
}

func (s Size) Half() Point {
	return Point{
		X: s.W / 2,
		Y: s.H / 2,
	}
}
