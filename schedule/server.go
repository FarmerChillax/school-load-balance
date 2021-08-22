package schedule

import (
	"bytes"
	"fmt"
	"net/http"
)

// Consul Config
type Consul struct {
	Protocol string          `json:"protocol"`
	Host     string          `json:"host"`
	Port     int             `json:"port"`
	Path     string          `json:"path"`
	Services []ConsulService `json:"services"`
}

type ConsulService struct {
	Name  string `json:"name"`
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func (c Consul) sendServices() error {
	// protocol + host + port + path
	consulAddres := fmt.Sprintf("%s://%s:%d%s", c.Protocol, c.Host, c.Port, c.Path)
	for _, item := range c.Services {
		// url + name
		addres := fmt.Sprintf("%s/%s/%s", consulAddres, item.Name, item.Key)
		buf := bytes.NewBuffer([]byte(item.Value))
		req, err := http.NewRequest(http.MethodPut, addres, buf)
		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "text/plain")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to registry consul. Registry service responded with code %v. URL: %v", res.StatusCode, addres)
		}
	}
	return nil
}
