package golgtm

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"testing"
)

func TestLgtmConvert(t *testing.T) {
	f, _ := os.Open("cat.png")
	defer f.Close()
	img, ty, _ := image.Decode(f)
	fmt.Println(ty)

	lgtm, err := NewLgtm()
	if err != nil {
		t.Fatalf("error %v", err)
	}

	result, err := lgtm.Convert(img)
	if err != nil {
		t.Fatalf("error %v", err)
	}

	out, _ := os.Create("cat_exist.png")
	defer out.Close()

	// write new image to file
	png.Encode(out, result)
}
