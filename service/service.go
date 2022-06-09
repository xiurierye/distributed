package service

import (
	"context"
	"distibuted/registry"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, r registry.Registration, host, port string, registerHandler func()) (context.Context, error) {
	registerHandler() //registerHandler 用于自身服务的 url 绑定,
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
		log.Println(srv.ListenAndServe()) //启动web 服务

		err := registry.ShutdowService(fmt.Sprintf("http://%s:%s", host, port)) //启动失败,向注册中心发送 关闭服务信号
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()

	go func() {
		fmt.Printf("%v started, Press any key to stop. \n", serviceName) // 暂停自身服务方式
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)

		err := registry.ShutdowService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		cancel()

	}()

	return ctx
}
