package node

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/AntanasMaziliauskas/grpc/api"
	"github.com/AntanasMaziliauskas/grpc/node/person"
	"github.com/phayes/freeport"
	"google.golang.org/grpc"
)

type Application struct {
	conn       *grpc.ClientConn
	Port       string
	ID         string
	ServerPort string
	lis        net.Listener
	grpcServer *grpc.Server
	Timeout    int64
	Path       string
	Person     person.PersonService
}

//Init function runs Person.Init, connects to server and ?????
func (a *Application) Init() {

	a.Person.Init() //Error ReadFile error

	a.ConnectToServer() //Error log.Fatalf("did not connect: %s", err)

	a.SettinggRPCServer()

}
func (a *Application) Start() {

	a.GreetingWithServer()

	a.PingServer()

	a.StartgRPCServer()

}

func (a *Application) Stop() {
	//closing connection to the server
	a.conn.Close()
}

//ConnectToServer function connects to server
func (a *Application) ConnectToServer() error {
	var err error

	a.conn, err = grpc.Dial(a.ServerPort, grpc.WithInsecure())

	return err
}

//Connectionas su serveriu | Start Pasisveikinimo issiuntimas serveriui // pervadinti i sayhello
func (a *Application) GreetingWithServer() {
	c := api.NewNodeClient(a.conn)
	response, err := c.AddNode(context.Background(), &api.NodeInfo{Id: a.ID, Source: a.Port})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Timeout in: %d seconds", response.Timeout)
	a.Timeout = response.Timeout
}

//?????
func (a *Application) SettinggRPCServer() error {
	var err error
	//Sugeneruoja random porta
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}
	a.Port = fmt.Sprintf(":%d", port)

	a.lis, err = net.Listen("tcp", a.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a gRPC server object
	a.grpcServer = grpc.NewServer()
	// attach the Greeting service to the server
	api.RegisterServerServer(a.grpcServer, a.Person)

	return err
}

//Listeneris serverio |START
func (a *Application) StartgRPCServer() {
	if err := a.grpcServer.Serve(a.lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

//Pinginimo siuntimas serveriui | Start GO rutina
func (a *Application) PingServer() {
	p := api.NewNodeClient(a.conn)
	go func() {
		ticker := time.NewTicker(time.Duration(a.Timeout) / 2 * time.Second)
		for {
			select {
			case <-ticker.C:
				log.Printf("Pinging")
				_, err := p.Ping(context.Background(), &api.PingMessage{Id: a.ID})
				if err != nil {
					log.Fatalf("Error when calling PingMe: %s", err)
				}
				//log.Printf("Response from server: %d", response)
				//wg.Done()
			}
		}
	}()
}

/*
//Ateinancios uzklausos priimimas ir vykdymas
func (a *Application) FindData(ctx context.Context, in *api.LookFor) (*api.Person, error) {
	name := in.Name

	person, _ := a.Person.GetOne(name)

	return &api.Person{Name: person.Name, Age: person.Age, Profession: person.Profession}, nil
}
*/
