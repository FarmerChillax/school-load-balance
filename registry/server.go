package registry

import (
	"bytes"
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
	mutex         *sync.RWMutex
}

// 全局变量
var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.RWMutex),
}

// 添加注册
func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	err := r.sendRequiredServices(reg)
	return err
}

// 注册中心向服务发送依赖相关内容
// e.g 依赖的url，名字等
func (r registry) sendRequiredServices(reg Registration) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	var p patch
	// 遍历已注册的服务
	for _, serviceReg := range r.registrations {
		// 遍历需要的服务
		for _, reqService := range reg.RequiredServices {
			// 找到依赖项，添加到patch
			if reqService == serviceReg.ServiceName {
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceURL,
				})
			}
		}
	}
	err := r.sendPatch(p, reg.ServiceUpdateURL)
	if err != nil {
		return err
	}
	return nil
}

// 发送依赖项
func (r registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return err
	}
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
	return fmt.Errorf("service at URL %s not found", url)
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
