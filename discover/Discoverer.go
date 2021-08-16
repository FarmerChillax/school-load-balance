package discover

import (
	"balance/registry"
	"balance/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Addr struct {
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	Status   bool   `json:"-"`
	Protocol string `json:"Protocol"`
	timeout  time.Duration
	ssl      bool
}

type Addrs []Addr

// // implement encoding.BinaryMarshaler
// // MarshalBinary use msgpack
func (s *Addr) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// // UnmarshalBinary use msgpack
func (s *Addr) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

// fast discover
func commonPortDiscover(host string) {
	commonPorts := []int{80, 443, 344, 4000, 5000, 8080, 8000, 8888}
	addrs := Addrs{}
	ports := make(chan Addr, 20)
	results := make(chan Addr)
	// 启动worker
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}
	// 发送内容
	go func() {
		for _, port := range commonPorts {
			ports <- Addr{
				Host:     host,
				Port:     port,
				Protocol: "http",
				timeout:  time.Second * 3,
				ssl:      false,
			}
		}
	}()
	// 接收
	for i := 0; i < len(commonPorts); i++ {
		workerRes := <-results
		if workerRes.Status {
			addrs = append(addrs, workerRes)
		}
	}
	if len(addrs) <= 0 {
		return
	}
	err := sendDiscover(addrs)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// full config discover (main discover)
func discoverer(protocol, host string, start, end int, timeout time.Duration, ssl bool) {
	var addrs Addrs
	ports := make(chan Addr, end-start)
	results := make(chan Addr)
	// 启动worker
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}
	// 发送地址
	go func() {
		for i := start; i < end; i++ {
			ports <- Addr{
				Host:     host,
				Port:     i,
				Protocol: protocol,
				timeout:  timeout,
				ssl:      ssl,
			}
		}
		close(ports) // 关闭发送频道
	}()

	// 处理结果
	for i := start; i < end; i++ {
		workerRes := <-results
		if workerRes.Status {
			addrs = append(addrs, workerRes)
		}
	}
	if len(addrs) == 0 {
		fmt.Println("host not found.")
		return
	}
	// 发送结果给redis
	fmt.Println("sending results to redis...")
	err := sendDiscover(addrs)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func sendDiscover(addrs Addrs) error {
	fmt.Println(addrs)
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(addrs)
	if err != nil {
		return err
	}
	// 写数据库（redis）
	redisURL, err := registry.GetProvide(registry.RedisService)
	if err != nil {
		return err
	}

	res, err := http.Post(redisURL+"/write", "application/json", buf)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to write results, response code with %v", res.StatusCode)
	}

	return nil
}

func worker(raw, results chan Addr) {
	// HTTPClient := makeHTTPClient(addr.timeout, addr.ssl)
	for addr := range raw {
		url := fmt.Sprintf("%s://%s:%d", addr.Protocol, addr.Host, addr.Port)
		addr.Status = false
		HTTPClient := utils.NewHTTPClient(addr.timeout, addr.ssl)
		fmt.Printf("Start scanning %s\n", url)
		res, err := HTTPClient.Get(url)
		if err != nil {
			results <- addr
			continue
		}
		if res.StatusCode == http.StatusOK {
			addr.Status = true
		}
		results <- addr
	}
}
