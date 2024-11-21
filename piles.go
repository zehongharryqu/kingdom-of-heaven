package main

import (
	"cmp"
	"image/color"
	"math/rand"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/zehongharryqu/kingdom-of-heaven/assets"
)

type PlayerCards struct {
	hand, deck, discard, decision []*Card
}

func InitPlayerCards() *PlayerCards {
	return &PlayerCards{discard: []*Card{Study, Study, Study, Study, Study, Study, Study, Parable, Parable, Parable}}
}

// draws n cards into dest and returns the result
func (pc *PlayerCards) drawNCards(n int, dest []*Card) []*Card {
	// fmt.Printf("hand %d, discard %d, deck %d, drawing %d cards\n", len(pc.hand), len(pc.discard), len(pc.deck), n)
	if len(pc.deck)+len(pc.discard) < n {
		// not enough in deck and discard, draw everything
		dest = slices.Concat(dest, pc.deck, pc.discard)
		pc.deck = nil
		pc.discard = nil
	} else {
		// enough cards in deck and discard, shuffle discard if necessary and draw from deck
		if len(pc.deck) < n {
			// not enough in just deck, shuffle discard and put it on the bottom of the deck
			rand.Shuffle(len(pc.discard), func(i, j int) {
				pc.discard[i], pc.discard[j] = pc.discard[j], pc.discard[i]
			})
			pc.deck = append(pc.discard, pc.deck...)
			pc.discard = nil
		}
		// draw into dest
		dest = append(dest, pc.deck[len(pc.deck)-n:]...)
		if len(pc.deck) == n {
			pc.deck = nil
		} else {
			pc.deck = pc.deck[:len(pc.deck)-n]
		}
	}
	// sort
	slices.SortFunc(dest, func(a, b *Card) int {
		return cmp.Or(
			cmp.Compare(slices.Min(a.cardTypes), slices.Min(b.cardTypes)),
			cmp.Compare(b.cost, a.cost),
			cmp.Compare(a.name, b.name),
		)
	})
	// fmt.Printf("hand %d, discard %d, deck %d\n", len(pc.hand), len(pc.discard), len(pc.deck))
	return dest
}

func (pc *PlayerCards) Draw(screen *ebiten.Image) {
	// decision
	offset := (ScreenWidth - ArtSmallWidth*len(pc.decision)) / 2
	for i, c := range pc.decision {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(offset+i*ArtSmallWidth), DecisionY)
		screen.DrawImage(c.artSmall, op)
	}
	// deck
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(DeckPileX, DiscardDeckPileY)
	screen.DrawImage(assets.DeckSmall, op)
	// discard
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(DiscardPileX, DiscardDeckPileY)
	if len(pc.discard) > 0 {
		screen.DrawImage(pc.discard[len(pc.discard)-1].artSmall, op)
	} else {
		screen.DrawImage(assets.DiscardSmall, op)
	}
	// hand
	offset = (ScreenWidth - ArtSmallWidth*len(pc.hand)) / 2
	for i, c := range pc.hand {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(offset+i*ArtSmallWidth), ScreenHeight-ArtSmallWidth)
		screen.DrawImage(c.artSmall, op)
	}
}

// given logical screen pixel location x,y returns the decision card there
func (pc *PlayerCards) inDecision(x, y int) (int, *Card) {
	localX := x - ((ScreenWidth - ArtSmallWidth*len(pc.decision)) / 2)
	localY := y - DecisionY
	if localX > 0 && localX < ArtSmallWidth*len(pc.decision) && localY > 0 && localY < ArtSmallWidth {
		return localX / ArtSmallWidth, pc.decision[localX/ArtSmallWidth]
	}
	return -1, nil
}

// given logical screen pixel location x,y returns the card in hand there
func (pc *PlayerCards) inHand(x, y int) (int, *Card) {
	localX := x - ((ScreenWidth - ArtSmallWidth*len(pc.hand)) / 2)
	localY := y - (ScreenHeight - ArtSmallWidth)
	if localX > 0 && localX < ArtSmallWidth*len(pc.hand) && localY > 0 && localY < ArtSmallWidth {
		return localX / ArtSmallWidth, pc.hand[localX/ArtSmallWidth]
	}
	return -1, nil
}

// given logical screen pixel location x,y returns the number of cards in deck if hovered
func (pc *PlayerCards) inDeck(x, y int) int {
	if x > DeckPileX && x < DeckPileX+ArtSmallWidth && y > DiscardDeckPileY && y < DiscardDeckPileY+ArtSmallWidth {
		return len(pc.deck)
	}
	return -1
}

// given logical screen pixel location x,y returns the detailed art and the number of cards in discard if hovered
func (pc *PlayerCards) inDiscard(x, y int) (*ebiten.Image, int) {
	if x > DiscardPileX && x < DiscardPileX+ArtSmallWidth && y > DiscardDeckPileY && y < DiscardDeckPileY+ArtSmallWidth {
		if n := len(pc.discard); n > 0 {
			return pc.discard[n-1].artBig, n
		} else {
			return nil, 0
		}
	}
	return nil, -1
}

// returns true if there are any works cards in hand
func (pc *PlayerCards) HasWorks() bool {
	for _, c := range pc.hand {
		if slices.Contains(c.cardTypes, WorkType) {
			return true
		}
	}
	return false
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
	KingdomMatX      = 330
	KingdomMatW      = 310
	KingdomMatH      = 240
)

// checks if the game is done
func (k *Kingdom) gameDone() bool {
	if k.v[6].n == 0 {
		return true
	}
	emptyPiles := 0
	for _, v := range k.v {
		if v.n == 0 {
			emptyPiles++
		}
	}
	return emptyPiles > 2
}

// draws the kingdom piles and released pile
func (k *Kingdom) Draw(screen *ebiten.Image) {
	// draw mat
	vector.DrawFilledRect(screen, KingdomMatX, 0, KingdomMatW, KingdomMatH+10+BigFontSize, color.RGBA{124, 54, 38, 255}, true)
	// draw mat label
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(KingdomMatX, KingdomMatH)
	textOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, "Verse Piles", &text.GoTextFace{
		Source: MPlusFaceSource,
		Size:   BigFontSize,
	}, textOp)
	// draw kingdom piles
	for i, v := range k.v {
		if v.n > 0 {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(KingdomPileX[i]), float64(KingdomPileY[i]))
			screen.DrawImage(v.c.artSmall, op)
		}
	}
	// draw released pile
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(ReleasedPileX, ReleasedPileY)
	if n := len(k.released); n > 0 {
		screen.DrawImage(k.released[n-1].artSmall, op)
	} else {
		screen.DrawImage(assets.ReleaseSmall, op)
	}
}

// given logical screen pixel location x,y returns the card if there is a kingdom card there
func (k *Kingdom) In(x, y int) *VersePile {
	for i, v := range k.v {
		if x > KingdomPileX[i] && x < KingdomPileX[i]+ArtSmallWidth && y > KingdomPileY[i] && y < KingdomPileY[i]+ArtSmallWidth {
			return v
		}
	}
	return nil
}

// removes a card from the kingdom (e.g. when gained)
func (k *Kingdom) RemoveCard(name string) {
	for _, v := range k.v {
		if v.c.name == name {
			v.n--
			return
		}
	}
}

// create a new kingdom given the 10 verses and number of players
func InitKingdom(verses []*Card, n int) *Kingdom {
	// starting amounts from the dominion wiki gameplay article
	var startingStudy, startingPrayer, startingDevotion, startingGlory, startingMiracle int
	switch n {
	case 2:
		startingStudy = 46
		startingPrayer = 40
		startingDevotion = 30
		startingGlory = 8
		startingMiracle = 8
	case 3:
		startingStudy = 39
		startingPrayer = 40
		startingDevotion = 30
		startingGlory = 12
		startingMiracle = 12
	case 4:
		startingStudy = 32
		startingPrayer = 40
		startingDevotion = 30
		startingGlory = 12
		startingMiracle = 12
	case 5:
		startingStudy = 85
		startingPrayer = 80
		startingDevotion = 60
		startingGlory = 12
		startingMiracle = 15
	default:
		startingStudy = 78
		startingPrayer = 80
		startingDevotion = 60
		startingGlory = 12
		startingMiracle = 18
	}
	versePiles := []*VersePile{
		{Temptation, (n - 1) * 10},
		{Study, startingStudy},
		{Prayer, startingPrayer},
		{Devotion, startingDevotion},
		{Parable, startingGlory},
		{Sermon, startingGlory},
		{Miracle, startingMiracle},
	}
	// sort kingdom by cost and name
	slices.SortFunc(verses, func(a, b *Card) int {
		return cmp.Or(
			cmp.Compare(a.cost, b.cost),
			cmp.Compare(a.name, b.name),
		)
	})
	for _, c := range verses {
		startingAmount := 10
		if slices.Contains(c.cardTypes, GloryType) {
			startingAmount = startingGlory
		}
		versePiles = append(versePiles, &VersePile{c, startingAmount})
	}
	return &Kingdom{v: versePiles}
}
