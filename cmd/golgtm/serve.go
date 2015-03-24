package main

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hiroosak/golgtm"
)

func serve(addr string) {
	log.Println("start http server(" + addr + ")")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	src := q.Get("src")
	if src == "" {
		badRequestHandler(w, r)
		return
	}

	bs, err := fetchURL(src)
	if err != nil {
		panic(err)
	}

	f := bytes.NewBuffer(bs)
	img, t, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	lgtm, err := golgtm.NewLgtm()
	if err != nil {
		panic(err)
	}

	max := img.Bounds().Max
	textSize := lgtm.Text.Rectangle().Max

	x := q.Get("x")
	if x != "" {
		px, _ := strconv.ParseInt(x, 32, 10)
		lgtm.X = int(px)
	}
	if lgtm.X == 0 {
		lgtm.X = (max.X / 2) - (textSize.X / 2)
	}

	y := q.Get("y")
	if y != "" {
		py, _ := strconv.ParseInt(y, 32, 10)
		lgtm.Y = int(py)
	}
	if lgtm.Y == 0 {
		lgtm.Y = max.Y - (textSize.Y / 2)
	}

	lgtm.Text.SetFontColor(q.Get("color"))

	result, err := lgtm.Convert(img)
	if err != nil {
		panic(err)
	}

	b := bytes.NewBufferString("")

	switch t {
	case "png":
		png.Encode(b, result)
	case "jpeg":
		jpeg.Encode(b, result, nil)
	case "gif":
		gif.Encode(b, result, nil)
	}

	w.Header().Set("Content-Type", "image/"+t)
	b.WriteTo(w)
}

func badRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Bad Request"))
}

func fetchURL(src string) ([]byte, error) {
	cli := &http.Client{Timeout: 1 * time.Second}
	resp, err := cli.Get(src)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// sizeが2M以上なら対象外
	// imageでなければerr

	return ioutil.ReadAll(resp.Body)
}
