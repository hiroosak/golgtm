package golgtm

import (
	"image"
	"image/draw"

	"code.google.com/p/freetype-go/freetype"
)

type Lgtm struct {
	Text *Text
	X    int
	Y    int
}

func NewLgtm() (*Lgtm, error) {
	str, err := NewText("LGTM")
	if err != nil {
		return nil, err
	}
	return &Lgtm{Text: str}, nil
}

func (l *Lgtm) Convert(dstImage image.Image) (image.Image, error) {
	ctx := freetype.NewContext()

	src := image.NewUniform(l.Text.FontColor)
	ctx.SetSrc(src)

	ctx.SetFont(l.Text.Font)
	ctx.SetFontSize(l.Text.FontSize)
	ctx.SetHinting(freetype.FullHinting)

	dst := image.NewRGBA(dstImage.Bounds())
	bounds := dst.Bounds()
	draw.Draw(dst, bounds, dstImage, image.ZP, draw.Src)

	ctx.SetDst(dst)
	ctx.SetClip(bounds)

	if _, err := ctx.DrawString(l.Text.String(), freetype.Pt(l.X, l.Y)); err != nil {
		return nil, err
	}

	return dst, nil
}
