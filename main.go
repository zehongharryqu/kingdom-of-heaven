package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/zehongharryqu/kingdom-of-heaven/assets"
)

// sizes
const (
	ScreenHeight = 480
	ScreenWidth  = 640

	ArtBigHeight  = 400
	ArtBigWidth   = 300
	ArtSmallWidth = 50

	MaxNameChars = 10
)

// game states
const (
	RoomName = "RoomName"
	Lobby    = "Lobby"
	Playing  = "Playing"
)

type Game struct {
	state   string
	t       Typewriter
	pc      *PulsarClient
	players map[string]PlayerData
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
						g.players[messageText[len(PlayerJoined):]] = PlayerData{false}
					} else if strings.HasPrefix(messageText, PlayerLeft) {
						delete(g.players, messageText[len(PlayerLeft):])
					} else if strings.HasPrefix(messageText, PlayerToggledReady) {
						g.players[messageText[len(PlayerToggledReady):]] = PlayerData{!g.players[messageText[len(PlayerToggledReady):]].ready}
					}
					g.pc.consumer.Ack(msg)
				}
			}()
		}
	case Lobby:
		if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
			producerSend(g.pc.producer, PlayerToggledReady+g.pc.playerName)
		}
		if len(g.players) > 1 {
			ready := true
			for _, playerData := range g.players {
				if !playerData.ready {
					ready = false
				}
			}
			if ready {
				g.state = Playing
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case RoomName:
		g.t.Draw(screen)
	case Lobby:
		lobbyMessage := "Room " + g.t.confirmedRoom + "\nHit enter when ready to start\n\nPlayers in this room:\n"
		// sort player names otherwise it keeps switching them around
		names := make([]string, len(g.players))

		i := 0
		for name := range g.players {
			names[i] = name
			i++
		}
		sort.Strings(names)
		for _, name := range names {
			lobbyMessage += name + strings.Repeat(" ", MaxNameChars+1-len(name)) + "| "
			if g.players[name].ready {
				lobbyMessage += "Ready!\n"
			} else {
				lobbyMessage += "Waiting...\n"
			}
		}
		ebitenutil.DebugPrint(screen, lobbyMessage)
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
	c := &Card{assets.CardBig, assets.CardSmall}
	h := &Hand{[]*Card{c, c, c, c, c}}
	g := &Game{RoomName, Typewriter{}, nil, make(map[string]PlayerData), h}

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}

	g.pc.Close()
}
