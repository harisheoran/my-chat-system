package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/harisheoran/my-chat-system/pkg/model"
	postgre "github.com/harisheoran/my-chat-system/pkg/model/postgre"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const version = "1.0.0"

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all origins (not secure, only for testing)
			return true
		},
	}

	// redis channel
	myChannel = "common-room-new"
	ctx       = context.Background()

	//This map keeps track of connected WebSocket clients.
	client = make(map[*websocket.Conn]bool)

	// publish channel
	publishChannel = make(chan Message)

	// subscribe channel
	broadcastChannel = make(chan string)

	// kafka Channel
	kafkaChannel = make(chan string)

	// for kafka
	TOPIC_NAME = "TEST"
	producer   *kafka.Writer
	consumer   *kafka.Reader
	kafkaUrl   = "kafka-15e11ee5-project-54d1.c.aivencloud.com:25415"
)

func main() {
	// create two loggers for info and error
	errorlogger := log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime)
	infologger := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)

	// load .env file
	err := godotenv.Load()
	if err != nil {
		errorlogger.Println("can't read the env files")
	}
	dsn := os.Getenv("DBURI")

	// create database connection pool
	databaseConnection, err := createDbConnectionPool(dsn)
	if err != nil {
		errorlogger.Println("unable to get a database connection from the pool", err)
	}

	/*
		establish the redis connection
	*/
	redis_db, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 2, 64)
	if err != nil {
		errorlogger.Println("unable to parse the redis db value from .env file")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       int(redis_db),
	})
	// check redis is connected or not
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		errorlogger.Println("unable to connect to redis", err)
	}

	// Runtime api configs
	appConfig := AppConfig{}
	flag.IntVar(&appConfig.port, "port", 1316, "runtime port for the api")
	flag.StringVar(&appConfig.env, "env", "test", "runtime environment name")
	flag.Parse()

	app := app{
		infologger:      infologger,
		errorlogger:     errorlogger,
		redisConnection: rdb,
		messageController: &postgre.MessageController{
			DbConnection: databaseConnection,
		},
		userController: &postgre.UserController{
			DbConnection: databaseConnection,
		},
		channelController: &postgre.ChannelController{
			DbConnection: databaseConnection,
		},
		kafkaProducer: createKafkaProducer(),
		kafkaConsumer: createKafkaConsumer(),
		appConfig:     &appConfig,
	}

	// start the chat worflow
	go app.publishToRedis()
	go app.subscribeToRedis()
	go app.broadcastMessages()
	go app.produceToKafka()
	go app.consumeFromKafka()

	// start the server
	port := fmt.Sprintf(":%d", app.appConfig.port)
	server := &http.Server{
		Addr:     port,
		ErrorLog: app.errorlogger,
		Handler:  app.router(),
	}

	app.infologger.Println("Server is starting on the port ", port)
	err = server.ListenAndServe()
	if err != nil {
		app.errorlogger.Fatal("unable to start the server on port 1316", err)
	}
}

// establish the db connection
func createDbConnectionPool(dsn string) (*gorm.DB, error) {
	dbConnectionPool, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	// Run the automigration for Project Model
	if err := dbConnectionPool.AutoMigrate(&model.Message{}, &model.User{}, &model.Channel{}); err != nil {
		return nil, err
	}

	return dbConnectionPool, nil
}

func kafkaInitialize() (*kafka.Dialer, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("can't read the env files")
	}

	var username = os.Getenv("KAFKA_USERNAME")
	var password = os.Getenv("KAFKA_PASSWORD")

	caCert, err := os.ReadFile("ca.pem")
	if err != nil {
		log.Println("Failed to read CA certificate file: ", err)
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		log.Println("Failed to parse CA certificate file: ", err)
		return nil, err
	}

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}
	scram, err := scram.Mechanism(scram.SHA512, username, password)
	if err != nil {
		log.Println("Failed to create scram mechanism: ", err)
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

// create a kafka producer
func createKafkaProducer() *kafka.Writer {
	dailer, err := kafkaInitialize()
	if err != nil {
		log.Println("unable to authenticate or initialize with Kafka: ", err)
	}

	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaUrl},
		Topic:    TOPIC_NAME,
		Balancer: &kafka.Hash{},
		Dialer:   dailer,
	})

	return producer
}

// create a kafka consumer
func createKafkaConsumer() *kafka.Reader {

	dialer, err := kafkaInitialize()
	if err != nil {
		log.Println("unable to authenticate or initialize with Kafka: ", err)
	}

	// init consumer
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaUrl},
		Topic:          TOPIC_NAME,
		GroupID:        "main-chat-consumer-group",
		Dialer:         dialer,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		MaxWait:        1 * time.Second,
		CommitInterval: 0, // no auto-commit
	})

	return consumer
}
