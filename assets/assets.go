package assets

import (
	"embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed *
var assets embed.FS

func loadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

var (
	TemptationBig     = loadImage("big1.png")
	StudyBig          = loadImage("big2.png")
	PrayerBig         = loadImage("big3.png")
	DevotionBig       = loadImage("big4.png")
	ParableBig        = loadImage("big5.png")
	SermonBig         = loadImage("big6.png")
	MiracleBig        = loadImage("big7.png")
	BezalelBig        = loadImage("big8.png")
	StumbleBig        = loadImage("big9.png")
	DoubtBig          = loadImage("big10.png")
	NewCreationBig    = loadImage("big11.png")
	PurificationBig   = loadImage("big12.png")
	Feed5000Big       = loadImage("big13.png")
	FestivalBig       = loadImage("big14.png")
	EdenBig           = loadImage("big15.png")
	LostCoinBig       = loadImage("big16.png")
	CraftBig          = loadImage("big17.png")
	CollectionBig     = loadImage("big18.png")
	MerchantBig       = loadImage("big19.png")
	BeliefBig         = loadImage("big20.png")
	DecreeBig         = loadImage("big21.png")
	GrowFaithBig      = loadImage("big22.png")
	ShieldBig         = loadImage("big23.png")
	WisdomBig         = loadImage("big24.png")
	DepletionBig      = loadImage("big25.png")
	TransformBig      = loadImage("big26.png")
	PlanBig           = loadImage("big27.png")
	IndustryBig       = loadImage("big28.png")
	DuplicationBig    = loadImage("big29.png")
	InspirationBig    = loadImage("big30.png")
	BethlehemBig      = loadImage("big31.png")
	DesiresBig        = loadImage("big32.png")
	GiftBig           = loadImage("big33.png")
	TemptationSmall   = loadImage("small1.png")
	StudySmall        = loadImage("small2.png")
	PrayerSmall       = loadImage("small3.png")
	DevotionSmall     = loadImage("small4.png")
	ParableSmall      = loadImage("small5.png")
	SermonSmall       = loadImage("small6.png")
	MiracleSmall      = loadImage("small7.png")
	BezalelSmall      = loadImage("small8.png")
	StumbleSmall      = loadImage("small9.png")
	DoubtSmall        = loadImage("small10.png")
	NewCreationSmall  = loadImage("small11.png")
	PurificationSmall = loadImage("small12.png")
	Feed5000Small     = loadImage("small13.png")
	FestivalSmall     = loadImage("small14.png")
	EdenSmall         = loadImage("small15.png")
	LostCoinSmall     = loadImage("small16.png")
	CraftSmall        = loadImage("small17.png")
	CollectionSmall   = loadImage("small18.png")
	MerchantSmall     = loadImage("small19.png")
	BeliefSmall       = loadImage("small20.png")
	DecreeSmall       = loadImage("small21.png")
	GrowFaithSmall    = loadImage("small22.png")
	ShieldSmall       = loadImage("small23.png")
	WisdomSmall       = loadImage("small24.png")
	DepletionSmall    = loadImage("small25.png")
	TransformSmall    = loadImage("small26.png")
	PlanSmall         = loadImage("small27.png")
	IndustrySmall     = loadImage("small28.png")
	DuplicationSmall  = loadImage("small29.png")
	InspirationSmall  = loadImage("small30.png")
	BethlehemSmall    = loadImage("small31.png")
	DesiresSmall      = loadImage("small32.png")
	GiftSmall         = loadImage("small33.png")

	DiscardSmall = loadImage("discard.png")
	DeckSmall    = loadImage("deck.png")
	ReleaseSmall = loadImage("release.png")

	EndWorkPhase     = loadImage("endworkphase.png")
	EndBlessingPhase = loadImage("endworkphase.png")
)
