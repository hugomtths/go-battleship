package components

// LayoutWidget interface para tratar col e row em certas ocasiões e evitar gambiarra
type LayoutWidget interface {
	IsLayout() bool
}
