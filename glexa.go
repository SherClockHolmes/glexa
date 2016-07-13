package glexa

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Decode the alexa request body
func Decode(r *http.Request) (*RequestBody, error) {
	body := &RequestBody{}
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		return body, err
	}

	return body, nil
}

// Encode the alexa response body
func Encode(w http.ResponseWriter, sessionAttributes *Attributes, resp *Response) error {
	// TODO: Validate response

	body := &ResponseBody{
		Version:           "1.0",
		SessionAttributes: sessionAttributes,
		Response:          resp,
	}

	return json.NewEncoder(w).Encode(body)
}

// Decorator for verifying if the request comes from AWS
func Verify(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sigCertChainURL := r.Header.Get("SignatureCertChainUrl")

		// Check for valid sig chain url
		if err := verifyCertURL(sigCertChainURL); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		bodyBuf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate timestamp
		if err := verifyBodyTimestamp(bytes.NewBuffer(bodyBuf)); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Vaidate certchain
		cert, err := validateCertChain(sigCertChainURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Verify signature
		signature := r.Header.Get("Signature")
		err = verifySignature(signature, cert.PublicKey.(*rsa.PublicKey), bytes.NewBuffer(bodyBuf))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Reset body
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBuf))

		// Everything checks out, run the handler
		next(w, r)
	}
}
