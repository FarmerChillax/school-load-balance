package tester

import (
	"net/http"
)

func RegisterHandlers() {
	http.HandleFunc("/start", startHandler)
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		go testService()
		msg := "tester start!"
		w.Write([]byte(msg))
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
