package main

import (
	"context"
	"distibuted/grades"
	"distibuted/log"
	"distibuted/registry"
	"distibuted/service"
	"fmt"
	stlog "log"
)

func main() {
	host, port := "localhost", "5001"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	r := registry.Registration{
		ServiceName:      registry.GradingService,
		ServiceURL:       serviceAddress,
		RequiredServices: []registry.ServiceName{registry.LogService},
		ServiceUpdateURl: serviceAddress + "/services ",
		HeartbeatURL:     serviceAddress + "/heartbeat",
	}

	ctx, err := service.Start(context.Background(), r, host, port, grades.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
		return
	}

	if logProvider, err := registry.GetProvoider(registry.LogService); err == nil {
		fmt.Printf("Logging service found at:%s\n", logProvider)
		log.SetClientLogger(logProvider, r.ServiceName)
	}
	<-ctx.Done()

	fmt.Println("Shutting down grading service")

}
