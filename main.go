package main

import (
	"bytes"
	"fmt"
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
	"github.com/hajimehoshi/ebiten/v2/vector"
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

	InPlayWorkY  = 250
	InPlayFaithY = 310

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
	// list of actions that have occured, last 10 of which are drawn
	actionLog []string
}

// show everyone your faith cards at the start of your blessing phase (end work phase)
func (g *Game) startBlessing() {
	var faithCards []string
	for _, c := range g.myCards.hand {
		if slices.Contains(c.cardTypes, FaithType) {
			faithCards = append(faithCards, c.name)
		}
	}
	producerSend(g.pc.producer, append([]string{EndPhase}, faithCards...))
	// spin until the phase changes
	for g.phase == WorkPhase {
	}
}

// local player actions on turn end (end blessing phase)
func (g *Game) rest() {
	// put hand and in play cards into discard
	g.myCards.discard = slices.Concat(g.myCards.discard, g.inPlayWork, g.myCards.hand)
	fmt.Printf("discard len: %d \n", len(g.myCards.discard))
	// draw new hand
	g.myCards.hand = nil
	g.myCards.DrawNCards(5)
	// tell everyone the blessing phase ended
	producerSend(g.pc.producer, []string{EndPhase})
	// spin until the phase changes
	for g.phase == BlessingPhase {
	}
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
			if name := message[1]; name == g.pc.playerName {
				// if we are leaving, close our producer and consumer
				g.pc.Close()
			} else {
				// if someone else is leaving, remove them
				delete(g.players, name)
			}
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
				// moving to blessing phase
				g.phase = BlessingPhase
				// see which faith cards they are playing and calculate their faith
				for _, c := range message[1:] {
					g.inPlayFaith = append(g.inPlayFaith, CardNameMap[c])
					g.ts.faith += CardNameMap[c].faith
				}
			case BlessingPhase:
				// turn ended
				g.turn++
				g.phase = WorkPhase
				g.inPlayFaith = nil
				g.inPlayWork = nil
				g.ts.reset()
			}
		case Gained:
			// write that the player gained the card
			g.actionLog = append(g.actionLog, message[1]+" gained "+message[2])
			// remove a card from supply
			g.kingdom.RemoveCard(message[2])
			// decrement the player's blessings and faith
			g.ts.blessings--
			g.ts.faith -= CardNameMap[message[2]].cost
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
		if g.myCards == nil {
			return nil
		}
		// can only interact if it's our turn
		if g.turnModulus[g.turn%len(g.players)] == g.pc.playerName {
			switch g.phase {
			case WorkPhase:
				// if no more works, auto start blessing
				if g.ts.works == 0 || !g.myCards.HasWorks() {
					fmt.Println(g.pc.playerName + " has no works, starting blessing")
					g.startBlessing()
					return nil
				}
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					cursorX, cursorY := ebiten.CursorPosition()
					// if clicked end phase, start blessing
					if cursorX > EndPhaseX && cursorX < EndPhaseX+EndPhaseWidth && cursorY > EndPhaseY && cursorY < EndPhaseY+EndPhaseHeight {
						fmt.Println(g.pc.playerName + " clicked end works phase, starting blessing")
						g.startBlessing()
						return nil
					}
					// TODO: click to play work cards
				}
			case BlessingPhase:
				// if no more blessings, auto rest
				if g.ts.blessings == 0 {
					fmt.Println(g.pc.playerName + " has no blessings, ending turn")
					g.rest()
					return nil
				}
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					cursorX, cursorY := ebiten.CursorPosition()
					// if clicked end phase, rest
					if cursorX > EndPhaseX && cursorX < EndPhaseX+EndPhaseWidth && cursorY > EndPhaseY && cursorY < EndPhaseY+EndPhaseHeight {
						fmt.Println(g.pc.playerName + " clicked end blessings phase, ending turn")
						g.rest()
						return nil
					}
					// if clicked card to buy, buy it
					if vp := g.kingdom.In(cursorX, cursorY); vp != nil {
						// only do something if you can afford it
						if g.ts.faith >= vp.c.cost {
							// gains to discard
							g.myCards.discard = append(g.myCards.discard, vp.c)
							// tell everyone you bought it so all kingdoms can decrement their supply
							producerSend(g.pc.producer, []string{Gained, g.pc.playerName, vp.c.name})
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
		// draw action log
		if n := len(g.actionLog); n > 10 {
			msg = strings.Join(g.actionLog[n-10:], "\n")
		} else {
			msg = strings.Join(g.actionLog, "\n")
		}
		op = &text.DrawOptions{}
		op.GeoM.Translate(0, BigFontSize+NormalFontSize)
		op.ColorScale.ScaleWithColor(color.White)
		op.LineSpacing = SmallFontSize
		text.Draw(screen, msg, &text.GoTextFace{
			Source: MPlusFaceSource,
			Size:   SmallFontSize,
		}, op)
		// draw in play cards
		for i, c := range g.inPlayWork {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(i*ArtSmallWidth), InPlayWorkY)
			screen.DrawImage(c.artSmall, op)
		}
		for i, c := range g.inPlayFaith {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(i*ArtSmallWidth), InPlayFaithY)
			screen.DrawImage(c.artSmall, op)
		}
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
		if art := g.myCards.inHand(cursorX, cursorY); art != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(displayX), 0)
			screen.DrawImage(art, op)
		} else if vp := g.kingdom.In(cursorX, cursorY); vp != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(displayX), 0)
			screen.DrawImage(vp.c.artBig, op)
			drawTextBox(screen, cursorX, cursorY, strconv.Itoa(vp.n))
		} else if n := g.myCards.inDeck(cursorX, cursorY); n != -1 {
			drawTextBox(screen, cursorX, cursorY, strconv.Itoa(n))
		} else if art, n := g.myCards.inDiscard(cursorX, cursorY); n != -1 {
			if art != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(displayX), 0)
				screen.DrawImage(art, op)
			}
			drawTextBox(screen, cursorX, cursorY, strconv.Itoa(n))
		}
	}
}

func drawTextBox(dst *ebiten.Image, x, y int, msg string) {
	vector.DrawFilledRect(dst, float32(x), float32(y), SmallFontSize, SmallFontSize, color.Black, true)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = SmallFontSize
	text.Draw(dst, msg, &text.GoTextFace{
		Source: MPlusFaceSource,
		Size:   SmallFontSize,
	}, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	g := &Game{state: RoomName, t: Typewriter{}, players: make(map[string]PlayerData), phase: WorkPhase, ts: TurnStats{works: 1, blessings: 1, faith: 0}}

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
	producerSend(g.pc.producer, []string{LeftLobby, g.pc.playerName})
}
