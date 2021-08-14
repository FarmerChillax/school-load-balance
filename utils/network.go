package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
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
