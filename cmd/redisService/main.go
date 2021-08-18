package main

import (
	"balance/log"
	"balance/registry"
	"balance/service"
	"balance/storage"
	"context"
	"encoding/json"
	"fmt"
	stlog "log"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

func (c *config) LoadConfig(path string) error {
	args := strings.Split(path, ".")
	if args[len(args)-1] != "json" {
		return fmt.Errorf("文件扩展名错误")
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

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)

	// host, port := "localhost", "7001"
	config, err := initConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	host, port := config.Host, fmt.Sprintf("%d", config.Port)

	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registration{
		ServiceName: registry.RedisService,
		ServiceURL:  serviceAddress,
		RequiredServices: []registry.ServiceName{
			registry.LogService,
		},
		ServiceUpdateURL: serviceAddress + "/services",
		HeartbeatURL:     serviceAddress + "/heartbeat",
	}

	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		storage.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatalln(err)
	}

	// 启动 log service
	if logProvider, err := registry.GetProvide(registry.LogService); err == nil {
		fmt.Printf("Logging service found at: %s\n", logProvider)
		log.SetClientLogger(logProvider, r.ServiceName)
	}

	<-ctx.Done()
	fmt.Println("Shutting down redis service.")

}
