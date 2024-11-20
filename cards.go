package main

import (
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/zehongharryqu/kingdom-of-heaven/assets"
)

type Card struct {
	name               string
	artBig             *ebiten.Image
	artSmall           *ebiten.Image
	cost, glory, faith int
	cardTypes          []int
}

// card types: for sorting hand
const (
	TemptationType = 5
	FaithType      = 3
	GloryType      = 4
	WorkType       = 1
	TrialType      = 0
	ReactionType   = 2
)

// cards
var (
	Temptation   = &Card{name: "Temptation", artBig: assets.TemptationBig, artSmall: assets.TemptationSmall, glory: -1, cardTypes: []int{TemptationType}}
	Study        = &Card{name: "Study", artBig: assets.StudyBig, artSmall: assets.StudySmall, faith: 1, cardTypes: []int{FaithType}}
	Prayer       = &Card{name: "Prayer", artBig: assets.PrayerBig, artSmall: assets.PrayerSmall, cost: 3, faith: 2, cardTypes: []int{FaithType}}
	Devotion     = &Card{name: "Devotion", artBig: assets.DevotionBig, artSmall: assets.DevotionSmall, cost: 6, faith: 3, cardTypes: []int{FaithType}}
	Parable      = &Card{name: "Parable", artBig: assets.ParableBig, artSmall: assets.ParableSmall, cost: 2, glory: 1, cardTypes: []int{GloryType}}
	Sermon       = &Card{name: "Sermon", artBig: assets.SermonBig, artSmall: assets.SermonSmall, cost: 5, glory: 3, cardTypes: []int{GloryType}}
	Miracle      = &Card{name: "Miracle", artBig: assets.MiracleBig, artSmall: assets.MiracleSmall, cost: 8, glory: 6, cardTypes: []int{GloryType}}
	Bezalel      = &Card{name: "Bezalel", artBig: assets.BezalelBig, artSmall: assets.BezalelSmall, cost: 6, cardTypes: []int{WorkType}}
	Stumble      = &Card{name: "Stumble", artBig: assets.StumbleBig, artSmall: assets.StumbleSmall, cost: 5, cardTypes: []int{WorkType, TrialType}}
	Doubt        = &Card{name: "Doubt", artBig: assets.DoubtBig, artSmall: assets.DoubtSmall, cost: 4, cardTypes: []int{WorkType, TrialType}}
	NewCreation  = &Card{name: "NewCreation", artBig: assets.NewCreationBig, artSmall: assets.NewCreationSmall, cost: 2, cardTypes: []int{WorkType}}
	Purification = &Card{name: "Purification", artBig: assets.PurificationBig, artSmall: assets.PurificationSmall, cost: 2, cardTypes: []int{WorkType}}
	Feed5000     = &Card{name: "Feed5000", artBig: assets.Feed5000Big, artSmall: assets.Feed5000Small, cost: 5, cardTypes: []int{WorkType}}
	Festival     = &Card{name: "Festival", artBig: assets.FestivalBig, artSmall: assets.FestivalSmall, cost: 5, cardTypes: []int{WorkType}}
	Eden         = &Card{name: "Eden", artBig: assets.EdenBig, artSmall: assets.EdenSmall, cost: 4, cardTypes: []int{GloryType}}
	LostCoin     = &Card{name: "LostCoin", artBig: assets.LostCoinBig, artSmall: assets.LostCoinSmall, cost: 3, cardTypes: []int{WorkType}}
	Craft        = &Card{name: "Craft", artBig: assets.CraftBig, artSmall: assets.CraftSmall, cost: 5, cardTypes: []int{WorkType}}
	Collection   = &Card{name: "Collection", artBig: assets.CollectionBig, artSmall: assets.CollectionSmall, cost: 5, cardTypes: []int{WorkType}}
	Merchant     = &Card{name: "Merchant", artBig: assets.MerchantBig, artSmall: assets.MerchantSmall, cost: 5, cardTypes: []int{WorkType}}
	Belief       = &Card{name: "Belief", artBig: assets.BeliefBig, artSmall: assets.BeliefSmall, cost: 3, cardTypes: []int{WorkType}}
	Decree       = &Card{name: "Decree", artBig: assets.DecreeBig, artSmall: assets.DecreeSmall, cost: 4, cardTypes: []int{WorkType, TrialType}}
	GrowFaith    = &Card{name: "GrowFaith", artBig: assets.GrowFaithBig, artSmall: assets.GrowFaithSmall, cost: 5, cardTypes: []int{WorkType}}
	Shield       = &Card{name: "Shield", artBig: assets.ShieldBig, artSmall: assets.ShieldSmall, cost: 2, cardTypes: []int{WorkType, ReactionType}}
	Wisdom       = &Card{name: "Wisdom", artBig: assets.WisdomBig, artSmall: assets.WisdomSmall, cost: 4, cardTypes: []int{WorkType}}
	Depletion    = &Card{name: "Depletion", artBig: assets.DepletionBig, artSmall: assets.DepletionSmall, cost: 4, cardTypes: []int{WorkType}}
	Transform    = &Card{name: "Transform", artBig: assets.TransformBig, artSmall: assets.TransformSmall, cost: 4, cardTypes: []int{WorkType}}
	Plan         = &Card{name: "Plan", artBig: assets.PlanBig, artSmall: assets.PlanSmall, cost: 5, cardTypes: []int{WorkType}}
	Industry     = &Card{name: "Industry", artBig: assets.IndustryBig, artSmall: assets.IndustrySmall, cost: 4, cardTypes: []int{WorkType}}
	Duplication  = &Card{name: "Duplication", artBig: assets.DuplicationBig, artSmall: assets.DuplicationSmall, cost: 4, cardTypes: []int{WorkType}}
	Inspiration  = &Card{name: "Inspiration", artBig: assets.InspirationBig, artSmall: assets.InspirationSmall, cost: 3, cardTypes: []int{WorkType}}
	Bethlehem    = &Card{name: "Bethlehem", artBig: assets.BethlehemBig, artSmall: assets.BethlehemSmall, cost: 3, cardTypes: []int{WorkType}}
	Desires      = &Card{name: "Desires", artBig: assets.DesiresBig, artSmall: assets.DesiresSmall, cost: 5, cardTypes: []int{WorkType, TrialType}}
	Gift         = &Card{name: "Gift", artBig: assets.GiftBig, artSmall: assets.GiftSmall, cost: 3, cardTypes: []int{WorkType}}
)

// cards for randomization
var NonBaseCards = []*Card{
	Bezalel,
	Stumble,
	Doubt,
	NewCreation,
	Purification,
	Feed5000,
	Festival,
	Eden,
	LostCoin,
	Craft,
	Collection,
	Merchant,
	Belief,
	Decree,
	GrowFaith,
	Shield,
	Wisdom,
	Depletion,
	Transform,
	Plan,
	Industry,
	Duplication,
	Inspiration,
	Bethlehem,
	Desires,
	Gift}

// convert string name into card
var CardNameMap = map[string]*Card{
	"Temptation":   Temptation,
	"Study":        Study,
	"Prayer":       Prayer,
	"Devotion":     Devotion,
	"Parable":      Parable,
	"Sermon":       Sermon,
	"Miracle":      Miracle,
	"Bezalel":      Bezalel,
	"Stumble":      Stumble,
	"Doubt":        Doubt,
	"NewCreation":  NewCreation,
	"Purification": Purification,
	"Feed5000":     Feed5000,
	"Festival":     Festival,
	"Eden":         Eden,
	"LostCoin":     LostCoin,
	"Craft":        Craft,
	"Collection":   Collection,
	"Merchant":     Merchant,
	"Belief":       Belief,
	"Decree":       Decree,
	"GrowFaith":    GrowFaith,
	"Shield":       Shield,
	"Wisdom":       Wisdom,
	"Depletion":    Depletion,
	"Transform":    Transform,
	"Plan":         Plan,
	"Industry":     Industry,
	"Duplication":  Duplication,
	"Inspiration":  Inspiration,
	"Bethlehem":    Bethlehem,
	"Desires":      Desires,
	"Gift":         Gift,
}

// which cards require decisions
const (
	DecisionBezalel1 = iota
	DecisionBezalel2
	DecisionStumble
)

// what runs when you play a card
func (g *Game) localCardEffect(c *Card) {
	producerSend(g.pc.producer, []string{Played, g.pc.playerName, c.name})
	switch c.name {
	case Bezalel.name:
		g.decision = DecisionBezalel1
	case Stumble.name:
		// gain Devotion if there are any
		if g.kingdom.v[3].n > 0 {
			g.myCards.discard = append(g.myCards.discard, Devotion)
			// tell everyone you gained it so all kingdoms can decrement their supply
			producerSend(g.pc.producer, []string{Gained, g.pc.playerName, Devotion.name})
		}
		g.otherDecisions = len(g.players) - 1
	case Doubt.name:
	case NewCreation.name:
	case Purification.name:
	case Feed5000.name:
	case Festival.name:
	case Eden.name:
	case LostCoin.name:
	case Craft.name:
	case Collection.name:
	case Merchant.name:
	case Belief.name:
	case Decree.name:
	case GrowFaith.name:
	case Shield.name:
	case Wisdom.name:
	case Depletion.name:
	case Transform.name:
	case Plan.name:
	case Industry.name:
	case Duplication.name:
	case Inspiration.name:
	case Bethlehem.name:
	case Desires.name:
	case Gift.name:
	}
}

// what runs when others play a card
func (g *Game) reactToCard(c *Card) {
	switch c.name {
	case Stumble.name:
		// reveal top 2 cards
		g.myCards.decision = g.myCards.drawNCards(2, g.myCards.decision)
		switch {
		// both cards are non-Study Faiths, make decision
		case (slices.Contains(g.myCards.decision[0].cardTypes, FaithType) && g.myCards.decision[0].name != "Study") && (slices.Contains(g.myCards.decision[1].cardTypes, FaithType) && g.myCards.decision[1].name != "Study"):
			g.decision = DecisionStumble
		// first card is a non-Study faith, release it
		case slices.Contains(g.myCards.decision[0].cardTypes, FaithType) && g.myCards.decision[0].name != "Study":
			// tell everyone to add to release pile
			producerSend(g.pc.producer, []string{CardSpecific, g.pc.playerName, Stumble.name, g.myCards.decision[0].name, g.myCards.decision[1].name, g.myCards.decision[0].name})
			// discard second card
			g.myCards.discard = append(g.myCards.discard, g.myCards.decision[1])
			g.myCards.decision = nil
		// second card is a non-Study faith, release it
		case slices.Contains(g.myCards.decision[1].cardTypes, FaithType) && g.myCards.decision[1].name != "Study":
			// tell everyone to add to release pile
			producerSend(g.pc.producer, []string{CardSpecific, g.pc.playerName, Stumble.name, g.myCards.decision[0].name, g.myCards.decision[1].name, g.myCards.decision[1].name})
			// discard first card
			g.myCards.discard = append(g.myCards.discard, g.myCards.decision[0])
			g.myCards.decision = nil
		// neither card is a non-Study faith, discard both
		default:
			producerSend(g.pc.producer, []string{CardSpecific, g.pc.playerName, Stumble.name, g.myCards.decision[0].name, g.myCards.decision[1].name})
			g.myCards.discard = append(g.myCards.discard, g.myCards.decision...)
			g.myCards.decision = nil
		}
	case Doubt.name:
	case NewCreation.name:
	case Purification.name:
	case Feed5000.name:
	case Festival.name:
	case Eden.name:
	case LostCoin.name:
	case Craft.name:
	case Collection.name:
	case Merchant.name:
	case Belief.name:
	case Decree.name:
	case GrowFaith.name:
	case Shield.name:
	case Wisdom.name:
	case Depletion.name:
	case Transform.name:
	case Plan.name:
	case Industry.name:
	case Duplication.name:
	case Inspiration.name:
	case Bethlehem.name:
	case Desires.name:
	case Gift.name:
	}
}

// react to local decisions made
func (g *Game) listenForDecision() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cursorX, cursorY := ebiten.CursorPosition()
		switch g.decision {
		case DecisionBezalel1:
			// if clicked card to gain, gain it
			if vp := g.kingdom.In(cursorX, cursorY); vp != nil {
				// only do something if cost is at most 5 and there are cards left
				if vp.c.cost <= 5 && vp.n > 0 {
					// gains to hand
					g.myCards.hand = append(g.myCards.hand, vp.c)
					// tell everyone you gained it so all kingdoms can decrement their supply
					producerSend(g.pc.producer, []string{Gained, g.pc.playerName, vp.c.name})
					// move to next part (put card on deck)
					g.decision = DecisionBezalel2
				}
			}
		case DecisionBezalel2:
			if i, c := g.myCards.inHand(cursorX, cursorY); c != nil {
				// no more decision
				g.decision = -1
				// put in front of deck
				g.myCards.deck = append([]*Card{c}, g.myCards.deck...)
				// remove from hand
				g.myCards.hand = append(g.myCards.hand[:i], g.myCards.hand[i+1:]...)
			}
		case DecisionStumble:
			if i, c := g.myCards.inDecision(cursorX, cursorY); c != nil {
				// no more decision
				g.decision = -1
				// tell everyone to add to release pile
				producerSend(g.pc.producer, []string{CardSpecific, g.pc.playerName, Stumble.name, g.myCards.decision[0].name, g.myCards.decision[1].name, c.name})
				// discard other card
				g.myCards.discard = append(g.myCards.discard, g.myCards.decision[1-i])
				g.myCards.decision = nil
			}
		}
	}
}

// what message should be shown to the player, and can they skip it?
func (g *Game) promptDecision() (string, bool) {
	switch g.decision {
	case -1:
		return "", false
	case DecisionBezalel1:
		return "Select a card to gain costing up to 5 Faith", false
	case DecisionBezalel2:
		return "Select a card from your hand to put on your deck", false
	case DecisionStumble:
		return "Select a card to release", false
	}
	return "", false
}
