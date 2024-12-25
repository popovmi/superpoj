package color

import (
	"fmt"
	"image/color"
)

//go:generate msgp

type RGBA [4]byte

func (r *RGBA) ExtensionType() int8 { return 99 }

func (r *RGBA) Len() int { return 4 }

func (r *RGBA) MarshalBinaryTo(b []byte) error {
	copy(b, (*r)[:])
	return nil
}

func (r *RGBA) UnmarshalBinary(b []byte) error {
	copy((*r)[:], b)
	return nil
}

func (r *RGBA) MarshalJSON() ([]byte, error) {
	b := *r
	return []byte(fmt.Sprintf("[%d, %d, %d, %d]", b[0], b[1], b[2], b[3])), nil
}

func (r *RGBA) ToColorRGBA() color.RGBA {
	return color.RGBA{R: r[0], G: r[1], B: r[2], A: r[3]}
}
