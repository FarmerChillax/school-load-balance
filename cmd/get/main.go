package main

import (
	"fmt"
	"net/http"
)

func main() {
	resp, err := http.Get("http://farmer233.top:80")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
}
