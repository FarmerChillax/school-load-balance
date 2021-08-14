package tester

import (
	"balance/storage"
	"fmt"
	"net/http"
)

func RegisterHandlers() {
	http.HandleFunc("/start", startHandler)
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// start tester
		batch := storage.Batch{
			Cursor: 0,
			Match:  "",
			Count:  10,
		}
		err := getAddrs(batch)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
