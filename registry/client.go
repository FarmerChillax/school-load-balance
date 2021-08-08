package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// func ShutdownService(url string) error {
// 	// req, err := http.NewRequest(http.MethodDelete, )
// }

func RegisterService(r Registration) error {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(r)
	if err != nil {
		return err
	}
	res, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to register service. Registry service "+
			"responded with code %v", res.StatusCode)
	}
	return nil
}
