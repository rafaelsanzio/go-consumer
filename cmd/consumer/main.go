package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/rafaelsanzio/go-full-cycle/part1/internal/infra/database"
	"github.com/rafaelsanzio/go-full-cycle/part1/internal/usecase"
	"github.com/rafaelsanzio/go-full-cycle/part1/pkg/kafka"
	"github.com/rafaelsanzio/go-full-cycle/part1/pkg/rabbitmq"

	_ "github.com/mattn/go-sqlite3"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	db, err := sql.Open("sqlite3", "./orders.db")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	repository := database.NewOrderRepository(db)
	usecase := usecase.CalculateFinalPrice{OrderRepository: repository}

	msgChanKafka := make(chan *ckafka.Message)
	topics := []string{"orders"}
	servers := "host.docker.internal:9094"

	go kafka.Consume(topics, servers, msgChanKafka)
	go kafkaWorker(msgChanKafka, usecase)

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	msgRabbitmqChannel := make(chan amqp.Delivery)
	go rabbitmq.Consume(ch, msgRabbitmqChannel)
	rabbitmqWorker(msgRabbitmqChannel, usecase)
}

func kafkaWorker(msgChan chan *ckafka.Message, uc usecase.CalculateFinalPrice) {
	fmt.Println("Kafka worker has started")

	for msg := range msgChan {
		var orderInputDTO usecase.OrderInputDTO
		err := json.Unmarshal(msg.Value, &orderInputDTO)
		if err != nil {
			panic(err)
		}

		orderOutputDTO, err := uc.Execute(orderInputDTO)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Kafka has processed order %s\n", orderOutputDTO.ID)
	}
}

func rabbitmqWorker(msgChan chan amqp.Delivery, uc usecase.CalculateFinalPrice) {
	fmt.Println("Rabbitmq worker has started")

	for msg := range msgChan {
		var OrderInputDTO usecase.OrderInputDTO
		err := json.Unmarshal(msg.Body, &OrderInputDTO)
		if err != nil {
			panic(err)
		}

		orderOutputDTO, err := uc.Execute(OrderInputDTO)
		if err != nil {
			panic(err)
		}

		msg.Ack(false)
		fmt.Printf("Rabbitmq has processed order %s\n", orderOutputDTO.ID)
	}
}
