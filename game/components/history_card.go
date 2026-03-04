package components

import (
	"fmt"
	"image/color"

	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

// como estava com preguiça fiz tudo e mandei o gpt modularizar (sinceramente nao acho que melhorou aslkdnalsdns)
type HistoryCard struct {
	pos, currentPos basic.Point
	body            StylableWidget
	entity.MatchResult
}

func NewHistoryCard(pos basic.Point, size basic.Size, result entity.MatchResult) *HistoryCard {
	return &HistoryCard{
		pos:  pos,
		body: buildCardContainer(pos, size, result),
	}
}

/* =======================
   Estrutura Principal
======================= */

func buildCardContainer(pos basic.Point, size basic.Size, result entity.MatchResult) StylableWidget {
	return NewContainer(
		pos, size, 20,
		colors.SeaCyan, basic.Center, basic.Center,
		buildCardContent(size, result),
	)
}

func buildCardContent(size basic.Size, result entity.MatchResult) Widget {
	return NewContainer(
		basic.Point{}, size, 0,
		colors.Transparent, basic.Center, basic.Center,
		NewColumn(
			basic.Point{}, 10, size,
			basic.Center, basic.Center,
			[]Widget{
				buildHeader(size, result),
				buildStatsSection(size, result),
			},
		),
	)
}

/* =======================
   Header
======================= */

func buildHeader(size basic.Size, result entity.MatchResult) Widget {
	rowSize := basic.Size{
		W: size.W * 0.9,
		H: size.H * 0.2,
	}

	title, resultColor := resolveResultLabel(result)
	scoreColor := resolveScoreColor(result.Score)

	return NewContainer(
		basic.Point{}, rowSize, 0,
		colors.Transparent, basic.Center, basic.Center,
		NewRow(
			basic.Point{}, 10, rowSize,
			basic.Start, basic.Center,
			[]Widget{
				mustImage("assets/icons/skull.png", 45, 45),
				NewText(basic.Point{}, title, resultColor, 35),
				sideSpacer(450),
				mustImage("assets/icons/star.png", 40, 40),
				NewText(basic.Point{}, "SCORE: ", colors.White, 35),
				NewText(basic.Point{}, fmt.Sprintf("%03d", result.Score), scoreColor, 35),
			},
		),
	)
}

/* =======================
   Estatísticas
======================= */

func buildStatsSection(size basic.Size, result entity.MatchResult) Widget {
	rowSize := basic.Size{
		W: size.W * 0.98,
		H: size.H * 0.65,
	}

	internalSize := basic.Size{
		W: size.W * 0.47,
		H: size.H * 0.5,
	}

	return NewContainer(
		basic.Point{}, rowSize, 12,
		colors.NightBlue, basic.Center, basic.Center,
		NewRow(
			basic.Point{}, 0, rowSize,
			basic.Center, basic.Center,
			[]Widget{
				buildLeftStats(internalSize, result),
				buildRightStats(internalSize, result),
			},
		),
	)
}

func buildLeftStats(size basic.Size, result entity.MatchResult) Widget {
	iconRowSize := basic.Size{W: size.W, H: 40}
	_, resultColor := resolveResultLabel(result)

	return NewContainer(
		basic.Point{}, size, 0,
		colors.Transparent, basic.Center, basic.Center,
		NewColumn(
			basic.Point{}, 5, size,
			basic.Start, basic.Center,
			[]Widget{
				buildIconRow("assets/icons/shot.png", "TIROS............................", fmt.Sprintf("%02d", result.PlayerShots), iconRowSize, resultColor),
				buildIconRow("assets/icons/target.png", "HITS.............",
					fmt.Sprintf("%02d (%02.2f%%)",
						result.Hits,
						safeHitPercent(result.Hits, result.PlayerShots)),
					iconRowSize,
					resultColor),
				buildIconRow("assets/icons/eye.png", "MAIOR SEQUENCIA......",
					fmt.Sprintf("%02d", result.HigherHitSequence),
					iconRowSize,
					resultColor),
			},
		),
	)
}

func buildRightStats(size basic.Size, result entity.MatchResult) Widget {
	iconRowSize := basic.Size{W: size.W, H: 40}
	_, resultColor := resolveResultLabel(result)

	return NewContainer(
		basic.Point{}, size, 0,
		colors.Transparent, basic.Start, basic.Start,
		NewColumn(
			basic.Point{}, 5, size,
			basic.Start, basic.Center,
			[]Widget{
				buildIconRow("assets/icons/clock.png", "DURAÇÃO...............",
					result.FormattedDuration(),
					iconRowSize,
					resultColor),
				buildIconRow("assets/icons/anchor.png", "NAVIOS PERDIDOS...",
					fmt.Sprintf("%02d/06", result.LostShips),
					iconRowSize,
					resultColor),
			},
		),
	)
}

/* =======================
   Utilitários
======================= */

func resolveResultLabel(result entity.MatchResult) (string, color.Color) {
	if result.Win {
		return "VITORIA", colors.GoldMedal
	}
	return "DERROTA", colors.Red
}

func resolveScoreColor(score int) color.Color {
	switch {
	case score >= 1000:
		return colors.Lighten(colors.GoldMedal, 0.3)
	case score >= 500:
		return colors.Lighten(colors.SilverMedal, 0.3)
	default:
		return colors.Lighten(colors.BronzeMedal, 0.3)
	}
}

func buildIconRow(iconPath, label, value string, size basic.Size, col color.Color) Widget {
	row, err := NewIconRow(iconPath, label, value, size, basic.Point{}, col)
	if err != nil {
		panic(err)
	}
	return row
}

func mustImage(path string, w, h float32) Widget {
	img, err := NewImage(path, basic.Point{}, basic.Size{W: w, H: h})
	if err != nil {
		panic(err)
	}
	return img
}

func sideSpacer(width float32) Widget {
	return NewContainer(
		basic.Point{},
		basic.Size{W: width, H: 1},
		0,
		colors.Transparent,
		basic.Center,
		basic.Center,
		nil,
	)
}

func safeHitPercent(hits, shots int) float64 {
	if shots == 0 {
		return 0
	}
	return (float64(hits) / float64(shots)) * 100
}

/* =======================
   Métodos do Card
======================= */

func (h *HistoryCard) GetPos() basic.Point {
	return h.pos
}

func (h *HistoryCard) SetPos(point basic.Point) {
	h.pos = point
}

func (h *HistoryCard) GetSize() basic.Size {
	return h.body.GetSize()
}

func (h *HistoryCard) Update(offset basic.Point) {
	h.currentPos = h.pos.Add(offset)
	h.body.Update(h.currentPos)
}

func (h *HistoryCard) Draw(screen *ebiten.Image) {
	h.body.Draw(screen)
}
