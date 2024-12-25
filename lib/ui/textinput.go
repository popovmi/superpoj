package ui

import (
	"image"
	"image/color"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/exp/textinput"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"wars/lib/game"
)

type TextField struct {
	bounds     image.Rectangle
	multilines bool
	field      textinput.Field
	fontFace   *text.GoTextFace
}

func NewTextField(
	bounds image.Rectangle, multilines bool, fontFace *text.GoTextFace,
) *TextField {
	return &TextField{
		bounds:     bounds,
		multilines: multilines,
		fontFace:   fontFace,
	}
}

func (t *TextField) Contains(x, y int) bool {
	return image.Pt(x, y).In(t.bounds)
}

func (t *TextField) SetSelectionStartByCursorPosition(x, y int) bool {
	idx, ok := t.textIndexByCursorPosition(x, y)
	if !ok {
		return false
	}
	t.field.SetSelection(idx, idx)
	return true
}

func (t *TextField) textIndexByCursorPosition(x, y int) (int, bool) {
	if !t.Contains(x, y) {
		return 0, false
	}

	x -= t.bounds.Min.X
	y -= t.bounds.Min.Y
	px, py := t.padding()
	x -= px
	y -= py
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	lineSpacingInPixels := int(t.fontFace.Metrics().HLineGap + t.fontFace.Metrics().HAscent + t.fontFace.Metrics().HDescent)
	var nlCount int
	var lineStart int
	var prevAdvance float64
	txt := t.field.Text()
	for i, r := range txt {
		var x0, x1 int
		currentAdvance := text.Advance(txt[lineStart:i], t.fontFace)
		if lineStart < i {
			x0 = int((prevAdvance + currentAdvance) / 2)
		}
		if r == '\n' {
			x1 = int(math.MaxInt32)
		} else if i < len(txt) {
			nextI := i + 1
			for !utf8.ValidString(txt[i:nextI]) {
				nextI++
			}
			nextAdvance := text.Advance(txt[lineStart:nextI], t.fontFace)
			x1 = int((currentAdvance + nextAdvance) / 2)
		} else {
			x1 = int(currentAdvance)
		}
		if x0 <= x && x < x1 && nlCount*lineSpacingInPixels <= y && y < (nlCount+1)*lineSpacingInPixels {
			return i, true
		}
		prevAdvance = currentAdvance

		if r == '\n' {
			nlCount++
			lineStart = i + 1
			prevAdvance = 0
		}
	}

	return len(txt), true
}

func (t *TextField) Focus() {
	t.field.Focus()
}

func (t *TextField) Blur() {
	t.field.Blur()
}

func (t *TextField) Update() error {
	if !t.IsFocused() {
		return nil
	}

	fieldVal := t.field.Text()
	x, y := t.bounds.Min.X, t.bounds.Min.Y
	cx, cy := t.cursorPos()
	px, py := t.padding()
	x += cx + px
	y += cy + py + int(t.fontFace.Metrics().HAscent)
	if len(fieldVal) < warsgame.MaxTextLength {
		handled, err := t.field.HandleInput(x, y)
		if err != nil {
			return err
		}
		if handled {
			return nil
		}
	}

	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter):
		if t.multilines {
			selectionStart, selectionEnd := t.field.Selection()
			fieldVal = fieldVal[:selectionStart] + "\n" + fieldVal[selectionEnd:]
			selectionStart += len("\n")
			selectionEnd = selectionStart
			t.field.SetTextAndSelection(fieldVal, selectionStart, selectionEnd)
		}

	case inpututil.IsKeyJustPressed(ebiten.KeyBackspace):
		selectionStart, selectionEnd := t.field.Selection()
		if selectionStart != selectionEnd {
			fieldVal = fieldVal[:selectionStart] + fieldVal[selectionEnd:]
		} else if selectionStart > 0 {
			// TODO: Remove a grapheme instead of a code point.
			_, l := utf8.DecodeLastRuneInString(fieldVal[:selectionStart])
			fieldVal = fieldVal[:selectionStart-l] + fieldVal[selectionEnd:]
			selectionStart -= l
		}
		selectionEnd = selectionStart
		t.field.SetTextAndSelection(fieldVal, selectionStart, selectionEnd)

	case inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		selectionStart, _ := t.field.Selection()
		if selectionStart > 0 {
			// TODO: Remove a grapheme instead of a code point.
			_, l := utf8.DecodeLastRuneInString(fieldVal[:selectionStart])
			selectionStart -= l
		}
		t.field.SetTextAndSelection(fieldVal, selectionStart, selectionStart)

	case inpututil.IsKeyJustPressed(ebiten.KeyRight):
		_, selectionEnd := t.field.Selection()
		if selectionEnd < len(fieldVal) {
			// TODO: Remove a grapheme instead of a code point.
			_, l := utf8.DecodeRuneInString(fieldVal[selectionEnd:])
			selectionEnd += l
		}
		t.field.SetTextAndSelection(fieldVal, selectionEnd, selectionEnd)
	}

	if !t.multilines {
		newVal := strings.ReplaceAll(fieldVal, "\n", "")
		if newVal != fieldVal {
			selectionStart, selectionEnd := t.field.Selection()
			selectionStart -= strings.Count(fieldVal[:selectionStart], "\n")
			selectionEnd -= strings.Count(fieldVal[:selectionEnd], "\n")
			t.field.SetSelection(selectionStart, selectionEnd)
		}
	}

	return nil
}

func (t *TextField) IsFocused() bool {
	return t.field.IsFocused()
}

func (t *TextField) Value() string {
	return t.field.Text()
}

func (t *TextField) cursorPos() (int, int) {
	var nlCount int
	lastNLPos := -1
	txt := t.field.TextForRendering()
	selectionStart, _ := t.field.Selection()
	if s, _, ok := t.field.CompositionSelection(); ok {
		selectionStart += s
	}
	txt = txt[:selectionStart]
	for i, r := range txt {
		if r == '\n' {
			nlCount++
			lastNLPos = i
		}
	}

	txt = txt[lastNLPos+1:]
	x := int(text.Advance(txt, t.fontFace))
	y := nlCount * int(t.fontFace.Metrics().HLineGap+t.fontFace.Metrics().HAscent+t.fontFace.Metrics().HDescent)
	return x, y
}

func (t *TextField) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(
		screen,
		float32(t.bounds.Min.X),
		float32(t.bounds.Min.Y),
		float32(t.bounds.Dx()),
		float32(t.bounds.Dy()),
		color.Black, false,
	)
	var clr color.Color = color.Gray{Y: 127}
	if t.field.IsFocused() {
		clr = color.White
	}
	vector.StrokeRect(
		screen, float32(t.bounds.Min.X), float32(t.bounds.Min.Y),
		float32(t.bounds.Dx()), float32(t.bounds.Dy()), 1, clr, false,
	)

	px, py := t.padding()
	selectionStart, _ := t.field.Selection()
	if t.field.IsFocused() && selectionStart >= 0 {
		x, y := t.bounds.Min.X, t.bounds.Min.Y
		cx, cy := t.cursorPos()
		x += px + cx
		y += py + cy
		h := int(t.fontFace.Metrics().HLineGap + t.fontFace.Metrics().HAscent + t.fontFace.Metrics().HDescent)
		vector.StrokeLine(
			screen, float32(x), float32(y), float32(x), float32(y+h), 1,
			color.Gray{Y: 127}, false,
		)
	}

	tx := t.bounds.Min.X + px
	ty := t.bounds.Min.Y + py
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(tx), float64(ty))
	op.ColorScale.ScaleWithColor(color.Gray{Y: 177})
	op.LineSpacing = t.fontFace.Metrics().HLineGap + t.fontFace.Metrics().HAscent + t.fontFace.Metrics().HDescent
	text.Draw(screen, t.field.TextForRendering(), t.fontFace, op)
}

func (t *TextField) padding() (int, int) {
	m := t.fontFace.Metrics()
	return 4, (warsgame.TextFieldHeight - int(m.HLineGap+m.HAscent+m.HDescent)) / 2
}
