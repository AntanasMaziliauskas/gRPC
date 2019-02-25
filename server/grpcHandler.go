package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/AntanasMaziliauskas/grpc/api"
	"google.golang.org/grpc"
)

type Node struct {
	Port       string
	LastSeen   time.Time
	Connection *grpc.ClientConn
}

//Pasiruosimas jungti serveri |INIT
func (a *Application) SetServer() {
	var err error
	var s Application // kažkas čia negerai
	//a.timeout = 30
	//Creating a listener
	a.lis, err = net.Listen("tcp", fmt.Sprintf(":%d", 7778))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a gRPC server object
	a.grpcServer = grpc.NewServer()
	// attach the Greeting service to the server
	api.RegisterGreetingServer(a.grpcServer, &s)
	// attach the Ping service to the server
	api.RegisterPingServer(a.grpcServer, &s)
}

//Listeneris serverio |START
func (a *Application) StartServer() {
	fmt.Println("Server Starting")
	go func() {
		if err := a.grpcServer.Serve(a.lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

}

//Pasisveikinimo priimimas is Node kliento
func (a *Application) SayHello(ctx context.Context, in *api.Node) (*api.Timeout, error) {

	log.Printf("Received message: %s. Port: %s", in.Id, in.Port)

	a.addNode(in)

	return &api.Timeout{Timeout: Timeout}, nil
}

//Pingo priimimas is Node kliento
func (a *Application) PingMe(ctx context.Context, in *api.PingMessage) (*api.Empty, error) {
	fmt.Println("I Got Pinged From ", in.Id)
	a.receivedPing(in)

	return &api.Empty{}, nil
}

//TURI KELIAUTI I BROKERI
func (a *Application) addNode(in *api.Node) {

	conn, err := grpc.Dial(fmt.Sprintf(":%s", in.Port), grpc.WithInsecure()) // Portas ateina is NODE
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	//fmt.Println("test")

	Nodes[in.Id] = &Node{
		Port:       in.Port,
		LastSeen:   time.Now(),
		Connection: conn,
	}
	fmt.Println(Nodes[in.Id])
}

func (a *Application) receivedPing(in *api.PingMessage) {
	Nodes[in.Id].LastSeen = time.Now()

	//	fmt.Println(Nodes[in.Id])
}
