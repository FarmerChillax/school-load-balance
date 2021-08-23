package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Redis struct {
	Server        string `json:"server"`
	Port          int    `json:"port"`
	Password      string `json:"password"`
	DB            int    `json:"db"`
	Key           string `json:"key"`
	SCORE_MAX     int    `json:"SCORE_MAX"`
	SCORE_MIN     int    `json:"SCORE_MIN"`
	SCORE_DEFAULT int    `json:"SCORE_DEFAULT"`
}

type Consul struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Path     string `json:"path"`
}

type config struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Protocol     string `json:"protocol"`
	RegistryAddr string `json:"registryAddr"`
	RegistryPort int    `json:"registryPort"`
	FileLogPath  string `json:"fileLogPath"`
	Redis        Redis  `json:"redis"`
	Consul       Consul `json:"consul"`
}

func (c *config) LoadConfig(path string) error {
	args := strings.Split(path, ".")
	if args[len(args)-1] != "json" {
		return fmt.Errorf("faild verification file extension")
	}
	file, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		return err
	}
	return nil
}

func initConfig() (c config, err error) {
	path := os.Args[len(os.Args)-1]
	err = c.LoadConfig(path)
	if err != nil {
		return c, err
	}
	return c, nil
}

var Config config

func init() {
	var err error
	Config, err = initConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
}
