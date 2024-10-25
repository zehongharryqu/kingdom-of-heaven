package main

import (
	"image"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin"
	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/admin"
	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/zehongharryqu/kingdom-of-heaven/assets"
)

const (
	ScreenHeight = 480
	ScreenWidth  = 640
)
const (
	ArtBigHeight  = 400
	ArtBigWidth   = 300
	ArtSmallWidth = 50
)

// game states
const (
	RoomName = "RoomName"
	Lobby    = "Lobby"
	Playing  = "Playing"
)

type Card struct {
	artBig     *ebiten.Image
	artSmall   *ebiten.Image
	alphaImage *image.Alpha
}

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

type Game struct {
	state   string
	t       Typewriter
	pa      admin.Client
	pc      *PulsarClient
	players []string
	h       *Hand
}

func (g *Game) Update() error {
	switch g.state {
	case RoomName:
		g.t.Update()
		if g.t.confirmedName != "" && g.t.confirmedRoom != "" {
			g.pc = newPulsarClient(g.t.confirmedRoom, g.t.confirmedName)
			g.state = Lobby // Todo: implement Lobby
		}
	case Lobby:
		topicName, err := utils.GetTopicName(g.pc.roomName)
		if err != nil {
			panic(err)
		}
		topicStats, err := g.pa.Topics().GetStats(*topicName)
		if err != nil {
			panic(err)
		}
		subscriptions := topicStats.Subscriptions
		keys := make([]string, len(subscriptions))
		i := 0
		for k := range subscriptions {
			keys[i] = k
			i++
		}
		g.players = keys
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case RoomName:
		g.t.Draw(screen)
	case Lobby:
		ebitenutil.DebugPrint(screen, "Players in this room:\n"+strings.Join(g.players, "\n"))
	case Playing:
		g.h.Draw(screen)
		cursorX, cursorY := ebiten.CursorPosition()
		var displayX int
		if cursorX > ScreenWidth/2 {
			displayX = cursorX - ArtBigWidth
		} else {
			displayX = cursorX
		}
		if displayArt := g.h.In(cursorX, cursorY); displayArt != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(displayX), 0)
			screen.DrawImage(displayArt, op)
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {

	cfg := &pulsaradmin.Config{
		AuthPlugin:     "oauth2",
		IssuerEndpoint: "https://auth.streamnative.cloud/",
		Audience:       "urn:sn:pulsar:o-hwa6o:kingdom-of-heaven-instance",
		KeyFile:        "file:///Users/harry/Downloads/o-hwa6o-harry.json",
		WebServiceURL:  "https://pc-de347430.gcp-shared-usce1.g.snio.cloud:6651"}
	admin, err := pulsaradmin.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	c := &Card{assets.CardBig, assets.CardSmall, nil}
	h := &Hand{[]*Card{c, c, c, c, c}}
	g := &Game{RoomName, Typewriter{}, admin, nil, nil, h}

	err = ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}

	g.pc.Close()
}
