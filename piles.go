package main

import "github.com/hajimehoshi/ebiten/v2"

type Hand struct {
	cards []*Card
}

func (h *Hand) Update() {

}

func (h *Hand) Draw(screen *ebiten.Image) {
	offset := (ScreenWidth - ArtSmallWidth*len(h.cards)) / 2
	for i, c := range h.cards {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(offset+i*ArtSmallWidth), ScreenHeight-ArtSmallWidth)
		screen.DrawImage(c.artSmall, op)
	}
}

// given logical screen pixel location x,y returns the detailed art if there is a card in hand there
func (h *Hand) In(x, y int) *ebiten.Image {
	localX := x - ((ScreenWidth - ArtSmallWidth*len(h.cards)) / 2)
	localY := y - (ScreenHeight - ArtSmallWidth)
	if localX > 0 && localX < ArtSmallWidth*len(h.cards) && localY > 0 && localY < ArtSmallWidth {
		return h.cards[localX/ArtSmallWidth].artBig
	}
	return nil
}
