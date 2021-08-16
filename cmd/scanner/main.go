package main

import (
	"balance/discover"
	"balance/registry"
	"balance/service"
	"context"
	"fmt"
	stlog "log"
)

func main() {
	host, port := "localhost", "8000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registration{
		ServiceName: registry.ScanService,
		ServiceURL:  serviceAddress,
		RequiredServices: []registry.ServiceName{
			registry.RedisService,
			registry.LogService,
		},
		ServiceUpdateURL: serviceAddress + "/services",
	}

	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		discover.RegisterHandlers,
	)

	if err != nil {
		stlog.Fatalln(err)
	}

	<-ctx.Done()
	fmt.Println("Shutting down scanner service.")

}
