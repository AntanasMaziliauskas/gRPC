package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/AntanasMaziliauskas/grpc/api"
	"google.golang.org/grpc"
)

//Node structure holds values of Port, LastSeen and Connection
type Node struct {
	Port       string
	LastSeen   time.Time
	Connection *grpc.ClientConn
}

//SetServer function creates listener, server object
func (a *Application) SetServer() error {
	var err error

	a.lis, err = net.Listen("tcp", Source)

	a.grpcServer = grpc.NewServer()

	api.RegisterNodeServer(a.grpcServer, a.Broker)
	api.RegisterControlServer(a.grpcServer, a.Broker)

	return err
}

//StartServer function starts the gRPC server
func (a *Application) StartServer() {
	fmt.Println("Server Starting")
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := a.grpcServer.Serve(a.lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

}
