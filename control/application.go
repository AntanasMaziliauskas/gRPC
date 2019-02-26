package control

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
	//var listen string

	source := "0.0.0.0:7778"
	a.conn, err = grpc.Dial(source, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
}

func (a *Application) GetOnePersonBroadcast() error {
	c := api.NewControlClient(a.conn)
	response, err := c.GetOnePersonBroadcast(context.Background(), &api.Person{Name: "Jonas"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}
func (a *Application) ListNodes() error {
	c := api.NewControlClient(a.conn)
	response, err := c.ListNodes(context.Background(), &api.Empty{})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}
