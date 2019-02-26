package server

import (
	"net"

	"github.com/AntanasMaziliauskas/grpc/server/broker"
	"google.golang.org/grpc"
)

const Timeout = 15

type Application struct {
	Broker     broker.BrokerService
	grpcServer *grpc.Server
	lis        net.Listener
	//TODO
	//Ar hardcodinam timeout? Ar kaip flaga?
}

//Init function get the server and the broker ready
func (a *Application) Init() {
	a.SetServer() // Error 	log.Fatalf("failed to listen: %v", err)
	a.Broker.Init()
}

//
func (a *Application) Start() {
	a.StartServer()
	a.Broker.Start(Timeout)

	_ = a.StartHTTPServer() // Error log.Fatal("ListenAndServe: ", err)

}

func (a *Application) Stop() {
	//uzdaryti serveri
}
