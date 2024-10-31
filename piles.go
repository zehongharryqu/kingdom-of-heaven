package main

import "github.com/hajimehoshi/ebiten/v2"

type PlayerCards struct {
	hand, deck, discard []*Card
}

func (pc *PlayerCards) Update() {

}

func (pc *PlayerCards) Draw(screen *ebiten.Image) {
	// deck
	if len(pc.deck) > 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(DeckPileX, DiscardDeckPileY)
		// TODO: back of card
		screen.DrawImage(pc.deck[len(pc.deck)-1].artSmall, op)
	}
	// discard
	if len(pc.discard) > 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(DiscardPileX, DiscardDeckPileY)
		screen.DrawImage(pc.discard[len(pc.discard)-1].artSmall, op)
	}
	// hand
	offset := (ScreenWidth - ArtSmallWidth*len(pc.hand)) / 2
	for i, c := range pc.hand {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(offset+i*ArtSmallWidth), ScreenHeight-ArtSmallWidth)
		screen.DrawImage(c.artSmall, op)
	}
}

// given logical screen pixel location x,y returns the detailed art if there is a card in hand there
func (pc *PlayerCards) In(x, y int) *ebiten.Image {
	localX := x - ((ScreenWidth - ArtSmallWidth*len(pc.hand)) / 2)
	localY := y - (ScreenHeight - ArtSmallWidth)
	if localX > 0 && localX < ArtSmallWidth*len(pc.hand) && localY > 0 && localY < ArtSmallWidth {
		return pc.hand[localX/ArtSmallWidth].artBig
	}
	return nil
}

type VersePile struct {
	// which card this is a pile of
	c *Card
	// how many are left in the pile
	n int
}

type Kingdom struct {
	v        []*VersePile
	released []*Card
}

// coordinates to draw kingdom
var (
	KingdomPileX = [...]int{460, 340, 340, 400, 580, 520, 580, 340, 400, 460, 520, 580, 340, 400, 460, 520, 580}
	KingdomPileY = [...]int{70, 10, 70, 70, 10, 70, 70, 130, 130, 130, 130, 130, 190, 190, 190, 190, 190}
)

const (
	ReleasedPileX    = 460
	ReleasedPileY    = 10
	DiscardPileX     = 520
	DeckPileX        = 580
	DiscardDeckPileY = 370
)

// draws the kingdom piles and released pile
func (k *Kingdom) Draw(screen *ebiten.Image) {
	for i, v := range k.v {
		if v.n > 0 {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(KingdomPileX[i]), float64(KingdomPileY[i]))
			screen.DrawImage(v.c.artSmall, op)
		}
	}
	if n := len(k.released); n > 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(ReleasedPileX, ReleasedPileY)
		screen.DrawImage(k.released[n-1].artSmall, op)
	}
}

// given logical screen pixel location x,y returns the detailed art if there is a kingdom card there
func (k *Kingdom) In(x, y int) *ebiten.Image {
	for i, v := range k.v {
		if x > KingdomPileX[i] && x < KingdomPileX[i]+ArtSmallWidth && y > KingdomPileY[i] && y < KingdomPileY[i]+ArtSmallWidth {
			return v.c.artBig
		}
	}
	return nil
}

// create a new kingdom given the 10 verses and number of players
func InitKingdom(verses []*Card, n int) *Kingdom {
	// starting amounts from the dominion wiki gameplay article
	var startingStudy, startingPrayer, startingDevotion, startingParable, startingSermon, startingMiracle int
	switch n {
	case 2:
		startingStudy = 46
		startingPrayer = 40
		startingDevotion = 30
		startingParable = 8
		startingSermon = 8
		startingMiracle = 8
	case 3:
		startingStudy = 39
		startingPrayer = 40
		startingDevotion = 30
		startingParable = 12
		startingSermon = 12
		startingMiracle = 12
	case 4:
		startingStudy = 32
		startingPrayer = 40
		startingDevotion = 30
		startingParable = 12
		startingSermon = 12
		startingMiracle = 12
	case 5:
		startingStudy = 85
		startingPrayer = 80
		startingDevotion = 60
		startingParable = 12
		startingSermon = 12
		startingMiracle = 15
	default:
		startingStudy = 78
		startingPrayer = 80
		startingDevotion = 60
		startingParable = 12
		startingSermon = 12
		startingMiracle = 18
	}
	versePiles := []*VersePile{
		{Temptation, (n - 1) * 10},
		{Study, startingStudy},
		{Prayer, startingPrayer},
		{Devotion, startingDevotion},
		{Parable, startingParable},
		{Sermon, startingSermon},
		{Miracle, startingMiracle},
	}
	for _, c := range verses {
		versePiles = append(versePiles, &VersePile{c, 10}) // todo: glory cards
	}
	return &Kingdom{v: versePiles}
}
