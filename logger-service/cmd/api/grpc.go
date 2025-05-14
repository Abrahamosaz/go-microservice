package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"logger/logs"
	"net"

	"google.golang.org/grpc"
)


type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}


func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {

	input := req.GetLogEntry();

	log.Println("input payload from grpc client: ", input)

	// write the log
	logEntry := data.LogEntry {
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		log.Println("Error writing log: ", err)
		return &logs.LogResponse{
			Result: "Failed",
		}, nil
	}

	return &logs.LogResponse{
		Result: "Logged",
	}, nil
}


func (app *Config) gRpcListen() {

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for grpc: %v", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{
		Models: app.Models,
	})

	log.Printf("Starting gRPC server on port %s", gRpcPort)

	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to listen for grpc: %v", err)
	}
}