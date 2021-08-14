package main

import (
	"balance/registry"
	"balance/service"
	"balance/tester"
	"context"
	"fmt"
	"log"
)

func main() {
	host, port := "localhost", "3500"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registration{
		ServiceName: registry.TesterService,
		ServiceURL:  serviceAddress,
		RequiredServices: []registry.ServiceName{
			registry.RedisService,
		},
		ServiceUpdateURL: serviceAddress + "/services",
	}

	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		tester.RegisterHandlers,
	)
	if err != nil {
		log.Fatalln(err)
	}

	<-ctx.Done()
	fmt.Println("Shutting down redis service.")
	// tester.GetAddrs(storage.Batch{
	// 	Cursor: 0,
	// 	Match:  "",
	// 	Count:  10,
	// })
}
