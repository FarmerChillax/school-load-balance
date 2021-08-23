package main

import (
	"balance/registry"
	"balance/schedule"
	"balance/service"
	"context"
	"fmt"
)

func main() {
	host, port := "localhost", "8848"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registration{
		ServiceName: registry.ScheduleService,
		ServiceURL:  serviceAddress,
		RequiredServices: []registry.ServiceName{
			registry.LogService,
			registry.RedisService,
		},
		ServiceUpdateURL: serviceAddress + "/services",
		HeartbeatURL:     serviceAddress + "/heartbeat",
	}

	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		schedule.RegisterHandlers,
	)
	if err != nil {
		fmt.Println(err)
	}

	<-ctx.Done()
	fmt.Println("Shutting down redis service.")
}
