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

// 传入停止服务的url，向注册中心发送注销请求
// url格式为-> ip:port
func ShutdownService(url string) error {
	req, err := http.NewRequest(http.MethodDelete, ServicesURL,
		// request请求体内容->url
		bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to deregister service. Registry service responded with code %v", res.StatusCode)
	}
	return nil
}
