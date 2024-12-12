package main

import (
	"encoding/json"
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

		// pass message to kafka channel
		kafkaChannel <- message.Payload
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
	for {
		payload := <-kafkaChannel
		// publish to Kafka broker
		err := app.produceMessage(payload)
		if err != nil {
			app.errorlogger.Println("unable to produce to kafka: ", err)
		} else {
			app.infologger.Println("Message published to kafka")
		}
	}
}
