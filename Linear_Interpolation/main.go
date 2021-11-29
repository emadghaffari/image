package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var uColor = color.RGBA{238, 164, 127, 0xff}
var dColor = color.RGBA{0, 83, 156, 0xff}

const outFilename = "image.png"
const wid = 2160
const hei = 1920

func main() {

	upLeft := image.Point{0, 0}
	downRight := image.Point{wid, hei}

	img := image.NewRGBA(image.Rectangle{upLeft, downRight})

	for y := 0; y < hei; y++ {
		uRatio := float32(hei-y) / float32(hei)
		dRatio := float32(y) / float32(hei)
		rowR := uint8(uRatio*float32(uColor.R) + dRatio*float32(dColor.R))
		rowG := uint8(uRatio*float32(uColor.G) + dRatio*float32(dColor.G))
		rowB := uint8(uRatio*float32(uColor.B) + dRatio*float32(dColor.B))

		rowColor := color.RGBA{rowR, rowG, rowB, 0xff}

		for x := 0; x < wid; x++ {
			img.Set(x, y, rowColor)
		}
	}

	outFile, _ := os.Create(outFilename)
	png.Encode(outFile, img)
}
