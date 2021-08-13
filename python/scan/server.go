package scan

import (
	"balance/registry"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func RegisterHandlers() {
	http.HandleFunc("/start", scanHandler)
}

// 请求体格式：
// {
//		hosts: [],
// 		start: <start>,
// 		end: <end>,
// }
func scanHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// 扫描常用端口(默认方法)
	case http.MethodPost:
		// 扫描指定端口
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func FastScan() {
	templateHost := "172.16.1.%d"
	hostList := make([]string, 260)
	for i := 1; i <= 255; i++ {
		hostList = append(hostList, fmt.Sprintf(templateHost, i))
	}
	go fastScan(hostList)
}

// 常规端口扫描
func fastScan(hostList []string) {
	address := make(chan Address, 1024)
	results := make(chan Address)
	commonPort := []int{80, 443, 344, 4000, 5000, 8080, 8000, 8888}
	// new worker
	for i := 0; i < cap(address); i++ {
		go worker(address, results)
	}
	// push value in chan
	go func() {
		for _, host := range hostList {
			for _, port := range commonPort {
				address <- Address{Host: host, Port: port, status: false}
			}
		}
		close(address)
	}()

	// collect worker results
	addrs := resultsHandler(len(hostList)*len(commonPort), results)
	close(results)
	// send res to redis
	err := sendResult(addrs)
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func sendResult(addrs Addrs) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(addrs)
	if err != nil {
		return err
	}
	redisURL, err := registry.GetProvide(registry.RedisService)
	if err != nil {
		return err
	}
	res, err := http.Post(redisURL+"/write", "application/json", buf)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("faild write redis, response status code: %v", res.StatusCode)
	}
	return nil
}

func resultsHandler(resLen int, results chan Address) (res Addrs) {
	for i := 0; i < resLen; i++ {
		ret := <-results
		if ret.status {
			res = append(res, ret)
		}
	}

	return res
}
