package scan

import (
	"fmt"
	"net/http"
	"time"
)

func segmentScan(segmentHost string) {
	
}

func worker(address, results chan Address) {
	for addr := range address {
		url := fmt.Sprintf("http://%s:%d", addr.Host, addr.Port)
		fmt.Printf("start %s\n", url)
		addr.status = false
		client := http.Client{Timeout: time.Second}
		resp, err := client.Get(url)
		if err != nil {
			results <- addr
			continue
		}

		if resp.Header.Get("Server") == "ZFSOFT.Inc" {
			addr.status = true
		}
		results <- addr
	}
}
