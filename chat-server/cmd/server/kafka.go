package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

/*
Contains all the kafka code
*/

var kafkaUrl = "kafka-my-chat-system-particleasw123-2262.c.aivencloud.com:15563"

func (app *app) kafkaInitialize() (*kafka.Dialer, error) {
	err := godotenv.Load()
	if err != nil {
		app.errorlogger.Println("can't read the env files")
	}

	var username = os.Getenv("KAFKA_USERNAME")
	var password = os.Getenv("KAFKA_PASSWORD")

	caCert, err := os.ReadFile("ca.pem")
	if err != nil {
		app.errorlogger.Println("Failed to read CA certificate file: ", err)
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		app.errorlogger.Println("Failed to parse CA certificate file: ", err)
		return nil, err
	}

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}
	scram, err := scram.Mechanism(scram.SHA512, username, password)
	if err != nil {
		app.errorlogger.Println("Failed to create scram mechanism: ", err)
		return nil, err
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		TLS:           tlsConfig,
		SASLMechanism: scram,
	}

	return dialer, nil
}

// create a producer
// check if there is already a producer then create if not
func (app *app) createProducer() *kafka.Writer {
	if producer != nil {
		return producer
	}

	dailer, err := app.kafkaInitialize()
	if err != nil {
		app.errorlogger.Println("unable to authenticate or initialize with Kafka: ", err)
	}

	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaUrl},
		Topic:    TOPIC_NAME,
		Balancer: &kafka.Hash{},
		Dialer:   dailer,
	})

	return producer
}

func (app *app) produceMessage(message string) error {
	app.infologger.Println("Starting produce message")
	if producer == nil {
		app.infologger.Println("Producer founded nil and created")
		producer = app.createProducer()
	}

	err := producer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(fmt.Sprintf("message-%s", time.Now().Format("2006-01-02"))),
		Value: []byte(message),
	})
	if err != nil {
		return err
	}
	app.infologger.Println("Message sent: " + message)

	//producer.Close()

	return nil
}

func (app *app) startKafkaConsumer() {
	app.infologger.Println("starting the kafka consumer func")

	dialer, err := app.kafkaInitialize()
	if err != nil {
		app.errorlogger.Println("unable to authenticate or initialize with Kafka: ", err)
	}

	// init consumer
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaUrl},
		Topic:   TOPIC_NAME,
		Dialer:  dialer,
	})

	for {
		message, err := consumer.ReadMessage(context.Background())
		if err != nil {
			app.errorlogger.Printf("Could not get message from kafka: %s", err)
		}
		app.infologger.Printf("Got message using SASL: %s", message.Value)

	}

}
