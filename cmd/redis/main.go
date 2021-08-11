package main

import (
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

	<-ctx.Done()
	fmt.Println("Shutting down redis service.")

}
