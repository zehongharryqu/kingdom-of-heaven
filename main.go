package main

import (
	"context"
	"fmt"
	"image"
	"log"
	"slices"
	"strings"

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
			g.state = Lobby
			go func() {
				for g.state == Lobby {
					msg, err := g.pc.consumer.Receive(context.Background())
					if err != nil {
						log.Fatal(err)
					}
					messageText := string(msg.Payload())
					fmt.Printf("Received message msgId: %v -- content: '%s'\n",
						msg.ID(), messageText)

					if strings.HasPrefix(messageText, PlayerJoined) {
						g.players = append(g.players, messageText[len(PlayerJoined):])
					} else if strings.HasPrefix(messageText, PlayerLeft) {
						g.players = slices.DeleteFunc(g.players, func(str string) bool { return str == messageText[len(PlayerLeft):] })
					}
					g.pc.consumer.Ack(msg)
				}
			}()
		}
	case Lobby:

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
	c := &Card{assets.CardBig, assets.CardSmall, nil}
	h := &Hand{[]*Card{c, c, c, c, c}}
	g := &Game{RoomName, Typewriter{}, nil, nil, h}

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}

	g.pc.Close()
}
