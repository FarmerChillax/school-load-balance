package discover

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type discoverConfig struct {
	Protocol, Host string
	Start, End     int
	Timeout        int
	SSL            bool
}

func RegisterHandlers() {
	http.HandleFunc("/start", DiscovererHandler)
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

func checkPort(start, end int) error {
	if start > end {
		return fmt.Errorf("end must big than start")
	}
	if start <= 0 || end > 65535 {
		return fmt.Errorf("scan port out of range")
	}
	return nil
}

func checkProtocol(protocol string) error {
	if len(protocol) == 0 {
		return nil
	}
	if protocol == "http" || protocol == "https" {
		return nil
	}
	return fmt.Errorf("protocol error")
}
