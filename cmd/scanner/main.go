package main

import (
	"balance/network"
	"fmt"
	"time"
)

// func main() {
// 	start := time.Now()
// 	hostStr := "172.16.1.%d"
// 	hosts := []string{}
// 	for i := 1; i <= 255; i++ {
// 		hosts = append(hosts, fmt.Sprintf(hostStr, i))
// 	}
// 	res := scan.DeepScanHandler(hosts)
// 	// res := scan.FastScanHandler(hosts)
// 	elsp := time.Since(start)
// 	fmt.Println(res)
// 	fmt.Printf("Scan used %ds.\n", elsp/1e9)
// }

func main() {
	// addr := network.Address{
	// 	Host:    "farmer233.top",
	// 	Status:  false,
	// 	Timeout: network.Timeout(time.Second),
	// }
	template := "172.16.1.%d"
	startTime := time.Now()
	// do something
	network.StartSegment(template, 60, 120)
	// end time
	elsp := time.Since(startTime)
	fmt.Printf("Scan used %ds.\n", elsp/1e9)
}
