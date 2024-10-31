package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/zehongharryqu/kingdom-of-heaven/assets"
)

type Card struct {
	name               string
	artBig             *ebiten.Image
	artSmall           *ebiten.Image
	cost, glory, faith int
	cardTypes          []string
}

// card types
const (
	TemptationType = "Temptation"
	FaithType      = "Faith"
	GloryType      = "Glory"
	WorkType       = "Work"
	TrialType      = "Trial"
	ReactionType   = "Reaction"
)

var (
	Temptation   = &Card{name: "Temptation", artBig: assets.TemptationBig, artSmall: assets.TemptationSmall, glory: -1, cardTypes: []string{TemptationType}}
	Study        = &Card{name: "Study", artBig: assets.StudyBig, artSmall: assets.StudySmall, faith: 1, cardTypes: []string{FaithType}}
	Prayer       = &Card{name: "Prayer", artBig: assets.PrayerBig, artSmall: assets.PrayerSmall, cost: 3, faith: 2, cardTypes: []string{FaithType}}
	Devotion     = &Card{name: "Devotion", artBig: assets.DevotionBig, artSmall: assets.DevotionSmall, cost: 6, faith: 3, cardTypes: []string{FaithType}}
	Parable      = &Card{name: "Parable", artBig: assets.ParableBig, artSmall: assets.ParableSmall, cost: 2, glory: 1, cardTypes: []string{GloryType}}
	Sermon       = &Card{name: "Sermon", artBig: assets.SermonBig, artSmall: assets.SermonSmall, cost: 5, glory: 2, cardTypes: []string{GloryType}}
	Miracle      = &Card{name: "Miracle", artBig: assets.MiracleBig, artSmall: assets.MiracleSmall, cost: 8, glory: 3, cardTypes: []string{GloryType}}
	Bezalel      = &Card{name: "Bezalel", artBig: assets.BezalelBig, artSmall: assets.BezalelSmall, cost: 6, cardTypes: []string{WorkType}}
	Stumble      = &Card{name: "Stumble", artBig: assets.StumbleBig, artSmall: assets.StumbleSmall, cost: 5, cardTypes: []string{WorkType, TrialType}}
	Doubt        = &Card{name: "Doubt", artBig: assets.DoubtBig, artSmall: assets.DoubtSmall, cost: 4, cardTypes: []string{WorkType, TrialType}}
	NewCreation  = &Card{name: "NewCreation", artBig: assets.NewCreationBig, artSmall: assets.NewCreationSmall, cost: 2, cardTypes: []string{WorkType}}
	Purification = &Card{name: "Purification", artBig: assets.PurificationBig, artSmall: assets.PurificationSmall, cost: 2, cardTypes: []string{WorkType}}
	Feed5000     = &Card{name: "Feed5000", artBig: assets.Feed5000Big, artSmall: assets.Feed5000Small, cost: 5, cardTypes: []string{WorkType}}
	Festival     = &Card{name: "Festival", artBig: assets.FestivalBig, artSmall: assets.FestivalSmall, cost: 5, cardTypes: []string{WorkType}}
	Eden         = &Card{name: "Eden", artBig: assets.EdenBig, artSmall: assets.EdenSmall, cost: 4, cardTypes: []string{GloryType}}
	LostCoin     = &Card{name: "LostCoin", artBig: assets.LostCoinBig, artSmall: assets.LostCoinSmall, cost: 3, cardTypes: []string{WorkType}}
	Craft        = &Card{name: "Craft", artBig: assets.CraftBig, artSmall: assets.CraftSmall, cost: 5, cardTypes: []string{WorkType}}
	Collection   = &Card{name: "Collection", artBig: assets.CollectionBig, artSmall: assets.CollectionSmall, cost: 5, cardTypes: []string{WorkType}}
	Merchant     = &Card{name: "Merchant", artBig: assets.MerchantBig, artSmall: assets.MerchantSmall, cost: 5, cardTypes: []string{WorkType}}
	Belief       = &Card{name: "Belief", artBig: assets.BeliefBig, artSmall: assets.BeliefSmall, cost: 3, cardTypes: []string{WorkType}}
	Decree       = &Card{name: "Decree", artBig: assets.DecreeBig, artSmall: assets.DecreeSmall, cost: 4, cardTypes: []string{WorkType, TrialType}}
	GrowFaith    = &Card{name: "GrowFaith", artBig: assets.GrowFaithBig, artSmall: assets.GrowFaithSmall, cost: 5, cardTypes: []string{WorkType}}
	Shield       = &Card{name: "Shield", artBig: assets.ShieldBig, artSmall: assets.ShieldSmall, cost: 2, cardTypes: []string{WorkType, ReactionType}}
	Wisdom       = &Card{name: "Wisdom", artBig: assets.WisdomBig, artSmall: assets.WisdomSmall, cost: 4, cardTypes: []string{WorkType}}
	Depletion    = &Card{name: "Depletion", artBig: assets.DepletionBig, artSmall: assets.DepletionSmall, cost: 4, cardTypes: []string{WorkType}}
	Transform    = &Card{name: "Transform", artBig: assets.TransformBig, artSmall: assets.TransformSmall, cost: 4, cardTypes: []string{WorkType}}
	Plan         = &Card{name: "Plan", artBig: assets.PlanBig, artSmall: assets.PlanSmall, cost: 5, cardTypes: []string{WorkType}}
	Industry     = &Card{name: "Industry", artBig: assets.IndustryBig, artSmall: assets.IndustrySmall, cost: 4, cardTypes: []string{WorkType}}
	Duplication  = &Card{name: "Duplication", artBig: assets.DuplicationBig, artSmall: assets.DuplicationSmall, cost: 4, cardTypes: []string{WorkType}}
	Inspiration  = &Card{name: "Inspiration", artBig: assets.InspirationBig, artSmall: assets.InspirationSmall, cost: 3, cardTypes: []string{WorkType}}
	Bethlehem    = &Card{name: "Bethlehem", artBig: assets.BethlehemBig, artSmall: assets.BethlehemSmall, cost: 3, cardTypes: []string{WorkType}}
	Desires      = &Card{name: "Desires", artBig: assets.DesiresBig, artSmall: assets.DesiresSmall, cost: 5, cardTypes: []string{WorkType, TrialType}}
	Gift         = &Card{name: "Gift", artBig: assets.GiftBig, artSmall: assets.GiftSmall, cost: 3, cardTypes: []string{WorkType}}
)

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
