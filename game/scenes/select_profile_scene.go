package scenes

import (
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	inputhelper "github.com/allanjose001/go-battleship/game/util"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/allanjose001/go-battleship/internal/service"
	"github.com/hajimehoshi/ebiten/v2"
)

type SelectProfileScene struct {
	root       components.Widget
	offset     int
	profiles   []entity.Profile
	screenSize basic.Size
	dragging   bool
	lastMouseY int
	dragAccum  float64
	StackHandler
}

func (s *SelectProfileScene) OnEnter(prev Scene, size basic.Size) {
	s.profiles = service.GetProfiles()
	s.offset = 0
	s.screenSize = size
	s.root = s.buildUI(size)
}

func (s *SelectProfileScene) OnExit(next Scene) {}

func (s *SelectProfileScene) Update() error {
	if s.root != nil {
		s.root.Update(basic.Point{})
	}

	if len(s.profiles) > 4 {
		_, my := ebiten.CursorPosition()
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			if !s.dragging {
				s.dragging = true
				s.lastMouseY = my
			} else {
				dy := my - s.lastMouseY
				if dy != 0 {
					s.dragAccum += float64(dy)
					s.lastMouseY = my

					step := 30.0

					for s.dragAccum <= -step && s.offset+4 < len(s.profiles) {
						s.offset++
						s.dragAccum += step
						s.root = s.buildUI(s.screenSize)
					}

					for s.dragAccum >= step && s.offset > 0 {
						s.offset--
						s.dragAccum -= step
						s.root = s.buildUI(s.screenSize)
					}
				}
			}
		} else {
			s.dragging = false
			s.dragAccum = 0
		}
	}

	return nil
}

func (s *SelectProfileScene) Draw(screen *ebiten.Image) {
	if s.root != nil {
		s.root.Draw(screen)
	}
}

type deleteProfileIcon struct {
	img      *components.Image
	pos      basic.Point
	size     basic.Size
	callback func()
}

func newDeleteProfileIcon(size basic.Size, cb func()) *deleteProfileIcon {
	img, err := components.NewImage("assets/images/apagar-simbolo.png", basic.Point{}, size)
	if err != nil {
		return nil
	}
	return &deleteProfileIcon{
		img:      img,
		pos:      basic.Point{},
		size:     size,
		callback: cb,
	}
}

func (d *deleteProfileIcon) GetPos() basic.Point {
	return d.pos
}

func (d *deleteProfileIcon) SetPos(p basic.Point) {
	d.pos = p
}

func (d *deleteProfileIcon) GetSize() basic.Size {
	return d.size
}

func (d *deleteProfileIcon) Update(offset basic.Point) {
	if d.img == nil {
		return
	}

	currentPos := d.pos.Add(offset)
	d.img.SetPos(currentPos)
	d.img.Update(basic.Point{})

	mx, my := ebiten.CursorPosition()
	if inputhelper.IsClicked(mx, my, currentPos, d.size) && d.callback != nil {
		d.callback()
	}
}

func (d *deleteProfileIcon) Draw(screen *ebiten.Image) {
	if d.img != nil {
		d.img.Draw(screen)
	}
}

type playProfileIcon struct {
	img      *components.Image
	pos      basic.Point
	size     basic.Size
	callback func()
}

func newPlayProfileIcon(size basic.Size, cb func()) *playProfileIcon {
	img, err := components.NewImage("assets/images/botao-play.png", basic.Point{}, size)
	if err != nil {
		return nil
	}
	return &playProfileIcon{
		img:      img,
		pos:      basic.Point{},
		size:     size,
		callback: cb,
	}
}

func (p *playProfileIcon) GetPos() basic.Point {
	return p.pos
}

func (p *playProfileIcon) SetPos(pos basic.Point) {
	p.pos = pos
}

func (p *playProfileIcon) GetSize() basic.Size {
	return p.size
}

func (p *playProfileIcon) Update(offset basic.Point) {
	if p.img == nil {
		return
	}

	currentPos := p.pos.Add(offset)
	p.img.SetPos(currentPos)
	p.img.Update(basic.Point{})

	mx, my := ebiten.CursorPosition()
	if inputhelper.IsClicked(mx, my, currentPos, p.size) && p.callback != nil {
		p.callback()
	}
}

func (p *playProfileIcon) Draw(screen *ebiten.Image) {
	if p.img != nil {
		p.img.Draw(screen)
	}
}

func (s *SelectProfileScene) buildUI(size basic.Size) components.Widget {
	title := components.NewText(
		basic.Point{},
		"Jogadores",
		colors.White,
		22,
	)

	spacer := components.NewContainer(
		basic.Point{},
		basic.Size{W: size.W, H: 50},
		0,
		colors.Transparent,
		basic.Center,
		basic.Center,
		nil,
	)

	listSize := basic.Size{W: size.W * 0.7, H: 200}
	list := s.buildList(listSize)
	listArea := components.NewContainer(
		basic.Point{},
		listSize,
		0,
		colors.Transparent,
		basic.Center,
		basic.Center,
		list,
	)

	backButton := components.NewButton(
		basic.Point{},
		basic.Size{W: 220, H: 55},
		"Voltar",
		colors.Dark,
		nil,
		func(b *components.Button) {
			if SwitchTo != nil {
				SwitchTo(&HomeScreen{})
			}
		},
	)

	newPlayer := components.NewButton(
		basic.Point{},
		basic.Size{W: 220, H: 55},
		"Novo Jogador",
		colors.Dark,
		nil,
		func(b *components.Button) {
			if SwitchTo != nil {
				SwitchTo(&CreateProfileScene{})
			}
		},
	)

	buttonRow := components.NewRow(
		basic.Point{},
		40,
		basic.Size{W: size.W, H: 80},
		basic.Center,
		basic.Center,
		[]components.Widget{
			backButton,
			newPlayer,
		},
	)

	children := []components.Widget{spacer, title}

	children = append(children, listArea)
	children = append(children, buttonRow)

	return components.NewColumn(
		basic.Point{X: 0, Y: 0},
		16,
		size,
		basic.Start,
		basic.Center,
		children,
	)
}

func (s *SelectProfileScene) buildList(size basic.Size) components.Widget {
	items := []components.Widget{}
	start := s.offset
	end := start + 4
	if end > len(s.profiles) {
		end = len(s.profiles)
	}
	for i := start; i < end; i++ {
		p := s.profiles[i]
		name := p.Username

		profileCopy := p

		nameBtn := components.NewButton(
			basic.Point{},
			basic.Size{W: size.W * 0.5, H: 50},
			name,
			colors.PlayerInput,
			nil,
			func(b *components.Button) {
				if SwitchTo != nil {
					SwitchTo(NewProfileSceneWithProfile(&profileCopy))
				}
			},
		)

		iconSize := basic.Size{W: 35, H: 35}
		deleteIcon := newDeleteProfileIcon(iconSize, func() {
			_ = service.RemoveProfile(profileCopy.Username)
			s.profiles = service.GetProfiles()

			if len(s.profiles) == 0 {
				s.offset = 0
			} else {
				maxOffset := 0
				if len(s.profiles) > 4 {
					maxOffset = len(s.profiles) - 4
				}
				if s.offset > maxOffset {
					s.offset = maxOffset
				}
			}

			s.root = s.buildUI(s.screenSize)
		})

		playIcon := newPlayProfileIcon(iconSize, func() {
			if SwitchTo != nil {
				SwitchTo(NewPlacementSceneWithProfile(&profileCopy))
			}
		})

		var rowChildren []components.Widget
		switch {
		case deleteIcon != nil && playIcon != nil:
			rowChildren = []components.Widget{deleteIcon, nameBtn, playIcon}
		case deleteIcon != nil:
			rowChildren = []components.Widget{deleteIcon, nameBtn}
		case playIcon != nil:
			rowChildren = []components.Widget{nameBtn, playIcon}
		default:
			rowChildren = []components.Widget{nameBtn}
		}

		row := components.NewRow(
			basic.Point{},
			0,
			size,
			basic.Center,
			basic.Center,
			rowChildren,
		)

		items = append(items, row)
	}

	return components.NewColumn(
		basic.Point{X: 1, Y: 0},
		5,
		size,
		basic.Center,
		basic.Center,
		items,
	)
}
