package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	webPort = "5004"
	rpcPort = "5001"
	gRpcPort = "50001"
	mongoURL = "mongodb://mongo:27017"
)


var client *mongo.Client

type Config struct {
	Models data.Models
	
}



func main() {
	//connect to mongoDB
	// dsn := os.Getenv("MONGO_URI")
	dsn := "mongodb://admin:password@logger-mongo:27017/"
	mongoClient, err := connectToMongo(dsn)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Connected to mongoDB")

	client = mongoClient

	//create a context in order to disconnect from mongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// close connection to mongoDB
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// Register the RPC server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		log.Panic(err)
	}

	go app.rpcListen()

	log.Println("Starting logger service on port", webPort)
	// app.serve()
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	go app.gRpcListen()
	
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) rpcListen() {
	log.Println("Starting RPC server on port", rpcPort)

	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		log.Panic(err)
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}

		go  rpc.ServeConn(rpcConn)
	}
}

func connectToMongo(dsn string) (*mongo.Client, error) {	
	clientOptions := options.Client().ApplyURI(dsn)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Println("Error connecting to mongoDB", err)
		return nil, err
	}

	return client, nil
}

