package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/AntanasMaziliauskas/grpc/api"
	"github.com/AntanasMaziliauskas/grpc/server/broker"
	"google.golang.org/grpc"
)

type Application struct {
	Broker broker.Broker
	//http       //
	grpcServer *grpc.Server
	lis        net.Listener
	//TODO
	//Ar hardcodinam timeout? Ar kaip flaga?
	timeout int64
}

func (a *Application) Init() {
	var err error
	var s Application
	a.timeout = 30
	//Creating a listener
	a.lis, err = net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a gRPC server object
	a.grpcServer = grpc.NewServer()
	// attach the Greeting service to the server
	api.RegisterGreetingServer(a.grpcServer, &s)
	// attach the Ping service to the server
	//api.RegisterPingServer(a.grpcServer, &s)

}

func (a *Application) Start() {

	go func() {
		if err := a.grpcServer.Serve(a.lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

}

func (a *Application) SayHello(ctx context.Context, in *api.Handshake) (*api.Timeout, error) {

	log.Printf("Received message: %s. Port: %s", in.Id, in.Port)

	/*if !sliceContainsString(in.Id, s.Nodes) {
		nodeStruct := Node{
			ID:             in.Id,
			LastConnection: time.Now(),
		}
		s.Nodes = append(s.Nodes, nodeStruct)
	}*/
	//TODO Tikrinam error?
	_ = a.Broker.AddNode(in)

	return &api.Timeout{Timeout: a.timeout}, nil
}

func (a *Application) HTTPHandleGet() {
	//	p, err := a.Broker.GetOnePersonBroadcast(...)
	//w.Write(json(p))
}

//a := Application{Broker: &broker.GRPCBroker{}}
