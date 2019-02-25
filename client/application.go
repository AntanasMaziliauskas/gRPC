package client

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/AntanasMaziliauskas/grpc/api"
	"github.com/AntanasMaziliauskas/grpc/client/person"
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

func (a *Application) Init() {

	a.Person.Init()

	a.ConnectToServer()

	a.SettingServer()

}
func (a *Application) Start() {

	a.GreetingWithServer()

	a.PingServer()

	a.StartServer()

}

func (a *Application) Stop() {
	//closing connection to the server
	a.conn.Close()
}

//Pasiruosimas connectionui su serveriu |Init Susijungimas su serveriu
func (a *Application) ConnectToServer() {
	var err error
	//Connecting to the server
	a.conn, err = grpc.Dial(fmt.Sprintf(":%s", a.ServerPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
}

//Connectionas su serveriu | Start Pasisveikinimo issiuntimas serveriui
func (a *Application) GreetingWithServer() {
	a.Person.Init()
	c := api.NewGreetingClient(a.conn)
	response, err := c.SayHello(context.Background(), &api.Node{Id: a.ID, Port: a.ServerPort})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Timeout in: %d seconds", response.Timeout)
	a.Timeout = response.Timeout
}

//Pasiruosimas jungti serveri |INIT
func (a *Application) SettingServer() {
	var err error
	var s Application
	a.lis, err = net.Listen("tcp", fmt.Sprintf(":%s", a.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a gRPC server object
	a.grpcServer = grpc.NewServer()
	// attach the Greeting service to the server
	api.RegisterLookForDataServer(a.grpcServer, &s)
}

//Listeneris serverio |START
func (a *Application) StartServer() {
	if err := a.grpcServer.Serve(a.lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

//Pinginimo siuntimas serveriui | Start GO rutina
func (a *Application) PingServer() {
	p := api.NewPingClient(a.conn)
	go func() {
		ticker := time.NewTicker(time.Duration(a.Timeout) / 2 * time.Second)
		for {
			select {
			case <-ticker.C:
				log.Printf("Pinging")
				_, err := p.PingMe(context.Background(), &api.PingMessage{Id: a.ID})
				if err != nil {
					log.Fatalf("Error when calling PingMe: %s", err)
				}
				//log.Printf("Response from server: %d", response)
				//wg.Done()
			}
		}
	}()
}

//Ateinancios uzklausos priimimas ir vykdymas
func (a *Application) FindData(ctx context.Context, in *api.LookFor) (*api.Found, error) {
	name := in.Name
	//TODO: Kodel negaliu naudoti butent cia?
	//a.Person.Init()
	//	found, _ := a.Person.GetOne(name)
	//return &api.Found{Name: found.Name, Age: found.Age, Profession: found.Profession}, nil
	return &api.Found{Name: name}, nil
}
