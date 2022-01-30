package imagerunner

func (p Pixel) IsValid() bool {
	return p.HasContrastChangeHorizontal && p.HasContrastChangeVertical
}
