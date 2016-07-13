GLEXA
=====

Glexa is an Alexa Skills webservice written in Go. Integrate Glexa with your go web server to create/host your own
Alexa commands (skills). Forked from [davinche/glexa](https://github.com/davinche/glexa).


Example Handler
-------------

```go
package main

import (
	"log"
	"net/http"

	"github.com/sherclockholmes/glexa"
)

func AlexaHandler(w http.ResponseWriter, r *http.Request) {
	body, err := glexa.Decode(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := &glexa.Response{}

	switch body.Request.Type {
	case glexa.LaunchRequest:
		resp.Tell("I did not understand your command. Please try again.")
	case glexa.IntentRequest:
		resp.Tell("You are awesome!")
	case glexa.SessionEndedRequest:
		resp.Tell("Good Bye!")
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = glexa.Encode(w, nil, resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func main() {
	// Use the glexa verification decorator
	http.HandleFunc("/test", glexa.Verify(AlexaHandler))
    log.Fatal(http.ListenAndServeTLS(":443", "mycert.pem", "mykey.pem", AlexaHandler))
}
```
