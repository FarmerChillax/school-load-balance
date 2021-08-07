package scan

import "time"

// 深度搜索教务系统
// 扫描网段中主机的每个端口（65535个）
func DeepScanHandler(hostList []string) OpenPorts {
	address := make(chan Address, 30000)
	results := make(chan Address, 1024)
	scanEndPort := 65535
	// new worker
	for i := 0; i < cap(address); i++ {
		go worker(address, results)
	}

	// push value in chan
	go func() {
		for _, item := range hostList {
			pushAddress(item, scanEndPort, address)
		}
	}()
	// collect worker results
	openPorts := resultsHandler(len(hostList), scanEndPort, results)

	close(address)
	close(results)
	return openPorts
}

// 常规端口扫描
func FastScanHandler(hostList []string) OpenPorts {
	address := make(chan Address, 30000)
	results := make(chan Address, 1024)
	commonPort := []Port{80, 443, 344, 4000, 5000, 8080, 8000, 8888}
	// new worker
	for i := 0; i < cap(address); i++ {
		go worker(address, results)
	}
	// push value in chan
	go func() {
		for _, host := range hostList {
			for _, port := range commonPort {
				address <- Address{Host: Host(host), Port: port, status: false}
			}
			time.Sleep(time.Second)
		}
	}()

	// collect worker results
	openPorts := resultsHandler(len(hostList), len(commonPort), results)
	close(address)
	close(results)
	return openPorts
}

// 处理worker结果
func resultsHandler(hostsLen, portsLen int, results chan Address) OpenPorts {
	res := make(OpenPorts)
	for i := 0; i < hostsLen; i++ {
		for j := 1; j <= portsLen; j++ {
			resAddr := <-results
			if resAddr.status {
				if _, ok := res[resAddr.Host]; !ok {
					res[resAddr.Host] = []Port{resAddr.Port}
				} else {
					res[resAddr.Host] = append(res[resAddr.Host], resAddr.Port)
				}
			}
		}
	}
	return res
}

// 推送地址进worker
func pushAddress(host string, scanEndPort int, address chan Address) {
	for i := 1; i <= scanEndPort; i++ {
		address <- Address{
			Host:   Host(host),
			Port:   Port(i),
			status: false,
		}
	}
}
