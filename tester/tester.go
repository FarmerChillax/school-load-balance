package tester

import (
	"balance/discover"
	"balance/registry"
	"balance/utils"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	TEST_BATCH   = 10
	TEST_TIMEOUT = 2
	TEST_SSL     = false
)

func tester() {
	recordCount, err := getCount()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// testTarget := make(chan Addr, recordCount/2)
	// results := make(chan Addr)
	fmt.Println(recordCount)
}

func getCount() (int, error) {
	redisURL, err := registry.GetProvide(registry.RedisService)
	if err != nil {
		return 0, err
	}
	resp, err := http.Get(redisURL + "/utils")
	if err != nil {
		return 0, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	count := binary.BigEndian.Uint64(respBody)
	fmt.Printf("%d %T\n", count, count)
	// return
	return 0, nil
}

// 打开网络io，测试地址
func testService() {
	fmt.Println("start tester.")
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
	time.Sleep(time.Second * 5)
	fmt.Println("push value success.")
	for i := 0; i < addrsCount; i++ {
		addrStatus := <-result
		if addrStatus.Status {
			// 验证成功，设置成满分
			fmt.Println(addrStatus, "is ok.")
		} else {
			fmt.Println(addrStatus, "not ok.")
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
