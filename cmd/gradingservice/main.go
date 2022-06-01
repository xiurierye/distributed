package main

import (
	"context"
	"distibuted/grades"
	"distibuted/registry"
	"distibuted/service"
	"fmt"
	stlog "log"
)

func main() {
	host, port := "localhost", "5001"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	r := registry.Registration{
		ServiceName: registry.GradingService,
		ServiceURL:  serviceAddress,
	}

	ctx, err := service.Start(context.Background(), r, host, port, grades.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
		return
	}
	<-ctx.Done()

	fmt.Println("Shutting down grading service")

}
