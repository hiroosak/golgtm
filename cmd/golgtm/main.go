package main

import (
	"flag"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/hiroosak/golgtm"
)

var (
	src  string
	dst  string
	text string
	x    int
	y    int
	addr string
)

func init() {
	flag.StringVar(&addr, "http", ":8000", "HTTP service address (e.g., ':8000'")

	flag.StringVar(&src, "src", "", "source file")
	flag.StringVar(&dst, "dst", "", "dest file")
	flag.StringVar(&text, "text", "LGTM", "insert text")
	flag.IntVar(&x, "x", 0, "x text position")
	flag.IntVar(&y, "y", 0, "y text position")
}

func main() {
	flag.Parse()

	if addr != "" {
		serve(addr)
	} else {
		cmd()
	}
}

// cmd is command line
func cmd() {
	if src == "" || dst == "" || text == "" {
		flag.Usage()
		return
	}

	if _, err := os.Stat(src); err != nil {
		flag.Usage()
		return
	}

	f, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, t, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	lgtm, err := golgtm.NewLgtm()
	if err != nil {
		panic(err)
	}

	result, err := lgtm.Convert(img)
	if err != nil {
		panic(err)
	}

	out, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	switch t {
	case "png":
		png.Encode(out, result)
	case "jpeg":
		jpeg.Encode(out, result, nil)
	case "gif":
		gif.Encode(out, result, nil)
	}
}
