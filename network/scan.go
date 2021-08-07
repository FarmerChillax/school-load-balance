package network

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func StartSegment(template string, start, end Port) {
	Seg := Segment{
		Template:    template,
		SegmentType: 'D',
		StartPort:   start,
		EndPort:     end,
	}
	Seg.SegmentHandlerFunc()
}

func (seg Segment) SegmentHandlerFunc() {
	segRes := make(chan SegmentResult, 100)
	for j := 1; j <= 255; j++ {
		addr := Address{
			Host:    Host(fmt.Sprintf(seg.Template, j)),
			Status:  false,
			Timeout: Timeout(time.Second), // timeout这里写死了
		}
		go addr.Handler(seg.StartPort, seg.EndPort, segRes)
	}
	for i := 1; i <= 255; i++ {
		res := <-segRes
		if len(res.Ports) > 0 {
			for _, item := range res.Ports {
				fmt.Printf("%s:%d opend!\n", res.Host, item)
			}
		}
	}
	close(segRes)
}

// port Handler；端口扫描器（单一ip）
func (addr Address) Handler(start, end Port, segRes chan SegmentResult) {
	res := SegmentResult{Host: addr.Host}
	address := make(chan Address, 1024)
	results := make(chan Address, 10)

	for i := 0; i < cap(address); i++ {
		go worker(address, results)
	}

	go addr.pushPort(start, end, address)

	ports := resultsHandler(start, end, results)

	res.Ports = ports
	segRes <- res
	close(address)
	close(results)
}

// 处理扫描结果
func resultsHandler(start, end Port, results chan Address) (ports []Port) {
	for i := start; i <= end; i++ {
		res := <-results
		if res.Status {
			log.Printf("URL: %s:%d find!\n", res.Host, res.Port)
			ports = append(ports, res.Port)
		}
	}
	return ports
}

// 推送端口进channel给worker
func (addr Address) pushPort(startPort, endPort Port, address chan Address) {
	if startPort <= 0 || endPort > 65535 || endPort < startPort {
		log.Fatalln("Port range oversize.")
		os.Exit(1)
	}

	for i := startPort; i <= endPort; i++ {
		addr.Port = i
		address <- addr
	}
}

// ...
func worker(address, results chan Address) {
	for addr := range address {
		url := fmt.Sprintf("http://%s:%d", addr.Host, addr.Port)
		addr.Status = false
		client := http.Client{Timeout: time.Duration(addr.Timeout)}
		resp, err := client.Get(url)
		if err != nil {
			results <- addr
			continue
		}
		resp.Body.Close()
		if resp.Header.Get("Server") == "ZFSOFT.Inc" {
			addr.Status = true
		}
		results <- addr
	}
}
