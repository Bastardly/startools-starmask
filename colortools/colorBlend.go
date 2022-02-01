package colortools

func ChannelBlendByProcentage(procentage float64, baseChannel, mixChannel uint32) uint32 {
	diff := float64(mixChannel) - float64(baseChannel)

	return baseChannel + uint32(diff*procentage)
}
