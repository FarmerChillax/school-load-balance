package main

import (
	"balance/log"
	"balance/registry"
	"balance/service"
	"balance/utils"
	"context"
	"fmt"
	stlog "log"
)

func main() {
	conf := utils.Config
	fileLogPath := conf.FileLogPath
	log.Run(fileLogPath)
	host, port := conf.Host, fmt.Sprintf("%d", conf.Port)
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registration{
		ServiceName:      registry.LogService,
		ServiceURL:       serviceAddress,
		RequiredServices: make([]registry.ServiceName, 0),
		ServiceUpdateURL: serviceAddress + "/services",
		HeartbeatURL:     serviceAddress + "/heartbeat",
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		log.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down log service.")
}
