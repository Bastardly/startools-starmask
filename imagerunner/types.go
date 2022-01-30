package imagerunner

type Pixel struct {
	R                           uint32
	G                           uint32
	B                           uint32
	A                           uint32
	brightness                  float32 // WCAG brightness made from 8bit estimate
	HasContrastChangeHorizontal bool
	HasContrastChangeVertical   bool
	HasBeenExplored             bool
	IsGlow                      bool
	Rejected                    bool // if not star, we reject all bright pixels in contact with this.
	starRadiusHorizontal        int
	starRadiusVertical          int
}

type Settings struct {
	maxStarSizeInPx     int
	maxStarGlowInPx     int
	wcagContrastMinimum float32 // The highter the star contrast is needed for detection.
}

type Pixels = [][]Pixel

type Store struct {
	Width int
	Height int
	settings Settings
	Pixels   [][]Pixel

}
