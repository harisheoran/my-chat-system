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
- Kafka producing and consuming
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

// subscribe to the Redis channel and pass to the broadcast function
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
// subscribe to the redis channel and produce it to kafka
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

		// publish to Kafka broker
		err = app.kafkaProducer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(fmt.Sprintf("message-%s", time.Now().Format("2006-01-02"))),
			Value: []byte(message.Payload),
		})
		if err != nil {
			app.errorlogger.Println("Unable to produce to kafka", err)
		} else {
			app.infologger.Println("message produced to kafka", message.Payload)
		}
	}

}

// consume from kafka and save to the database
func (app *app) consumeFromKafka() {

	for {
		consumedMessage, err := app.kafkaConsumer.ReadMessage(context.Background())
		if err != nil {
			app.errorlogger.Println("Could not read message: ", err)
		} else {
			app.infologger.Println("message comnsumed from kafka", consumedMessage.Value)
		}

		this := string(consumedMessage.Value[0])
		fmt.Printf("THIS %s", this)

		/*
			messageToSave := model.Message{
				Msg: string(consumedMessage.Value),
			}

			err = app.messageController.InsertMessage(&messageToSave)
			if err != nil {
				// send json response also here
				//
				//
				app.errorlogger.Println("unable to save the consumed messsage into the database")
			}
		*/

		app.infologger.Println("consumed message saved into db successfully")
	}

}
