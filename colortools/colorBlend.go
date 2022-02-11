package colortools

func ChannelBlendByProcentage(procentage float64, baseChannel, mixChannel uint32) uint32 {
	diff := float64(mixChannel) - float64(baseChannel)

	return baseChannel + uint32(diff*procentage)
}

type Color struct {
	R uint32
	G uint32
	B uint32
}

func MixColors(firstColor, oppositeColor Color, steps int) Color {
	procentage := 1 / float64(steps)

	return Color{
		R: ChannelBlendByProcentage(procentage, firstColor.R, oppositeColor.R),
		G: ChannelBlendByProcentage(procentage, firstColor.G, oppositeColor.G),
		B: ChannelBlendByProcentage(procentage, firstColor.B, oppositeColor.B),
	}
}
