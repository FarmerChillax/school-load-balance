package main

import (
	"balance/log"
	"balance/registry"
	"balance/service"
	"balance/tester"
	"balance/utils"
	"context"
	"fmt"
	stlog "log"
)

func main() {
	host, port := utils.Config.Host, fmt.Sprintf("%d", utils.Config.Port)
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registration{
		ServiceName: registry.TesterService,
		ServiceURL:  serviceAddress,
		RequiredServices: []registry.ServiceName{
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
		tester.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatalln(err)
	}
	// 启动测试器
	go tester.Start()

	// 启动 log service
	if logProvider, err := registry.GetProvide(registry.LogService); err == nil {
		fmt.Printf("Logging service found at: %s\n", logProvider)
		log.SetClientLogger(logProvider, r.ServiceName)
	}

	<-ctx.Done()
	fmt.Println("Shutting down redis service.")
}
