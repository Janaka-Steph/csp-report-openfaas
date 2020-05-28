package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Handle CSP report
func Handle(w http.ResponseWriter, r *http.Request) {
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()
		// read request payload
		reqBody, err := ioutil.ReadAll(r.Body)
		checkError(w, err)
		input = reqBody
	}

	// log to stdout
	// TODO: not printing
	fmt.Printf("request body: %s", string(input))

	type Message struct {
		Payload string              `json:"payload"`
		Headers map[string][]string `json:"headers"`
	}

	type Response struct {
		FromName  string  `json:"from_name"`
		FromEmail string  `json:"from_email"`
		ToEmail   string  `json:"to_email"`
		Message   Message `json:"message"`
		Subject   string  `json:"subject"`
	}

	response := Response{
		Subject:   `CSP Violation Report`,
		FromName:  `CSP Violation OpenFaaS Function`,
		FromEmail: `function@email.com`,
		ToEmail:   `rstephane@protonmail.com`,
		Message: Message{
			Payload: string(input),
			Headers: r.Header,
		},
	}

	resBody, err := json.Marshal(response)
	checkError(w, err)

	// Send POST request
	req, _ := http.NewRequest("POST", "https://openfaas.bitcoin-studio.com/function/send-email", bytes.NewBuffer(resBody))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	checkError(w, err)
	defer res.Body.Close()

	// write result
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}

func checkError(w http.ResponseWriter, err error) {
	if err != nil {
		log.Printf("panic: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
