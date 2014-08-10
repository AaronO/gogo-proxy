package main

import (
	"github.com/AaronO/gogo-proxy/proxy"
	"net/http"
)

func main() {
	p, _ := proxy.New(proxy.ProxyOptions{
		Balancer: func(req *http.Request) (string, error) {
			return "https://www.google.com"+req.URL.Path, nil
		},
	})

    http.ListenAndServe(":8080", p)
}
