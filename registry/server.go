package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const ServerPort = ":6660"
const ServicesURL = "http://localhost" + ServerPort + "/services"

type registry struct {
	registrations []Registration
	mutex         *sync.Mutex
}

// 全局变量
var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.Mutex),
}

// 添加注册
func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	return nil
}

// 移除注册（注销）
func (r *registry) remove(url string) error {
	for index := range reg.registrations {
		if reg.registrations[index].ServiceURL == url {
			r.mutex.Lock()
			tmp := reg.registrations[index]
			reg.registrations = append(reg.registrations[:index], reg.registrations[index+1:]...)
			r.mutex.Unlock()
			log.Printf("Logout serveice: %v success.", tmp.ServiceName)
			return nil
		}
	}
	return fmt.Errorf("Service at URL %s not found", url)
}

// 注册组件的 web服务
type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request received")
	switch r.Method {
	// 注册服务
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var r Registration
		err := decoder.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v with URL: %s\n", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)
		log.Printf("Removing service at URL: %s", url)
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
