package main

import (
	"balance/log"
	"balance/registry"
	"balance/service"
	"context"
	"fmt"
	stlog "log"
)

func main() {
	fileLogPath := "./distributed.log"
	log.Run(fileLogPath)
	host, port := "localhost", "6500"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registration{
		ServiceName:      registry.LogService,
		ServiceURL:       serviceAddress,
		RequiredServices: make([]registry.ServiceName, 0),
		ServiceUpdateURL: serviceAddress + "/services",
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
