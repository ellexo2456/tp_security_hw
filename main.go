package main

import (
	"github.com/ellexo2456/tp_security_hw/proxy"
	"log"
	"net/http"
)

func main() {
	p := proxy.Server{}

	server := &http.Server{
		Addr:    ":8080",
		Handler: p,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf(err.Error())
	}
}
