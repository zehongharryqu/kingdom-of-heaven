package main

import (
	"bytes"
	"image/color"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// sizes
const (
	ScreenHeight = 480
	ScreenWidth  = 640

	ArtBigHeight  = 400
	ArtBigWidth   = 300
	ArtSmallWidth = 50

	MaxNameChars = 10

	NormalFontSize = 20

	PlayedCardsY   = 250
	UnplayedCardsY = 310
)

// game states
const (
	RoomName = "RoomName"
	Lobby    = "Lobby"
	Playing  = "Playing"
)

// turn phases
const (
	WorkPhase     = "Work Phase"
	BlessingPhase = "Blessing Phase"
)

var (
	MPlusFaceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	MPlusFaceSource = s
}

type Game struct {
	state       string
	t           Typewriter
	pc          *PulsarClient
	players     map[string]PlayerData
	turnModulus []string
	turn        int
	phase       string
	kingdom     *Kingdom
	myCards     *PlayerCards
}

func (g *Game) ReceiveMessages() {
	for {
		message := consumerReceive(g.pc.consumer)

		switch message[0] {
		case JoinedLobby:
			pid, err := strconv.Atoi(message[2])
			if err != nil {
				log.Fatal(err)
			}
			g.players[message[1]] = PlayerData{pid: pid, ready: false}
		case LeftLobby:
			delete(g.players, message[1])
		case ToggledReady:
			g.players[message[1]] = PlayerData{pid: g.players[message[1]].pid, ready: !g.players[message[1]].ready}
		case SetKingdom:
			// generate local kingdom from message
			cards := make([]*Card, 10)
			for i, c := range message[1:] {
				idx, _ := strconv.Atoi(c)
				cards[i] = NonBaseCards[idx]
			}
			g.kingdom = InitKingdom(cards, len(g.players))
			// create deck and discard
			g.myCards = InitPlayerCards()
			g.myCards.DrawNCards(5)
		case Clicked:
		}
	}
}

func (g *Game) Update() error {
	switch g.state {
	case RoomName:
		g.t.Update()
		if g.t.confirmedName != "" && g.t.confirmedRoom != "" {
			g.pc = newPulsarClient(g.t.confirmedRoom, g.t.confirmedName)
			g.state = Lobby
			go g.ReceiveMessages()
		}
	case Lobby:
		if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
			producerSend(g.pc.producer, []string{ToggledReady, g.pc.playerName})
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
				// create turn order by pid
				names := make([]string, len(g.players))
				i := 0
				for name := range g.players {
					names[i] = name
					i++
				}
				sort.SliceStable(names, func(i, j int) bool {
					return g.players[names[i]].pid < g.players[names[j]].pid
				})
				// the first player generates the kingdom
				if names[0] == g.pc.playerName {
					// create strings out of non base card indices, get a random 10, and send it to everyone
					nonBaseCardStrings := make([]string, len(NonBaseCards))
					for i := range NonBaseCards {
						nonBaseCardStrings[i] = strconv.Itoa(i)
					}
					rand.Shuffle(len(nonBaseCardStrings), func(i, j int) {
						nonBaseCardStrings[i], nonBaseCardStrings[j] = nonBaseCardStrings[j], nonBaseCardStrings[i]
					})
					producerSend(g.pc.producer, append([]string{SetKingdom}, nonBaseCardStrings[:10]...))
					// g.kingdom = InitKingdom(NonBaseCards[], len(names))
				}
				g.turnModulus = names
			}
		}
	case Playing:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			producerSend(g.pc.producer, []string{Clicked, g.pc.playerName})
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
		// draw player's turn message
		msg := g.turnModulus[g.turn%len(g.players)] + "'s turn: " + g.phase
		op := &text.DrawOptions{}
		// op.GeoM.Translate(10, 10)
		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, msg, &text.GoTextFace{
			Source: MPlusFaceSource,
			Size:   NormalFontSize,
		}, op)
		// draw player cards
		g.myCards.Draw(screen)
		// calculate mouse position to determine hover
		cursorX, cursorY := ebiten.CursorPosition()
		var displayX int
		if cursorX > ScreenWidth/2 {
			displayX = cursorX - ArtBigWidth
		} else {
			displayX = cursorX
		}
		if displayArt := g.myCards.In(cursorX, cursorY); displayArt != nil {
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
	g := &Game{state: RoomName, t: Typewriter{}, players: make(map[string]PlayerData)}

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}

	g.pc.Close()
}
