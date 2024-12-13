package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

/*
this file contains logic for realtime chat
- Web Socket listening and broadcasting
- Redis publishing and subscribing
- Kafka producing
*/

// broadcast the message to every connection
func (app *app) broadcastMessages() {
	//	var message Message
	for {
		recieveMessage := <-broadcastChannel

		var message Message
		err := json.Unmarshal([]byte(recieveMessage), &message)
		if err != nil {
			app.errorlogger.Println("Unable to unmarshal message:", err)
			continue
		}

		for clientConnection := range client {
			err := clientConnection.WriteMessage(message.PayloadType, []byte(message.Payload))
			if err != nil {
				app.errorlogger.Println("unable to broadcast the message", err)
				clientConnection.Close()
				delete(client, clientConnection)
			}
		}

	}
}

// publish the message to the Redis channel
func (app *app) publishToRedis() {
	for {
		payload := <-publishChannel

		payloadJson, err := json.Marshal(payload)
		if err != nil {
			app.errorlogger.Println("unable top marshal the payload json")
		}
		err = app.redisConnection.Publish(ctx, myChannel, payloadJson).Err()
		if err != nil {
			panic(err)
		}
		app.infologger.Println("message published to redis")

	}
}

// subscribe to the Redis channel
func (app *app) subscribeToRedis() {
	// There is no error because go-redis automatically reconnects on error.
	pubsub := app.redisConnection.Subscribe(ctx, myChannel)

	// Close the subscription when we are done.
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			app.errorlogger.Println("unable to subscribe to the redis channel", err)
		} else {
			app.infologger.Println("message subscribed from redis")
		}

		// pass message to the broadcast channel
		broadcastChannel <- msg.Payload
	}
}

// produce message to the kafka channel
func (app *app) produceToKafka() {

	// There is no error because go-redis automatically reconnects on error.
	pubsub := app.redisConnection.Subscribe(ctx, myChannel)

	// Close the subscription when we are done.
	defer pubsub.Close()
	defer app.kafkaProducer.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			app.errorlogger.Println("unable to subscribe to the redis channel for kafka", err)
		} else {
			app.infologger.Println("message subscribed from redis for kafka")
		}

		var message Message
		err = json.Unmarshal([]byte(msg.Payload), &message)
		if err != nil {
			app.errorlogger.Println("Unable to unmarshal message:", err)
			continue
		}
		fmt.Println("From Redis subscribed value is: ", message)

		// publish to Kafka broker
		err = app.kafkaProducer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(fmt.Sprintf("message-%s", time.Now().Format("2006-01-02"))),
			Value: []byte(message.Payload),
		})
		if err != nil {
			app.errorlogger.Println("Unable to produce to kafka", err)
		} else {
			app.infologger.Println("MESSAGE SENT: ", message.Payload)
		}
	}

}
