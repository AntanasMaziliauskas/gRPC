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
	//timeout int64
}

//TODO Kazkaip protingiau turbut galima apsirasyti
var Nodes map[string]*Node

func (a *Application) Init() {
	Nodes = make(map[string]*Node)
	a.SetServer()
	a.Broker.Init()
}

func (a *Application) Start() {
	a.StartServer()
	//a.Broker.Start()
	a.StartHTTPServer()

}

func (a *Application) Stop() {
	//uzdaryti serveri
}
