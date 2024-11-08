// sets up the pulsar client for the player
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"log"

	"github.com/apache/pulsar-client-go/pulsar"
)

// message types
const (
	JoinedLobby  = "J"
	LeftLobby    = "L"
	ToggledReady = "TR"
	SetKingdom   = "SK"
	EndPhase     = "E"
	Gained       = "G"
)

type PulsarClient struct {
	roomName, playerName string
	client               pulsar.Client
	producer             pulsar.Producer
	consumer             pulsar.Consumer
	// tableView            pulsar.TableView
	// consumeCh            chan pulsar.ConsumerMessage
	// // exclude type
	// exclusiveObstacleConsumer pulsar.Consumer
	// // to read the latest obstacle graph
	// obstacleReader pulsar.Reader
	// // subscribe the obstacle topic,
	// closeCh chan struct{}
}

func (c *PulsarClient) Close() {
	if err := c.consumer.Unsubscribe(); err != nil {
		log.Fatal(err)
	}
	c.producer.Close()
	c.consumer.Close()
	c.client.Close()
	// c.closeCh <- struct{}{}
	// c.tableView.Close()
	// close(c.closeCh)
	// close(c.consumeCh)
}

func newPulsarClient(roomName, playerName string) *PulsarClient {
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

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "persistent://public/default/" + roomName,
	})

	if err != nil {
		log.Fatal(err)
	}

	producerSend(producer, []string{JoinedLobby, playerName, strconv.Itoa(rand.Intn(10))})

	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:                       "persistent://public/default/" + roomName,
		SubscriptionName:            playerName,
		SubscriptionInitialPosition: pulsar.SubscriptionPositionEarliest,
	})

	if err != nil {
		log.Fatal(err)
	}

	return &PulsarClient{roomName, playerName, client, producer, consumer}
}

func producerSend(producer pulsar.Producer, message []string) {
	msg, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	if msgId, err := producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: msg,
	}); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Published message: %v \n", msgId)
	}
}

func consumerReceive(consumer pulsar.Consumer) []string {
	msg, err := consumer.Receive(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Received message msgId: %v -- content: '%s'\n",
		msg.ID(), string(msg.Payload()))
	consumer.Ack(msg)
	message := &[]string{}
	json.Unmarshal(msg.Payload(), message)
	return *message
}
