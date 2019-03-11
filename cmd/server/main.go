package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/AntanasMaziliauskas/grpc/server"
	"github.com/AntanasMaziliauskas/grpc/server/broker"
)

// main start a gRPC server and waits for connection
func main() {

	app := server.Application{Broker: &broker.GRPCBroker{}}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM,/* syscall.SIGSTOP,*/ syscall.SIGKILL)

	app.Init()

	app.Start()

	<-stop

	app.Stop()
}
