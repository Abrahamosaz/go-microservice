package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "8080"

type Config struct {
	RabbitMQ *amqp.Connection
}


func main() {
	rabbitMQConn, err := connectRabbitMQ()
	if err != nil {
		log.Panic(err)
	}
	defer rabbitMQConn.Close()

	fmt.Println("Connected to RabbitMQ")

	app := Config{
		RabbitMQ: rabbitMQConn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
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