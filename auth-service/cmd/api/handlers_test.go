package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)


type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func  NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func Test_Authenticate(t *testing.T) {

	jsonToReturn := `
	{
		"error": false,
		"message": "some message",
		"data": null
	}
	`

	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Body: io.NopCloser(bytes.NewBufferString(jsonToReturn)),
			Header: make(http.Header),
		}
	})

	testApp.Client = client
	postBody := map[string]any{
		"email": "me@here.com",
		"password": "verysecret",
	}

	jsonData, _ := json.Marshal(postBody)

	request, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))

	rr := httptest.NewRecorder()
	
	handler := http.HandlerFunc(testApp.Authenticate)

	handler.ServeHTTP(rr, request)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status code of %d, but got %d", http.StatusAccepted, rr.Code)
	}

}

