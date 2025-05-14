package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)


func main() {
	// try to connect to rabbitmq
	rabbitMQConn, err := connectRabbitMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitMQConn.Close()

	log.Println("Connected to RabbitMQ")

	// create a consumer
	consumer, err := event.NewConsumer(rabbitMQConn)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume messages
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}

}




func connectRabbitMQ() (*amqp.Connection, error) {

	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	amqpUrl := "amqp://guest:guest@rabbitmq:5672"

	for {
		c, err := amqp.Dial(amqpUrl)
		if err != nil {
			fmt.Println("RabbitMQ not ready yet")
			counts++
		} else {
			log.Println("Connected to RabbitMQ")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
