// sets up the pulsar client for the player
package main

import (
	"context"
	"fmt"

	"log"

	"github.com/apache/pulsar-client-go/pulsar"
)

// message types
const (
	PlayerJoined       = "Player joined: "
	PlayerLeft         = "Player left: "
	PlayerToggledReady = "Player toggled ready: "
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
	if msgId, err := c.producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: []byte(PlayerLeft + c.playerName),
	}); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Published message: %v \n", msgId)
	}
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

	producerSend(producer, PlayerJoined+playerName)

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

func producerSend(producer pulsar.Producer, message string) {
	if msgId, err := producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: []byte(message),
	}); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Published message: %v \n", msgId)
	}
}
