package scan

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func worker(address, results chan Address) {
	for addr := range address {
		url := fmt.Sprintf("http://%s:%d", addr.Host, addr.Port)
		log.Printf("start %s\n", url)
		addr.status = false
		client := http.Client{Timeout: time.Second}
		resp, err := client.Get(url)
		// resp, err := http.Get(url)
		if err != nil {
			results <- addr
			continue
		}
		resp.Body.Close()

		if resp.Header.Get("Server") == "ZFSOFT.Inc" {
			addr.status = true
		}
		results <- addr
	}
}
