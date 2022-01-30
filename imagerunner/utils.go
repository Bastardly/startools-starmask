package imagerunner

import (
	"starkiller/colortools"
)

func rgbaToPixel(R uint32, G uint32, B uint32, A uint32) Pixel {
	r := uint8(R / bitDivider)
	g := uint8(G / bitDivider)
	b := uint8(B / bitDivider)
	brightness := colortools.GetRGBBrightNess(r, g, b)

	return Pixel{R, G, B, A, brightness, false, false, false, false, false, -1, -1, -1, -1}
}
