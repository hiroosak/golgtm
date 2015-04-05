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

	http.HandleFunc("/form", formHandler)
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

	fontSize := q.Get("fontSize")
	if fontSize != "" {
		fs, _ := strconv.ParseInt(fontSize, 64, 10)
		lgtm.Text.FontSize = float64(fs)
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

func formHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8" />
		<title>go lgtm</title>
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.4/css/bootstrap.min.css">
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.4/css/bootstrap-theme.min.css">
		<style>
		.main {
			margin-top: 120px;
		}
		</style>
	</head>
	<body>
		<nav class="navbar navbar-inverse navbar-fixed-top">
			<div class="container">
				<div class="navbar-header">
					<a class="navbar-brand" href="#">Go LGTM</a>
				</div>
			</div>
		</nav>
		<div class="container">
			<div class="main">
				<div class="row"> 
					<div class="col-md-10">
						<img src="/?src=" id="image" height="500">
					</div>
					<div class="col-md-2">
						<form action="GET">
							<div class="form-group">
								<label for="src">src</label>
								<input type="text" id="src" name="src" placeholder="http://...">
							</div>
							<div class="form-group">
								<label for="fontSize">font size</label>
								<input type="text" id="fontSize" name="fontSize" placeholder="94">
							</div>
							<div class="form-group">
								<label for="color">color</label>
								<input type="color" id="color" value="#ffffff" name="color">
							</div>
							<div class="form-group">
								<label for="x">x</label>
								<input type="text" id="x" name="x" placeholder="100">
							</div>
							<div class="form-group">
								<label for="y">y</label>
								<input type="text" id="y" name="y" placeholder="100">
							</div>
						</form>
					</div>
				</div>
			</div>
		</div>
		<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.4/js/bootstrap.min.js"></script>
	</body>
	</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
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

	return ioutil.ReadAll(resp.Body)
}
