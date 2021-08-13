package storage

import (
	"bytes"
	"fmt"
	"net/http"
)

func SetRedisClient(serviceURL string) (*RedisClient, error) {
	cr := &RedisClient{
		url: serviceURL,
	}
	return cr, nil
}

// func
type RedisClient struct {
	url string
}

func (cr RedisClient) Writer(b *bytes.Buffer) error {
	// b := bytes.NewBuffer([]byte(data))
	res, err := http.Post(cr.url+"/write", "application/json", b)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send redis. Service responded with %v", res.StatusCode)
	}
	return nil
}
