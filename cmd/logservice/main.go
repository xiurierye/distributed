package main

import (
	"context"
	"distibuted/log"
	"distibuted/registry"
	"distibuted/service"
	"fmt"
	stlog "log"
)

func main() {

	log.Run("./distributed.log")
	host, port := "localhost", "4000"

	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	var r = registry.Registration{
		ServiceName: registry.LogService,
		ServiceURL:  serviceAddress,
	}

	ctx, err := service.Start(context.Background(), r, host, port, log.RegisterHandlers)

	if err != nil {
		stlog.Fatalln(err)
	}

	<-ctx.Done()

	fmt.Print("Shuting down log service.")
}
