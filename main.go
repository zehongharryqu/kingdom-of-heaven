package main

import (
	"context"
	"fmt"
	"image"
	"log"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/hajimehoshi/ebiten/v2"
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
	state string
	t     Typewriter
	h     *Hand
}

func (g *Game) Update() error {
	switch g.state {
	case RoomName:
		g.t.Update()
		if g.t.confirmedName != "" && g.t.confirmedRoom != "" {
			g.state = Playing // Todo: implement Lobby
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case RoomName:
		g.t.Draw(screen)
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
	oauth := pulsar.NewAuthenticationOAuth2(map[string]string{
		"type":       "client_credentials",
		"issuerUrl":  "https://auth.streamnative.cloud/",
		"audience":   "urn:sn:pulsar:o-hwa6o:kingdom-of-heaven-instance",
		"privateKey": "file:///Users/harry/Downloads/o-hwa6o-harry.json", // Absolute path of your downloaded key file
	})

	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:            "pulsar+ssl://pc-de347430.gcp-shared-usce1.g.snio.cloud:6651",
		Authentication: oauth,
	})

	if err != nil {
		log.Fatalf("Could not instantiate Pulsar client: %v", err)
	}

	defer client.Close()

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "persistent://public/default/koh-topic",
	})

	if err != nil {
		log.Fatal(err)
	}

	defer producer.Close()

	for i := 0; i < 10; i++ {
		if msgId, err := producer.Send(context.Background(), &pulsar.ProducerMessage{
			Payload: []byte(fmt.Sprintf("hello-%d", i)),
		}); err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("Published message: %v \n", msgId)
		}
	}

	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:                       "persistent://public/default/koh-topic",
		SubscriptionName:            "test-sub",
		SubscriptionInitialPosition: pulsar.SubscriptionPositionEarliest,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer consumer.Close()

	for i := 0; i < 10; i++ {
		msg, err := consumer.Receive(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Received message msgId: %v -- content: '%s'\n",
			msg.ID(), string(msg.Payload()))

		consumer.Ack(msg)
	}

	c := &Card{assets.CardBig, assets.CardSmall, nil}
	h := &Hand{[]*Card{c, c, c, c, c}}
	g := &Game{RoomName, Typewriter{}, h}

	err = ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}

	if err := consumer.Unsubscribe(); err != nil {
		log.Fatal(err)
	}
}
