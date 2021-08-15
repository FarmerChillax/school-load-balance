package tester

import (
	"balance/discover"
	"balance/registry"
	"balance/utils"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	TEST_BATCH   = 10
	TEST_TIMEOUT = 2
	TEST_SSL     = false
	TEST_RATE    = 10
)

func Start() {
	log.Println("Tester start success, test rate:", TEST_RATE)
	for {
		testService()
		time.Sleep(TEST_RATE * time.Second)
	}
}

// 打开网络io，测试地址
func testService() {
	addrs, err := getAddrs()
	if err != nil {
		log.Fatalln(err)
	}
	addrsCount := len(addrs)
	testAddrs := make(chan discover.Addr, addrsCount)
	result := make(chan discover.Addr)
	for i := 0; i < cap(testAddrs); i++ {
		go worker(testAddrs, result)
	}

	for _, addr := range addrs {
		testAddrs <- addr
	}
	close(testAddrs)

	for i := 0; i < addrsCount; i++ {
		addrStatus := <-result
		if addrStatus.Status {
			// 验证成功，设置成满分
			err := max(addrStatus)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			// 验证失败，降分
			err := decrease(addrStatus)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}

}

func worker(raw, results chan discover.Addr) {
	for addr := range raw {
		url := fmt.Sprintf("%s://%s:%d", addr.Protocol, addr.Host, addr.Port)
		addr.Status = false
		HTTPClient := utils.NewHTTPClient(TEST_TIMEOUT*time.Second, TEST_SSL)
		fmt.Printf("Testing %s\n", url)
		res, err := HTTPClient.Get(url)
		if err != nil {
			results <- addr
			continue
		}
		if res.StatusCode == http.StatusOK {
			addr.Status = true
		}
		results <- addr
	}
}

func decrease(addr discover.Addr) error {
	redisURL, err := registry.GetProvide(registry.RedisService)
	if err != nil {
		return err
	}
	buf, err := addr.MarshalBinary()
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodDelete, redisURL+"/tester", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("faild to decrease with status code: %v", resp.StatusCode)
	}
	return nil
}

func max(addr discover.Addr) error {
	respStatusCode, err := sendRequest(addr, "/tester", http.MethodPut)
	if err != nil {
		return err
	}
	if respStatusCode != http.StatusOK {
		return fmt.Errorf("faild to decrease with status code: %v", respStatusCode)
	}
	return nil
}

func sendRequest(addr discover.Addr, url, method string) (int, error) {
	redisURL, err := registry.GetProvide(registry.RedisService)
	if err != nil {
		return 0, err
	}
	buf, err := addr.MarshalBinary()
	if err != nil {
		return 0, err
	}
	request, err := http.NewRequest(method, redisURL+url, bytes.NewBuffer(buf))
	if err != nil {
		return 0, err
	}
	request.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	return resp.StatusCode, err
}
