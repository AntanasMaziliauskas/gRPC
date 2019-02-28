package broker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/AntanasMaziliauskas/grpc/api"
	"google.golang.org/grpc"
)

type BrokerService interface {
	Init() error

	Start(int)

	Stop() error

	AddNode(context.Context, *api.NodeInfo) (*api.Timeout, error)

	DropNode(context.Context, *api.NodeInfo) (*api.Empty, error)

	Ping(context.Context, *api.PingMessage) (*api.Empty, error)

	ListNodes(context.Context, *api.Empty) (*api.NodesList, error)

	ListPersonsBroadcast(context.Context, *api.Empty) (*api.MultiPerson, error)
	ListPersonsNode(context.Context, *api.NodeInfo) (*api.MultiPerson, error)

	GetOnePersonBroadcast(context.Context, *api.Person) (*api.Person, error)
	GetOnePersonNode(context.Context, *api.Person) (*api.Person, error)

	GetMultiPersonBroadcast(context.Context, *api.MultiPerson) (*api.MultiPerson, error)
	GetMultiPersonNode(context.Context, *api.MultiPerson) (*api.MultiPerson, error)

	DropOnePersonBroadcast(context.Context, *api.Person) (*api.Empty, error)
	DropOnePersonNode(context.Context, *api.Person) (*api.Empty, error)

	DropMultiPersonBroadcast(context.Context, *api.MultiPerson) (*api.Empty, error)
	DropMultiPersonNode(context.Context, *api.MultiPerson) (*api.Empty, error)

	InsertOnePersonNode(context.Context, *api.Person) (*api.Empty, error)
	InsertMultiPersonNode(context.Context, *api.MultiPerson) (*api.Empty, error)
}

type GRPCBroker struct {
	Nodes      map[string]*Node
	ctx        context.Context
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
	GrpcServer *grpc.Server
}

type Node struct {
	Port       string
	LastSeen   time.Time
	Connection *grpc.ClientConn
}

//Init function makes Nodes map, assigns WaitGroup and Context values
func (g *GRPCBroker) Init() error {
	g.Nodes = make(map[string]*Node)
	g.wg = &sync.WaitGroup{}
	g.ctx, g.cancel = context.WithCancel(context.Background())

	//g.GrpcServer = grpc.NewServer()

	//api.RegisterNodeServer(g.GrpcServer, g)
	return nil
}

//Start function runs go routine and checks for connections that are timed out
//deleted nodes that are timed out from the Nodes list
func (g *GRPCBroker) Start(timeout int) {
	fmt.Println("Starting TimeOut service")

	go func() {
		g.wg.Add(1)
		ticker := time.NewTicker(time.Duration(timeout) * time.Second)
		for {
			select {
			case <-ticker.C:
				for k, v := range g.Nodes {
					diff := time.Now().Sub(v.LastSeen)
					sec := int(diff.Seconds())
					if sec > timeout {
						delete(g.Nodes, k)
						fmt.Println(k, "Node has timed out and was deleted from the list.")
					}
				}
			case <-g.ctx.Done():
				log.Println("TimeOut service has stopped.")
				g.wg.Done()
			}
		}
	}()
}

//Stop function stops Go Routine
func (g *GRPCBroker) Stop() error {
	g.cancel()
	g.wg.Wait()

	return nil
}

//AddNode function adds connected Node to a list
func (g *GRPCBroker) AddNode(ctx context.Context, in *api.NodeInfo) (*api.Timeout, error) {
	log.Printf("Received message: %s. Port: %s", in.Id, in.Source)
	conn, err := grpc.Dial(in.Source, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	g.Nodes[in.Id] = &Node{
		Port:       in.Source,
		LastSeen:   time.Now(),
		Connection: conn,
	}

	return &api.Timeout{Timeout: 15}, err
}

//DropNode deletes Node from the list
func (g *GRPCBroker) DropNode(ctx context.Context, in *api.NodeInfo) (*api.Empty, error) {

	for k := range g.Nodes {
		if k == in.Id {
			delete(g.Nodes, in.Id)

			return &api.Empty{Response: "Node successfully deleted"}, nil
		}
	}
	err := errors.New("Unable to drop Node: " + in.Id + ". It is not connected to the server")
	return &api.Empty{}, err
}

//ReceivedPing function updates LastSeen for the Node that pinged the server
func (g *GRPCBroker) Ping(ctx context.Context, in *api.PingMessage) (*api.Empty, error) {
	//fmt.Println("I Got Pinged From ", in.Id)

	//TODO: Ka darome, jeigu Node yra dropintas ir visvien pingina?
	for k := range g.Nodes {
		if k == in.Id {
			fmt.Println("I Got Pinged From ", in.Id)
			g.Nodes[in.Id].LastSeen = time.Now()
		}
	}

	return &api.Empty{}, nil
}

//ListNodes function return a list of Nodes that are connected to the server
func (g *GRPCBroker) ListNodes(ctx context.Context, in *api.Empty) (*api.NodesList, error) {
	var err error

	nodesList := &api.NodesList{}

	for k, v := range g.Nodes {
		nodesList.Nodes = append(nodesList.Nodes, &api.NodeInfo{Id: k, Source: v.Port})
	}
	if len(nodesList.Nodes) < 1 {
		err = errors.New("There are no nodes connected to the server")
	}
	return nodesList, err
}

func (g *GRPCBroker) ListPersonsBroadcast(ctx context.Context, in *api.Empty) (*api.MultiPerson, error) {
	var (
		response *api.MultiPerson
		err      error
	)
	response1 := &api.MultiPerson{}
	for _, v := range g.Nodes {
		a := api.NewServerClient(v.Connection)
		if response, err = a.ListPersons(context.Background(), &api.Empty{}); err == nil {
			for _, r := range response.Persons {
				response1.Persons = append(response1.Persons, r)
			}
		}
	}

	return response1, err
}

func (g *GRPCBroker) ListPersonsNode(ctx context.Context, in *api.NodeInfo) (*api.MultiPerson, error) {
	var (
		response *api.MultiPerson
		err      error
	)

	for k, v := range g.Nodes {
		if k == in.Id {
			a := api.NewServerClient(v.Connection)
			if response, err = a.ListPersons(context.Background(), &api.Empty{}); err == nil {
				return response, err
			}

			return response, err
		}
	}
	err = errors.New("Given Node is not connected to the server")

	return response, err
}

//GetOnePersonBroadcast function go through the list of connected Nodes
//requests to look for data with the name given, retrieved the data and return it.
func (g *GRPCBroker) GetOnePersonBroadcast(ctx context.Context, in *api.Person) (*api.Person, error) {
	var (
		response *api.Person
		err      error
	)

	for _, v := range g.Nodes {
		a := api.NewServerClient(v.Connection)
		if response, err = a.GetOnePerson(context.Background(), &api.Person{Name: in.Name}); err == nil {
			return response, err
		}
	}

	return response, err
}

func (g *GRPCBroker) GetOnePersonNode(ctx context.Context, in *api.Person) (*api.Person, error) {
	var (
		response *api.Person
		err      error
	)

	for k, v := range g.Nodes {
		if k == in.Node {
			a := api.NewServerClient(v.Connection)
			response, err = a.GetOnePerson(context.Background(), &api.Person{Name: in.Name})

			return response, err
		}
	}
	err = errors.New("Given Node is not connected to the server")

	return response, err
}

func (g *GRPCBroker) GetMultiPersonBroadcast(ctx context.Context, in *api.MultiPerson) (*api.MultiPerson, error) {
	var (
		response *api.MultiPerson
		err      error
	)

	response1 := &api.MultiPerson{}
	for _, v := range g.Nodes {
		a := api.NewServerClient(v.Connection)
		if response, err = a.GetMultiPerson(context.Background(), in); err == nil {
			//TODO: Kaip gudriau sudeti i slice?
			for _, r := range response.Persons {
				response1.Persons = append(response1.Persons, r)
			}
		}
	}
	if len(response1.Persons) < 1 {
		err = errors.New("Unable to locate given persons")
		return response1, err
	}
	return response, err
}

func (g *GRPCBroker) GetMultiPersonNode(ctx context.Context, in *api.MultiPerson) (*api.MultiPerson, error) {
	var (
		response *api.MultiPerson
		err      error
	)

	for k, v := range g.Nodes {
		if k == in.Persons[1].Node {
			a := api.NewServerClient(v.Connection)
			response, err = a.GetMultiPerson(context.Background(), in)
			return response, err
		}
	}
	err = errors.New("Given Node is not connected to the server")

	return response, err
}

func (g *GRPCBroker) DropOnePersonBroadcast(ctx context.Context, in *api.Person) (*api.Empty, error) {
	var (
		response *api.Empty
		err      error
	)

	for _, v := range g.Nodes {
		a := api.NewServerClient(v.Connection)
		response, err = a.DropOnePerson(context.Background(), &api.Person{Name: in.Name})
		if err != nil {
			//log.Fatalf("Error when trying to get response from server: %s", err)
		}
	}

	return response, err
}

func (g *GRPCBroker) DropOnePersonNode(ctx context.Context, in *api.Person) (*api.Empty, error) {
	var (
		response *api.Empty
		err      error
	)

	for k, v := range g.Nodes {
		if k == in.Node {
			a := api.NewServerClient(v.Connection)
			response, err = a.DropOnePerson(context.Background(), &api.Person{Name: in.Name})
			if err != nil {
				//log.Fatalf("Error when trying to get response from server: %s", err)
			}
		}
	}

	return response, err
}

func (g *GRPCBroker) DropMultiPersonBroadcast(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {
	var (
		response *api.Empty
		err      error
	)

	for _, v := range g.Nodes {
		a := api.NewServerClient(v.Connection)
		response, err = a.DropMultiPerson(context.Background(), in)
		if err != nil {
			//log.Fatalf("Error when trying to get response from server: %s", err)
		}
	}

	return response, err
}

func (g *GRPCBroker) DropMultiPersonNode(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {
	var (
		response *api.Empty
		err      error
	)

	for k, v := range g.Nodes {
		if k == in.Persons[1].Node {
			a := api.NewServerClient(v.Connection)
			response, err = a.DropMultiPerson(context.Background(), in)
			if err != nil {
				//log.Fatalf("Error when trying to get response from server: %s", err)
			}
		}
	}

	return response, err
}

func (g *GRPCBroker) InsertOnePersonNode(ctx context.Context, in *api.Person) (*api.Empty, error) {
	var (
		response *api.Empty
		err      error
	)

	for k, v := range g.Nodes {
		if k == in.Node {
			a := api.NewServerClient(v.Connection)
			response, err = a.InsertOnePerson(context.Background(), in)
			if err != nil {
				//log.Fatalf("Error when trying to get response from server: %s", err)
			}
		}
	}

	return response, err
}

func (g *GRPCBroker) InsertMultiPersonNode(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {
	var (
		response *api.Empty
		err      error
	)

	for k, v := range g.Nodes {
		if k == in.Persons[1].Node {
			a := api.NewServerClient(v.Connection)
			response, err = a.InsertMultiPerson(context.Background(), in)
			if err != nil {
				//log.Fatalf("Error when trying to get response from server: %s", err)
			}
		}
	}

	return response, err
}
