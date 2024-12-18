package image

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math/rand"
)

func GenerateRandomImage(width int, height int) ([]byte, error) {
	buf := new(bytes.Buffer)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	clr := color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255}
	draw.Draw(img, img.Bounds(), image.NewUniform(clr), image.Point{}, draw.Src)
	err := jpeg.Encode(buf, img, nil)

	return buf.Bytes(), err
}
