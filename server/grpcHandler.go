package server

import (
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

//SetServer function creates listener, server object
func (a *Application) SetServer() error {
	var err error

	a.lis, err = net.Listen("tcp", fmt.Sprintf(":%d", 7778))

	a.grpcServer = grpc.NewServer()

	api.RegisterNodeServer(a.grpcServer, a.Broker)
	api.RegisterControlServer(a.grpcServer, a.Broker)
	//	api.RegisterPingServer(a.grpcServer, a)

	return err
}

//StartServer function starts the gRPC server
func (a *Application) StartServer() {
	fmt.Println("Server Starting")
	go func() {
		if err := a.grpcServer.Serve(a.lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

}
