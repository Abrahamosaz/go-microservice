package main

import (
	"broker/event"
	"broker/logs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RequestPaylod struct {
	Action	string 			`json:"action"`
	Auth 	AuthPayload 	`json:"auth,omitempty"`
	Log 	LogPayload 		`json:"log,omitempty"`
	Mail 	MailPayload 	`json:"mail,omitempty"`
}


type MailPayload struct {
	From string `json:"from"`
	To string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}


type AuthPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type RPCPayload struct {
	Name string
	Data string
}


func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse {
		Error: false,
		Message: "Hit the broker",
		Data: nil,
	}

	app.writeJSON(w, http.StatusOK, payload)
}


func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	var requestPayload RequestPaylod

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	log.Println("Request Payload", requestPayload)

	switch requestPayload.Action {
	case "auth":
		app.Authenticate(w, requestPayload.Auth)
	case "log":
		// log via HTTP
		// app.logItem(w, requestPayload.Log)

		// log via RPC
		app.logEventViaRPC(w, requestPayload.Log)

		// log via RabbitMQ
		// app.logEventViaRabbitMQ(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknown action"), http.StatusBadRequest)
	}
}


func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	log.Println("Logging data", string(jsonData))

	logServiceURL := "http://logger-service:5004/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling log service"), http.StatusBadRequest)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Logged"

	app.writeJSON(w, http.StatusAccepted, payload)
}


func (app  *Config) Authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to auth microservice

	jsonData, _ := json.MarshalIndent(a, "", "\t")
	
	log.Println("Authenticating", string(jsonData))
	// call the service
	request, err := http.NewRequest("POST", "http://auth-service:5003/auth/login", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	
	fmt.Println("Response from Authenticating", response)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	
	defer response.Body.Close()
	
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"), http.StatusUnauthorized)
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	log.Println("Decoded Response from Authenticating", jsonFromService)

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}


func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	log.Println("Sending mail", string(jsonData))

	mailServiceURL := "http://mail-service:5005/send"

	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	
	defer response.Body.Close()
	
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mail service"), http.StatusBadRequest)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Mail sent to " + msg.To

	app.writeJSON(w, http.StatusAccepted, payload)
}


func (app *Config) logEventViaRabbitMQ(w http.ResponseWriter, l LogPayload) {

	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Logged via RabbitMQ"

	app.writeJSON(w, http.StatusAccepted, payload)
}


func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.RabbitMQ)

	if err != nil {
		return err
	}


	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	jsonData, _ := json.MarshalIndent(payload, "", "\t")
	
	err = emitter.Push(string(jsonData), "log.INFO")
	if err != nil {
		return err
	}
	
	return nil
}

func (app *Config) logEventViaRPC(w http.ResponseWriter, l LogPayload) {
	// create a new RPC client
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	rpcPayload := RPCPayload(l)

	log.Println("RPC Payload", rpcPayload)
	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
	}


	log.Println("Result from logger service", result)
	payload := jsonResponse{
		Error: false,
		Message: result,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}


func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {

	var requestPayload RequestPaylod

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}


	conn, err := grpc.NewClient("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	defer conn.Close()

	client := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = client.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})

	if err != nil {
		fmt.Println("Error writing to log", err)	
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Logged via GRPC"

	app.writeJSON(w, http.StatusAccepted, payload)
}