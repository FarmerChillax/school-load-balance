package scan

// 深度搜索教务系统
// 扫描网段中主机的每个端口（65535个）
// func DeepScanHandler(hostList []string) OpenPorts {
// 	address := make(chan Address, 30000)
// 	results := make(chan Address, 1024)
// 	scanEndPort := 65535
// 	// new worker
// 	for i := 0; i < cap(address); i++ {
// 		go worker(address, results)
// 	}

// 	// push value in chan
// 	go func() {
// 		for _, item := range hostList {
// 			pushAddress(item, scanEndPort, address)
// 		}
// 	}()
// 	// collect worker results
// 	openPorts := resultsHandler(len(hostList), scanEndPort, results)

// 	close(address)
// 	close(results)
// 	return openPorts
// }
