package service

import (
	"balance/registry"
	"context"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, host, port string,
	reg registry.Registration,
	RegisterHandlersFunc func()) (context.Context, error) {

	RegisterHandlersFunc()
	ctx = startService(ctx, reg.ServiceName, host, port)
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func startService(ctx context.Context, ServiceName registry.ServiceName,
	host, port string) context.Context {

	ctx, cancel := context.WithCancel(ctx)
	var srv http.Server
	srv.Addr = ":" + port

	go func() {
		log.Println(srv.ListenAndServe())
		// 关闭的时候要取消注册
		// ...
		cancel()
	}()

	go func() {
		fmt.Printf("Service is running in http://%v\n", srv.Addr)
		fmt.Println("Registry service started. Press any key to stop.")
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}