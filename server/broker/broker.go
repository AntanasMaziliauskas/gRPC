package broker

import (
	"context"
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

	ListNodes(context.Context, *api.Empty) (*api.NodeInfo, error)

	GetOnePersonBroadcast(context.Context, *api.Person) (*api.Person, error)
	//GetOnePersonNode(ctx context.Context, node, person string) (*api.Person, error)
	/*
		GetMultiPersonBroadcast(ctx context.Context, person string) ([]api.Person, error)
		GetMultiPersonNode(ctx context.Context, node, person string) ([]api.Person, error)

		DropOnePersonBroadcast(ctx context.Context, person string) error
		DropOnePersonNode(ctx context.Context, node, person string) error

		DropMultiPersonBroadcast(ctx context.Context, person []string) error
		DropMultiPersonNode(ctx context.Context, node, person []string) error

		InsertOnePersonNode(ctx context.Context, node string, person api.Person) error
		InsertMultiPersonNode(ctx context.Context, node string, person []api.Person) error*/
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

//DropNode deleted Node from the list
func (g *GRPCBroker) DropNode(ctx context.Context, in *api.NodeInfo) (*api.Empty, error) {
	delete(g.Nodes, in.Id)

	return &api.Empty{}, nil
}

//ReceivedPing function updates LastSeen for the Node that pinged the server
func (g *GRPCBroker) Ping(ctx context.Context, in *api.PingMessage) (*api.Empty, error) {
	fmt.Println("I Got Pinged From ", in.Id)
	g.Nodes[in.Id].LastSeen = time.Now()

	return &api.Empty{}, nil
}

//ListNodes function return a list of Nodes that are connected to the server
func (g *GRPCBroker) ListNodes(ctx context.Context, in *api.Empty) (*api.NodesList, error) {
	var nodesList *api.NodesList
	for k, v := range g.Nodes {
		nodesList = append(nodesList, &api.NodesList{
			Id:     k,
			Source: v.Port,
		})
	}
	return &g.Nodes, nil
}

//GetOnePersonBroadcast function go through the list of connected Nodes
//requests to look for data with the name given, retrieved the data and return it.
func (g *GRPCBroker) GetOnePersonBroadcast(ctx context.Context, in *api.Person) (*api.Person, error) {
	//var response *api.Person
	var err error

	for _, v := range g.Nodes {
		a := api.NewServerClient(v.Connection)
		_, err = a.GetOnePersonBroadcast(context.Background(), &api.Person{Name: in.Name})
		if err != nil {
			//log.Fatalf("Error when trying to get response from server: %s", err)
		}
	}

	return &api.Person{Name: "Petras"}, err
}

/*
func (g *GRPCBroker) GetOnePersonNode(node string, name string) (*api.Person, error) {

	a := api.NewLookForDataClient(g.Nodes[node].Connection)

	return a.FindData(context.Background(), &api.LookFor{Name: name})
}*/
