package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

// func ShutdownService(url string) error {
// 	// req, err := http.NewRequest(http.MethodDelete, )
// }

func RegisterService(r Registration) error {
	serviceUpdateURL, err := url.Parse(r.ServiceUpdateURL)
	if err != nil {
		return err
	}
	// 依赖接口绑定处理函数
	http.Handle(serviceUpdateURL.Path, &serviceUpdateHandler{})

	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	err = encoder.Encode(r)
	if err != nil {
		return err
	}
	res, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service. Registry service "+
			"responded with code %v", res.StatusCode)
	}
	return nil
}

// 传入停止服务的url，向注册中心发送注销请求
// url格式为-> ip:port
func ShutdownService(url string) error {
	req, err := http.NewRequest(http.MethodDelete, ServicesURL,
		// request请求体内容->url
		bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to deregister service. Registry service responded with code %v", res.StatusCode)
	}
	return nil
}

// /////// 更新依赖部分 ///////

type serviceUpdateHandler struct{}

func (suh serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dec := json.NewDecoder(r.Body)
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("Updated received %v\n", p)
	prov.Update(p)
}

// 存放服务需要的依赖
type providers struct {
	// 服务名->服务url
	services map[ServiceName][]string
	mutex    *sync.RWMutex
}

// 更新依赖
func (p *providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, patchEntry := range pat.Added {
		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = make([]string, 0)
		}
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL)
	}

	for _, patchEntry := range pat.Removed {
		if providerURLs, ok := p.services[patchEntry.Name]; ok {
			for index := range providerURLs {
				if providerURLs[index] == patchEntry.URL {
					p.services[patchEntry.Name] = append(
						providerURLs[:index], providerURLs[index+1:]...)
				}
			}
		}
	}

}

// 通过服务名称，找到依赖的urls，从依赖项里面随机返回一个url
// (后期改urls->返回服务的全部url)
// 私有函数
func (p providers) get(name ServiceName) (string, error) {
	providers, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("no providers available for service %v", name)
	}
	idx := int(rand.Float32() * float32(len(providers)))
	return providers[idx], nil
}

func GetProvide(name ServiceName) (string, error) {
	return prov.get(name)
}

var prov = providers{
	// 一个服务有多个组件
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}
