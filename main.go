package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		httpAddr = flag.String("http", ":8080", "http listen address")
	)
	flag.Parse()
	context.Background()
	srv := NewService()
	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	endpoints := Endpoints{
		GetEndpoint:      MakeGetEndpoint(srv),
		StatusEndpoint:   MakeStatusEndpoint(srv),
		ValidateEndpoint: MakeValidateEndpoint(srv),
	}

	go func() {
		log.Println("http is listening on port:", *httpAddr)
		handler := NewHTTPServer(endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	log.Fatalln(<-errChan)
}
