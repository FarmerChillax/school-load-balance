package schedule

import (
	"fmt"
	"net/http"
)

// type Scheduls struct {
// 	SchedulName string
// }

type SchedulHandler struct{}

func (sh SchedulHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// w.WriteHeader(http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := Consul.SendServices()
	// fmt.Println(Consul.Services)
	if err != nil {
		fmt.Println(err)
	}
}

func RegisterHandlers() {
	sh := SchedulHandler{}
	http.Handle("/test", &sh)
}
