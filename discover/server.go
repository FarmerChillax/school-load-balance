package discover

import (
	"fmt"
	"net/http"
	"time"
)

func RegisterHandlers() {
	http.HandleFunc("/start", DiscovererHandler)
}

func DiscovererHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// 默认扫描
		// fastDiscover()
		commonDiscover()
	case http.MethodPost:
		// 扫描指定段
		w.Header().Add("test", "ok")
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func fastDiscover() {
	hostTemplate := "192.168.2.%d"
	for i := 1; i <= 255; i++ {
		host := fmt.Sprintf(hostTemplate, i)
		go commonPortDiscover(host)
	}
}

func commonDiscover() {
	hostTemplate := "192.168.%d.%d"
	for cSegment := 1; cSegment <= 255; cSegment++ {
		for dSegment := 1; dSegment <= 255; dSegment++ {
			host := fmt.Sprintf(hostTemplate, cSegment, dSegment)
			go commonPortDiscover(host)
		}
		time.Sleep(time.Second * 2)
	}
}
