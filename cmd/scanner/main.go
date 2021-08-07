package main

import (
	"balance/scan"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	hostStr := "172.16.1.%d"
	hosts := []string{}
	for i := 1; i <= 255; i++ {
		hosts = append(hosts, fmt.Sprintf(hostStr, i))
	}
	// fmt.Println(hosts)
	res := scan.FastScanHandler(hosts)
	elsp := time.Since(start)
	fmt.Println(res)
	fmt.Printf("Scan used %ds.\n", elsp/1e9)
}
