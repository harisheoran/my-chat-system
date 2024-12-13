package main

/*
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
*/
