package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func MakeResponse(obj interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize value: %q", err)
	}
	return b.Bytes(), nil
}
