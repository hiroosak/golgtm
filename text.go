package golgtm

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"strconv"
	"unicode/utf8"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
)

const (
	DefaultFontSize = 94
	DefaultFontType = "Arial"
)

type Text struct {
	Font      *truetype.Font
	FontSize  float64
	Scale     int32
	Dpi       float64
	FontColor color.Color
	Hinting   freetype.Hinting
	Message   string
}

func NewText(str string) (*Text, error) {
	t := &Text{
		FontSize:  DefaultFontSize,
		FontColor: color.Black,
		Message:   str,
		Dpi:       72,
		Hinting:   freetype.NoHinting,
	}
	if err := t.FontType(DefaultFontType); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Text) String() string {
	return s.Message
}

func (s *Text) Rectangle() image.Rectangle {
	n := utf8.RuneCountInString(s.Message)
	size := int(s.FontSize)

	return image.Rectangle{
		Min: image.Point{},
		Max: image.Point{
			X: (size / 2) * n,
			Y: size,
		},
	}
}

func (s *Text) SetFontColor(colorName string) {
	switch colorName {
	case "blue":
		s.FontColor = color.RGBA{0, 0, 255, 255}
	case "green":
		s.FontColor = color.RGBA{0, 255, 0, 255}
	case "red":
		s.FontColor = color.RGBA{255, 0, 0, 255}
	case "white":
		s.FontColor = color.White
	case "black":
		s.FontColor = color.Black
	default:
		r, g, b := hexToRGB(colorName)
		s.FontColor = color.RGBA{r, g, b, 255}
	}
}

func hexToRGB(hex string) (uint8, uint8, uint8) {
	if len(hex) == 6 {
		if rgb, err := strconv.ParseUint(hex, 16, 32); err == nil {
			return uint8(rgb >> 16), uint8((rgb >> 8) & 0xFF), uint8(rgb & 0xFF)
		}
	}
	return 0, 0, 0
}

func (s *Text) FontType(ttf string) error {
	path := fmt.Sprintf("%v%v.ttf", config.FontBasePath, ttf)
	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if s.Font, err = freetype.ParseFont(bs); err != nil {
		return err
	}
	return nil
}
