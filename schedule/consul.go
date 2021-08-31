package schedule

import (
	"balance/registry"
	"balance/storage"
	"balance/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Consul Config
type ConsulConfig struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Path     string `json:"path"`
	// ServiceName string          `json:"service_name"`
	Services []ConsulService `json:"services"`
}

type ConsulService struct {
	Key   string      `json:"Key"`
	Value NginxConfig `json:"Value"`
}

type NginxConfig struct {
	Weight      int `json:"weight"`
	MaxFails    int `json:"max_fails"`
	FailTimeout int `json:"fail_timeout"`
}

func (c ConsulConfig) SendServices() error {
	// protocol + host + port + path
	records, err := getRecords()
	if err != nil {
		return err
	}
	c.Services = setup(records)
	consulAddres := fmt.Sprintf("%s://%s:%d%s", c.Protocol, c.Host, c.Port, c.Path)
	for _, item := range c.Services {
		// url + name
		addres := fmt.Sprintf("%s/%s", consulAddres, item.Key)
		value, err := utils.RestJson(item.Value)
		if err != nil {
			return err
		}
		buf := bytes.NewBuffer(value)
		req, err := http.NewRequest(http.MethodPut, addres, buf)
		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to registry consul. Registry service responded with code %v. URL: %v", res.StatusCode, addres)
		}
	}
	return nil
}

// 设置发送消息内容
func setup(rs storage.Records) (newSerives []ConsulService) {
	for _, record := range rs {
		nc := NginxConfig{}
		nc = nc.New()
		nc.Weight = record.Score
		addr := fmt.Sprintf("%s:%d", record.Member.Host, record.Member.Port)
		consulservice := ConsulService{
			Key:   addr,
			Value: nc,
		}
		newSerives = append(newSerives, consulservice)
	}
	return newSerives
}

// 构造结构体
func (nc NginxConfig) New() NginxConfig {
	nc.Weight = 1
	nc.FailTimeout = 10
	nc.MaxFails = 2
	return nc
}

var Consul *ConsulConfig

func init() {
	Consul.Protocol = utils.Config.Consul.Protocol
	Consul.Host = utils.Config.Consul.Host
	Consul.Port = utils.Config.Consul.Port
	Consul.Path = utils.Config.Consul.Path
}

// 将数据库内容发送到consul
func getRecords() (records storage.Records, err error) {
	// 获取数据库内容
	redisURL, err := registry.GetProvide(registry.RedisService)
	if err != nil {
		return records, err
	}
	resp, err := http.Get(redisURL + "/redis")
	if err != nil {
		return records, err
	}
	if resp.StatusCode != http.StatusOK {
		return records, fmt.Errorf("fail to get values, status code: %v", resp.StatusCode)
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&records)
	if err != nil {
		return records, err
	}
	return records, nil
}
