package scan

import (
	"fmt"
	"net/http"
)

func worker(address, results chan Address) {
	for addr := range address {
		url := fmt.Sprintf("http://%s:%d", addr.Host, addr.Port)
		addr.status = false
		resp, err := http.Get(url)
		if err != nil {
			results <- addr
			continue
		}
		resp.Body.Close()

		if resp.Header.Get("Server") == "nginx/1.18.0 (Ubuntu)" {
			addr.status = true
		}
		results <- addr
	}
}
