package discover

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	Protocol     = "http"
	HostTemplate = "172.16.%d.%d"
)

type discoverConfig struct {
	Protocol, Host string
	Start, End     int
	Timeout        int
	SSL            bool
}

func RegisterHandlers() {
	http.HandleFunc("/start", DiscovererHandler)
	http.HandleFunc("/randome", randHandler)
}

func randHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		go randDisvocer()
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func randDisvocer() {
	randNumber := MakeRandSegment()
	for i := 1; i <= 255; i++ {
		host := fmt.Sprintf(HostTemplate, randNumber, i)
		go commonPortDiscover(host)
	}
}

func DiscovererHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// 默认扫描
		go commonDiscover()
		w.Write([]byte("start success."))
	case http.MethodPost:
		// 扫描指定段
		var dConfig discoverConfig
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&dConfig)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = Discover(dConfig)
		if err != nil {
			log.Fatalln(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodHead:
		go fastDiscover()
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func fastDiscover() {
	for i := 1; i <= 255; i++ {
		host := fmt.Sprintf(HostTemplate, 1, i)
		go commonPortDiscover(host)
	}
}

func commonDiscover() {
	for cSegment := 1; cSegment <= 255; cSegment++ {
		for dSegment := 1; dSegment <= 255; dSegment++ {
			host := fmt.Sprintf(HostTemplate, cSegment, dSegment)
			go commonPortDiscover(host)
		}
		time.Sleep(time.Second * 2)
	}
}

func Discover(discoverConfig discoverConfig) error {
	err := checkProtocol(discoverConfig.Protocol)
	if err != nil {
		return err
	}
	err = checkPort(discoverConfig.Start, discoverConfig.End)
	if err != nil {
		return err
	}
	timeout := time.Second * time.Duration(discoverConfig.Timeout)
	go discoverer(discoverConfig.Protocol, discoverConfig.Host,
		discoverConfig.Start, discoverConfig.End,
		timeout, discoverConfig.SSL)

	return nil
}
