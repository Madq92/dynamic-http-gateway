package main

import (
	"context"
	"dynamic-http-gateway/gateway"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//start business gateway
	b := gateway.NewBusiness()
	b.Start(":8080")

	//start admin api
	a := gateway.NewAdmin(b)
	a.Start(":8081")

	sigs := make(chan os.Signal)
	done := make(chan bool)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Println("stoping:", sig)

		ctx := context.Background()
		b.Shutdown(ctx)
		a.Shutdown(ctx)
		ctx.Done()
		close(done)
	}()
	<-done
}
