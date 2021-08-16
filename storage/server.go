package storage

import (
	"balance/discover"
	"balance/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Resp struct {
	StatusCode int         `json:"status_code"`
	Msg        string      `json:"message"`
	Data       interface{} `json:"data"`
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var addrs discover.Addrs
		err := decoder.Decode(&addrs)
		fmt.Printf("recive data: %v\n", len(addrs))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// 写数据库
		err = add(addrs)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type Batch struct {
	Cursor uint64
	Match  string
	Count  int64
}

// utils
func testHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// 获取proxies数量
	case http.MethodGet:
		recordCount := fmt.Sprintf("%d", count())
		w.Write([]byte(recordCount))
	// 获取proxies详细内容
	case http.MethodPost:
		var b Batch
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&b)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res, err := GetBatch(b.Cursor, b.Match, b.Count)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, err := utils.RestJson(res)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(data)
	// 减少分数
	case http.MethodDelete:
		var addr discover.Addr
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&addr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = decrease(addr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	// 更新分数为最大值
	case http.MethodPut:
		var addr discover.Addr
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&addr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = max(addr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

// redis controller
func redisHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodHead:
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		protocol := r.Form.Get("protocol")
		host := r.Form.Get("host")
		port, err := strconv.Atoi(r.Form.Get("port"))

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		isExits := exists(discover.Addr{
			Protocol: protocol,
			Host:     host,
			Port:     port,
		})
		if !isExits {
			w.WriteHeader(http.StatusNotFound)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// ping redis server
func pingHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		respData := pingDB()
		rw.Write([]byte(respData))
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func RegisterHandlers() {
	http.HandleFunc("/ping", pingHandler)
	// http.HandleFunc("/address", writeHandler)
	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/redis", redisHandler)
	http.HandleFunc("/tester", testHandler)
}
