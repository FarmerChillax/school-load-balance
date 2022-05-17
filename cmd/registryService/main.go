package main

import (
	"balance/registry"
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.MemProfile, profile.MemProfileRate(1)).Stop()
	go func() {
		http.ListenAndServe(":5122", nil)
	}()
	registry.SetupRegistryService()
	http.Handle("/services", &registry.RegistryService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server
	srv.Addr = registry.ServerPort

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Printf("Service is running in %v\n", registry.ServicesURL)
		fmt.Println("Registry service started. Press any key to stop.")
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()
	}()

	<-ctx.Done()
	fmt.Println("Shutting down registry service.")
}
