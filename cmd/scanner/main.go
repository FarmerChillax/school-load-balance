package main

import (
	"balance/network"
	"fmt"
	"time"
)

func main() {
	template := "172.16.1.%d"
	startTime := time.Now()
	// do something
	network.StartSegment(template, 60, 100)
	// end time
	elsp := time.Since(startTime)
	
	fmt.Printf("Scan used %ds.\n", elsp/1e9)
}
