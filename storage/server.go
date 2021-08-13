package storage

import (
	"balance/network"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func writeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var addrs network.Addrs
		err := decoder.Decode(&addrs)
		log.Printf("recive data: %v\n", len(addrs))
		if err != nil {
			log.Println(err)
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

// utils
func utilsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		recordCount := fmt.Sprintf("%d", count())
		w.Write([]byte(recordCount))
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
		isExits := exists(network.Addr{
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
	http.HandleFunc("/utils", utilsHandler)
}
