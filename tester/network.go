package tester

import (
	"balance/discover"
	"balance/registry"
	"balance/storage"
	"bytes"
	"encoding/json"
	"net/http"
)

// 获取redis中的 Addr
func getAddrs() (addrs discover.Addrs, err error) {
	batch := storage.Batch{
		Cursor: 0,
		Match:  "",
		Count:  TEST_BATCH,
	}
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(batch)
	if err != nil {
		return addrs, err
	}
	redisURL, err := registry.GetProvide(registry.RedisService)
	if err != nil {
		return addrs, err
	}
	// 获取redis数据
	resp, err := http.Post(redisURL+"/tester", "application/json", buf)
	if err != nil {
		return addrs, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&addrs)
	if err != nil {
		return addrs, err
	}

	return addrs, nil
}

func GetAddrs() (discover.Addrs, error) {
	return getAddrs()
}
