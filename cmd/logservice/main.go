package main

import (
	"balance/log"
	"balance/registry"
	"fmt"
)

func main() {
	fileLogPath := "./distributed.log"
	log.Run(fileLogPath)
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registration{
		ServiceName: registry.LogService,
		ServiceURL:  serviceAddress,
	}
	// ctx, err := service
}
