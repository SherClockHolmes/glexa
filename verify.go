package glexa

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const verifyURL = "s3.amazonaws.com"

func verifyCertURL(certURL string) error {
	parsed, err := url.Parse(certURL)
	if err != nil {
		return fmt.Errorf("Could not parse SignatureCertChainUrl: %q\n", err)
	}

	if parsed.Scheme != "https" {
		return fmt.Errorf("Scheme is not https: %q\n", parsed.Scheme)
	}

	if host, port, err := net.SplitHostPort(parsed.Host); err == nil {
		if port != "443" || host != verifyURL {
			return fmt.Errorf("Invalid hostname or port")
		}
	}

	if !strings.HasPrefix(strings.ToLower(parsed.Host), verifyURL) {
		return fmt.Errorf("Invalid hostname")
	}

	if !strings.HasPrefix(parsed.Path, "/echo.api/") {
		return fmt.Errorf("Invalid path")
	}

	return nil
}

func verifyBodyTimestamp(body io.Reader) error {
	// check that body is not expired
	bodyStruct := struct {
		Request struct {
			Timestamp string `json:"timestamp"`
		} `json:"request"`
	}{}

	if err := json.NewDecoder(body).Decode(&bodyStruct); err != nil {
		return fmt.Errorf("Could not decode body: %q\n", err)
	}

	requestTime, err := time.Parse("2006-01-02T15:04:05Z", bodyStruct.Request.Timestamp)
	if err != nil {
		return fmt.Errorf("Could not parse timestamp: %q\n", err)
	}

	if time.Now().Sub(requestTime) > time.Second*150 {
		return fmt.Errorf("Timestamp is stale")
	}
	return nil
}

func validateCertChain(chainURL string) (*x509.Certificate, error) {
	resp, err := http.Get(chainURL)
	if err != nil {
		return nil, fmt.Errorf("Could not get cert chain pem: %q\n", err)
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read cert chain pem: %q\n", err)
	}

	block, _ := pem.Decode(buf)
	if err != nil {
		return nil, fmt.Errorf("Could not decode cert chain: %q\n", err)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Could not parse cert chain: %q\n", err)
	}

	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(buf); !ok {
		return nil, fmt.Errorf("Could not parse root cert: %q\n", err)
	}

	opts := x509.VerifyOptions{
		DNSName: "echo-api.amazon.com",
		Roots:   roots,
	}

	if _, err := cert.Verify(opts); err != nil {
		return nil, fmt.Errorf("Could not verify cert chain: %q\n", err)
	}
	return cert, nil
}

func verifySignature(signature string, pubKey *rsa.PublicKey, body io.Reader) error {
	data, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("Could not base64 decode signature: %q\n", err)
	}

	buf, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("Could not read request body: %q\n", err)
	}

	hashed := sha1.Sum(buf)
	if err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA1, hashed[:], data); err != nil {
		return fmt.Errorf("Cerification error: %q\n", err)
	}

	return nil
}
