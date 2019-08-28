package main

import (
	"context"
	"crypto/tls"
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

	certificate, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err!= nil {

	}
	certificate.

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
