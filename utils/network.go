package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func RestJson(obj interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize value: %q", err)
	}
	return b.Bytes(), nil
}

func MakeRequest(obj interface{}) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to encode value: %q", err)
	}
	return b, err

}

func NewHTTPClient(timeout time.Duration, ssl bool) (client http.Client) {
	tr := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: ssl},
	}
	client = http.Client{
		Timeout:   timeout,
		Transport: &tr,
	}
	return client
}
