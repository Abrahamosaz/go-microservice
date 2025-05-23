package main

import (
	"context"
	"log"
	"logger/data"
	"time"
)

type RPCServer struct {
}

type RPCPayload struct {
	Name string
	Data string
}


func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Println("Error inserting into log: ", err)
		return err
	}


	log.Println("Payload processed via RPC:", payload)

	*resp = "Processed payload via RPC:" + payload.Name
	return nil
}