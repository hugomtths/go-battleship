package components

import (
	"image/color"
	"strings"
	"unicode/utf8"

	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type TextWrap struct {
	pos, currentPos basic.Point
	color           color.Color
	rawText         string
	lines           []string
	maxWidth        float32
	fontSize        int
	face            font.Face
	size            basic.Size
	lineHeight      float32
}

func NewTextWrap(pos basic.Point, str string, c color.Color, fontSize int, maxWidth float32) *TextWrap {
	tw := &TextWrap{
		pos:      pos,
		color:    c,
		rawText:  str,
		maxWidth: maxWidth,
		fontSize: fontSize,
	}
	tw.face = createFace(float64(fontSize))
	tw.buildLines()
	return tw
}

func (tw *TextWrap) buildLines() {
	tw.lines = []string{}
	if tw.face == nil || tw.rawText == "" {
		tw.size = basic.Size{}
		return
	}

	words := strings.Fields(tw.rawText)
	if len(words) == 0 {
		tw.size = basic.Size{}
		return
	}

	currentLine := ""
	for _, word := range words {
		candidate := word
		if currentLine != "" {
			candidate = currentLine + " " + word
		}

		lineWidth := measureTextWidth(tw.face, candidate)
		if lineWidth > tw.maxWidth && currentLine != "" {
			tw.lines = append(tw.lines, currentLine)
			currentLine = word
		} else {
			currentLine = candidate
		}
	}
	if currentLine != "" {
		tw.lines = append(tw.lines, currentLine)
	}

	tw.recalcSize()
}

func (tw *TextWrap) recalcSize() {
	if len(tw.lines) == 0 {
		tw.size = basic.Size{}
		return
	}

	metrics := tw.face.Metrics()
	tw.lineHeight = float32(metrics.Height.Round())

	maxW := float32(0)
	for _, line := range tw.lines {
		w := float32(measureTextWidth(tw.face, line))
		if w > maxW {
			maxW = w
		}
	}

	tw.size = basic.Size{
		W: maxW,
		H: tw.lineHeight * float32(len(tw.lines)),
	}
}

func measureTextWidth(face font.Face, str string) float32 {
	if face == nil || str == "" {
		return 0
	}
	advance := font.MeasureString(face, str)

	if advance == 0 {
		total := 0
		prev := rune(-1)
		for _, r := range str {
			if prev >= 0 {
				if kern, ok := face.(interface {
					Kern(rune, rune) int64
				}); ok {
					total += int(kern.Kern(prev, r))
				}
			}
			a, ok := face.GlyphAdvance(r)
			if ok {
				total += int(a)
			}
			prev = r
			_ = utf8.RuneLen(r)
		}
		return float32(total >> 6)
	}

	return float32(advance >> 6)
}

func (tw *TextWrap) GetPos() basic.Point  { return tw.pos }
func (tw *TextWrap) SetPos(p basic.Point) { tw.pos = p }
func (tw *TextWrap) GetSize() basic.Size  { return tw.size }

func (tw *TextWrap) Update(offset basic.Point) {
	tw.currentPos = tw.pos.Add(offset)
}

func (tw *TextWrap) Draw(screen *ebiten.Image) {
	if tw.face == nil || len(tw.lines) == 0 {
		return
	}

	metrics := tw.face.Metrics()
	ascent := float32(metrics.Ascent.Round())

	for i, line := range tw.lines {
		lineW := measureTextWidth(tw.face, line)

		// cada linha é centralizada em relação à linha mais longa (tw.size.W)
		offsetX := (tw.size.W - lineW) / 2

		x := int(tw.currentPos.X + offsetX)
		y := int(tw.currentPos.Y + ascent + tw.lineHeight*float32(i))
		text.Draw(screen, line, tw.face, x, y, tw.color)
	}
}
