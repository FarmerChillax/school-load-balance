package storage

import (
	"balance/network"
	"encoding/json"
	"log"
	"net/http"
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
		err = writeDB(addrs)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

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
	http.HandleFunc("/write", writeHandler)
}
