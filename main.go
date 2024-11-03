package main

import (
	"bytes"
	"image/color"
	"log"
	"math/rand"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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

	BigFontSize    = 20
	NormalFontSize = 16
	SmallFontSize  = 12

	PlayedCardsY   = 250
	UnplayedCardsY = 310

	EndPhaseX      = 245
	EndPhaseY      = 370
	EndPhaseWidth  = 150
	EndPhaseHeight = 50
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

// stats for the current player, for display
type TurnStats struct {
	works, blessings, faith int
}

func (ts *TurnStats) reset() {
	ts.works = 1
	ts.blessings = 1
	ts.faith = 0
}

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
	// what state the game is in
	state string
	// for typing in the lobby
	t Typewriter
	// for sending messages between players
	pc *PulsarClient
	// all the players
	players map[string]PlayerData
	// which player gets which turn
	turnModulus []string
	// which turn are we on
	turn int
	// which phase is this turn in
	phase string
	// the active player's stats
	ts TurnStats
	// the kingdom piles
	kingdom *Kingdom
	// our cards
	myCards *PlayerCards
	// which cards are currently in play, to draw
	inPlayWork, inPlayFaith []*Card
}

// local player actions on turn end (end blessing phase)
func (g *Game) rest() {
	// put hand and in play cards into discard
	g.myCards.discard = slices.Concat(g.myCards.discard, g.inPlayWork, g.myCards.hand)
	// draw new hand
	g.myCards.DrawNCards(5)
	// tell everyone the blessing phase ended
	producerSend(g.pc.producer, []string{EndPhase})
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
		case EndPhase:
			switch g.phase {
			case WorkPhase:
			case BlessingPhase:
				// turn ended
				g.turn++
				g.phase = WorkPhase
				g.inPlayFaith = nil
				g.inPlayWork = nil
				g.ts.reset()
			}
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
				}
				g.turnModulus = names
			}
		}
	case Playing:
		// can only interact if it's our turn
		if g.turnModulus[g.turn%len(g.players)] == g.pc.playerName {
			switch g.phase {
			case WorkPhase:
			case BlessingPhase:
				// if no more blessings, auto rest
				if g.ts.blessings == 0 {
					g.rest()
					return nil
				}
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					cursorX, cursorY := ebiten.CursorPosition()
					// if clicked end phase, rest
					if cursorX > EndPhaseX && cursorX < EndPhaseX+EndPhaseWidth && cursorY > EndPhaseY && cursorY < EndPhaseY+EndPhaseHeight {
						g.rest()
						return nil
					}
					// if clicked card to buy, buy it
					if c := g.kingdom.In(cursorX, cursorY); c != nil {
						// only do something if you can afford it
						if g.ts.faith >= c.cost {
							// gains to discard
							g.myCards.discard = append(g.myCards.discard, c)
							// tell everyone you bought it so all kingdoms can decrement their supply
							producerSend(g.pc.producer, []string{Gained, g.pc.playerName, c.name})
							// TODO: check if game is done
						}
					}
				}
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
		// draw player's turn message
		msg := g.turnModulus[g.turn%len(g.players)] + "'s turn: " + g.phase
		op := &text.DrawOptions{}
		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, msg, &text.GoTextFace{
			Source: MPlusFaceSource,
			Size:   BigFontSize,
		}, op)
		// draw turn stats
		msg = "Works: " + strconv.Itoa(g.ts.works) + " Blessings: " + strconv.Itoa(g.ts.blessings) + " Faith: " + strconv.Itoa(g.ts.faith)
		op = &text.DrawOptions{}
		op.GeoM.Translate(0, BigFontSize)
		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, msg, &text.GoTextFace{
			Source: MPlusFaceSource,
			Size:   NormalFontSize,
		}, op)
		// draw player cards
		if g.myCards == nil {
			return
		}
		g.myCards.Draw(screen)
		// draw kingdom
		if g.kingdom == nil {
			return
		}
		g.kingdom.Draw(screen)
		// draw end phase button if it's our turn
		if g.turnModulus[g.turn%len(g.players)] == g.pc.playerName {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(EndPhaseX, EndPhaseY)
			switch g.phase {
			case WorkPhase:
				screen.DrawImage(assets.EndWorkPhase, op)
			case BlessingPhase:
				screen.DrawImage(assets.EndBlessingPhase, op)
			}
		}
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
		} else if c := g.kingdom.In(cursorX, cursorY); c != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(displayX), 0)
			screen.DrawImage(c.artBig, op)
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	g := &Game{state: RoomName, t: Typewriter{}, players: make(map[string]PlayerData), phase: BlessingPhase, ts: TurnStats{works: 1, blessings: 1, faith: 0}}

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}

	g.pc.Close()
}
