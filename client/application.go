package client

import (
	"context"
	"log"

	"github.com/AntanasMaziliauskas/grpc/api"
	"google.golang.org/grpc"
)

type Application struct {
	conn *grpc.ClientConn
}

func (a *Application) Init() {
	var err error
	//Connecting to the server
	a.conn, err = grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

}
func (a *Application) Start() {
	c := api.NewGreetingClient(a.conn)
	response, err := c.SayHello(context.Background(), &api.Handshake{Id: "Node003", Port: "7778"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Timeout in: %d seconds", response.Timeout)
}

func (a *Application) Stop() {
	//closing connection to the server
	a.conn.Close()
}
