package service

import (
	"context"
	"distibuted/registry"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, r registry.Registration, host, port string, registerHandler func()) (context.Context, error) {
	registerHandler()
	ctx = startService(ctx, r.ServiceName, host, port)

	err := registry.RegisterService(r)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func startService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context {

	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server

	srv.Addr = ":" + port
	go func() {
		log.Println(srv.ListenAndServe())

		cancel()
	}()

	go func() {
		fmt.Printf("%v staredc, Press any key to stop. \n", serviceName)
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()

	}()

	return ctx
}
