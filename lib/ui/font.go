package ui

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	//go:embed assets/RobotoMono-Regular.ttf
	robotoMonoTtf []byte

	fontFaceSource *text.GoTextFaceSource
	FontFace22     *text.GoTextFace
	FontFace18     *text.GoTextFace
	FontFace16     *text.GoTextFace
)

func InitFont() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(robotoMonoTtf))
	if err != nil {
		log.Fatal(err)
	}

	fontFaceSource = s
	FontFace22 = &text.GoTextFace{
		Source: fontFaceSource,
		Size:   22,
	}
	FontFace18 = &text.GoTextFace{
		Source: fontFaceSource,
		Size:   18,
	}
	FontFace16 = &text.GoTextFace{
		Source: fontFaceSource,
		Size:   16,
	}
}
