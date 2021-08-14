package tester

import (
	"balance/registry"
	"balance/storage"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const TEST_BATCH = 10

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

// 获取redis中的url
func getAddrs(batch storage.Batch) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(batch)
	if err != nil {
		return err
	}
	redisURL, err := registry.GetProvide(registry.RedisService)
	if err != nil {
		return err
	}

	resp, err := http.Post(redisURL+"/utils", "application/json", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var batchResp storage.Resp
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&batchResp)
	if err != nil {
		return err
	}
	// fmt.Printf("%v %T\n", batchResp.Data, batchResp.Data["Results"])
	return nil
}

// func GetAddrs(batch storage.Batch) error {
// 	return getAddrs(batch)
// }

// 打开网络io，测试地址
func testService() error {
	count, err := getCount()
	if err != nil {
		return err
	}
	fmt.Printf("%d proxies to test.\n", count)
	cursor := 0
	for {
		fmt.Printf("testing proxies use cursor %d, count %d\n", cursor, TEST_BATCH)
		// cursor, proxies :=
	}
}
