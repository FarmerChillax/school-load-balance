package main

import (
	"fmt"
	"net/http"
)

func main() {
	HelloAPI := HelloAPI{}
	server := http.Server{
		Addr: "localhost:8080",
	}
	http.Handle("/api", &HelloAPI)

	fmt.Printf("API Serve is running in http://%v\n", server.Addr)
	server.ListenAndServe()
}

type HelloAPI struct{}

func (h *HelloAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API")
}