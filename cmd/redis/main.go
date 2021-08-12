package main

import (
	"balance/log"
	"balance/registry"
	"balance/service"
	"balance/storage"
	"context"
	"fmt"
	stlog "log"
)

func main() {
	host, port := "localhost", "7000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registration{
		ServiceName: registry.RedisService,
		ServiceURL:  serviceAddress,
		RequiredServices: []registry.ServiceName{
			registry.LogService,
		},
		ServiceUpdateURL: serviceAddress + "/services",
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
