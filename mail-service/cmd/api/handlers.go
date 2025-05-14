package main

import (
	"fmt"
	"net/http"
)


type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}



func (app *Config ) sendMail(w http.ResponseWriter, r *http.Request) {
	var requestPayload MailPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}


	fmt.Println("requestPayload", requestPayload)
	msg := Message{
		From: requestPayload.From,
		To: []string{requestPayload.To},
		Subject: requestPayload.Subject,
		Data: requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		fmt.Println("error sending mail", err)
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error: false,
		Message: "Message sent to " + requestPayload.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
