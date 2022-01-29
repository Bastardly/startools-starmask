package imagerunner

import (
	"image"
	"starkiller/colortools"
)

const bitDivider = 257 // 16 bit to 8 bit divider

func rgbaToPixel(R uint32, G uint32, B uint32, A uint32) Pixel {
	r := uint8(R / bitDivider)
	g := uint8(G / bitDivider)
	b := uint8(B / bitDivider)
	brightness := colortools.GetRGBBrightNess(r, g, b)

	return Pixel{R, G, B, A, brightness, false, false, false, false, false, -1, -1}
}

// Get the bi-dimensional pixel array and image width and height
func getPixels(img image.Image) ([][]Pixel, int, int) {

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			R, G, B, A := img.At(x, y).RGBA()
			column := rgbaToPixel(R, G, B, A)
			row = append(row, column)
		}
		pixels = append(pixels, row)
	}

	return pixels, width, height
}
