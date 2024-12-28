package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/harisheoran/my-chat-system/pkg/model"
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

/*
consume from kafka and save to the database
*/
func (app *app) consumeFromKafka() {
	consumer := app.kafkaConsumer
	// number of messages per batch
	batchSize := 5

	// message buffer
	var batch []kafka.Message
	//	batchStart := time.Now()

	for {
		consumedMessage, err := consumer.FetchMessage(context.Background())
		if err != nil {
			app.errorlogger.Println("Could not read message: ", err)
		} else {
			app.infologger.Println("message consumed from kafka", consumedMessage.Value)
		}

		batch = append(batch, consumedMessage)

		if len(batch) >= batchSize {
			err := app.saveMessageToDatabase(batch)

			if err != nil {
				app.errorlogger.Println("unable to prouduce batch", err)

				// retry here
				continue
			}

			// commit after saving message
			err = consumer.CommitMessages(context.Background(), batch...)
			if err != nil {
				app.errorlogger.Println("Unable to commit message after saving to the database")
			}

			batch = nil
		}

		app.infologger.Println("CONSUMED: ", string(consumedMessage.Value))

		//app.infologger.Println("consumed message saved into db successfully")
	}

}

// save message to database
func (app *app) saveMessageToDatabase(messages []kafka.Message) error {

	fmt.Printf("Processing batch of %d messages\n", len(messages))
	fmt.Println()

	messagesToSave := []model.Message{}

	for _, value := range messages {
		messagesToSave = append(messagesToSave, model.Message{Data: string(value.Value)})
	}

	err := app.messageController.BulkInsertMessage(&messagesToSave)
	if err != nil {
		app.errorlogger.Println("Unable to run bulk insert the data into database ", err)
	}
	return nil
}

/*
Retrieve message from the database
*/
func (app *app) messageHistory() error {
	messages, err := app.messageController.GetLastMessages()
	if err != nil {
		app.errorlogger.Println("unable to retrieve the message for history ", err)
	}

	fmt.Println("DATA is here of length", len(messages))
	for i, value := range messages {
		fmt.Println(i, value.Data)
	}

	return nil
}
