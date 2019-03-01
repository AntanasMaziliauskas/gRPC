package server

import (
	"log"
	"net"
	"sync"

	"github.com/AntanasMaziliauskas/grpc/server/broker"
	"google.golang.org/grpc"
)

const Timeout = 15

//Source is a server address
const Source = ":7778"

type Application struct {
	Broker     broker.BrokerService
	grpcServer *grpc.Server
	lis        net.Listener
	wg         *sync.WaitGroup
}

//Init function get the server and the broker ready
func (a *Application) Init() {
	a.wg = &sync.WaitGroup{}

	if err := a.SetServer(); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	if err := a.Broker.Init(); err != nil {
		log.Fatalf("Failed to launch Broker: %v", err)
	}
}

//Start function starts the server and broker services also start HTTP server
func (a *Application) Start() {
	a.StartServer()
	a.Broker.Start(Timeout)

	/*	if err := a.StartHTTPServer(); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}*/

}

//Stop function stops gRPC server, Broker services.
func (a *Application) Stop() {
	a.grpcServer.Stop()
	a.Broker.Stop()
	a.wg.Wait()
}
