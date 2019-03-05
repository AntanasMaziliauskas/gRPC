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

	AddNode(context.Context, *api.NodeInfo) (*api.Empty, error)

	DropNode(context.Context, *api.NodeInfo) (*api.Empty, error)

	//	Ping(context.Context, *api.PingMessage) (*api.Empty, error)

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
	LastPing   time.Time
	IsOnline   bool
	Connection *grpc.ClientConn
}

//Init function makes Nodes map, assigns WaitGroup and Context values
func (g *GRPCBroker) Init() error {
	g.Nodes = make(map[string]*Node)
	g.wg = &sync.WaitGroup{}
	g.ctx, g.cancel = context.WithCancel(context.Background())

	return nil
}

//Start function runs go routine and checks for connections that are timed out
//deleted nodes that are timed out from the Nodes list
func (g *GRPCBroker) Start(timeout int) {
	fmt.Println("Starting Ping service")

	g.wg.Add(1)
	go func() {
		ticker := time.NewTicker(time.Duration(timeout) * time.Second)
		for {
			select {
			case <-ticker.C:

				for k, v := range g.Nodes {
					a := api.NewServerClient(v.Connection)
					_, err := a.Ping(context.Background(), &api.PingMessage{})
					if err != nil && v.IsOnline == true {
						v.IsOnline = false
						fmt.Printf("Node %s went offline\n", k)
					} /*else if !v.IsOnline {
						g.Nodes[k].IsOnline = true
						fmt.Printf("Node %s went online\n", k)
					}*/
					g.Nodes[k].LastSeen = time.Now()
				}
			case <-g.ctx.Done():

				log.Println("Ping service has stopped.")
				g.wg.Done()

				return

			}
		}
	}()
}

//Stop function stops Go Routine, closes connections with Nodes
func (g *GRPCBroker) Stop() error {
	for _, v := range g.Nodes {
		v.Connection.Close()
	}
	g.cancel()
	g.wg.Wait()
	return nil
}

//AddNode function adds connected Node to a list
func (g *GRPCBroker) AddNode(ctx context.Context, in *api.NodeInfo) (*api.Empty, error) {
	log.Printf("Received message: %s. Port: %s", in.Id, in.Source)
	conn, err := grpc.Dial(in.Source, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	g.Nodes[in.Id] = &Node{
		Port:       in.Source,
		LastSeen:   time.Now(),
		IsOnline:   true,
		Connection: conn,
	}

	return &api.Empty{}, err
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

//ListNodes function return a list of Nodes that are connected to the server
func (g *GRPCBroker) ListNodes(ctx context.Context, in *api.Empty) (*api.NodesList, error) {
	var err error

	nodesList := &api.NodesList{}

	for k, v := range g.Nodes {
		nodesList.Nodes = append(nodesList.Nodes, &api.NodeInfo{Id: k, Source: v.Port, Isonline: v.IsOnline})
	}
	if len(nodesList.Nodes) < 1 {
		err = errors.New("There are no nodes connected to the server")
	}
	return nodesList, err
}

//ListPersonsBroadcast gets a list of all the persons in every connected Node and return information about him
func (g *GRPCBroker) ListPersonsBroadcast(ctx context.Context, in *api.Empty) (*api.MultiPerson, error) {
	var (
		response *api.MultiPerson
		err      error
		con      bool
	)

	response1 := &api.MultiPerson{}

	for k, v := range g.Nodes {
		if !v.IsOnline {
			continue
		}
		a := api.NewServerClient(v.Connection)
		if response, err = a.ListPersons(context.Background(), &api.Empty{}); err != nil {
			g.Nodes[k].IsOnline = false
			log.Println("Error while trying to call ListPersons: ", err)
			continue
		}
		con = true
		response1.Persons = append(response1.Persons, response.Persons...)
		g.Nodes[k].LastSeen = time.Now()
	}
	if !con {
		log.Println("There are no Online Nodes.")
	}
	return response1, nil
}

//ListPersonsNode gets a list of persons in a specific node and return information
func (g *GRPCBroker) ListPersonsNode(ctx context.Context, in *api.NodeInfo) (*api.MultiPerson, error) {
	var (
		response *api.MultiPerson
		err      error
	)

	for k, v := range g.Nodes {
		if k == in.Id && v.IsOnline == true {
			a := api.NewServerClient(v.Connection)
			if response, err = a.ListPersons(context.Background(), &api.Empty{}); err != nil {
				g.Nodes[k].IsOnline = false
				log.Println("Error while trying to call ListPersons for Node: ", k, ". Error:", err)

				return &api.MultiPerson{}, nil
			}
			g.Nodes[k].LastSeen = time.Now()

			return response, nil
		}
	}
	log.Println("Node ", in.Id, " is not connected.")

	return &api.MultiPerson{}, nil
}

//GetOnePersonBroadcast function go through the list of connected Nodes
//requests to look for data with the name given, retrieved the data and return it.
func (g *GRPCBroker) GetOnePersonBroadcast(ctx context.Context, in *api.Person) (*api.Person, error) {
	var (
		response *api.Person
		err      error
		connect  bool
	)

	for k, v := range g.Nodes {
		if !v.IsOnline {
			continue
		}
		a := api.NewServerClient(v.Connection)
		if response, err = a.GetOnePerson(context.Background(), in); err != nil {
			g.Nodes[k].IsOnline = false
			log.Println("Error while trying to call GetOnePerson for Node: ", k, ". Error:", err)
			continue
		}
		g.Nodes[k].LastSeen = time.Now()
		connect = true
	}
	if !connect {
		log.Println("There are no Nodes connected.")

		return response, nil
	}

	return response, nil
}

//GetOnePersonNode function gets information about a specific person from specified Node
//Returns received information
func (g *GRPCBroker) GetOnePersonNode(ctx context.Context, in *api.Person) (*api.Person, error) {
	var (
		response *api.Person
		err      error
	)

	for k, v := range g.Nodes {
		if k == in.Node && v.IsOnline == true {
			a := api.NewServerClient(v.Connection)
			if response, err = a.GetOnePerson(context.Background(), in); err != nil {
				g.Nodes[k].IsOnline = false
				log.Println("Error while trying to run GetOnePerson: ", err)
				continue
			}
			g.Nodes[k].LastSeen = time.Now()

			return response, nil
		}
	}
	log.Println("Given Node is not connected.")

	return &api.Person{}, nil
}

//GetMultiPersonBroadcast function gets information about multiple persons
//Looks through every Node that is Online
func (g *GRPCBroker) GetMultiPersonBroadcast(ctx context.Context, in *api.MultiPerson) (*api.MultiPerson, error) {
	var (
		response *api.MultiPerson
		err      error
		connect  bool
	)

	response1 := &api.MultiPerson{}
	for k, v := range g.Nodes {
		if !v.IsOnline {
			continue
		}
		a := api.NewServerClient(v.Connection)
		if response, err = a.GetMultiPerson(context.Background(), in); err != nil {
			g.Nodes[k].IsOnline = false
			log.Println("Error while trying to run GetMultiPerson: ", err)

			continue
		}
		g.Nodes[k].LastSeen = time.Now()
		response1.Persons = append(response1.Persons, response.Persons...)
		connect = true
	}
	if !connect {
		log.Println("There are no Nodes connected.")

		return &api.MultiPerson{}, nil
	}

	return response, nil
}

//GetMultiPersonNode function get information about multiple persons from a specific Node
func (g *GRPCBroker) GetMultiPersonNode(ctx context.Context, in *api.MultiPerson) (*api.MultiPerson, error) {
	var (
		response *api.MultiPerson
		err      error
	)

	for k, v := range g.Nodes {
		if k == in.Persons[0].Node && v.IsOnline == true {
			a := api.NewServerClient(v.Connection)
			if response, err = a.GetMultiPerson(context.Background(), in); err != nil {
				g.Nodes[k].IsOnline = false
				log.Println("Error while trying to run GetMultiPerson: ", err)

				return &api.MultiPerson{}, nil
			}
			g.Nodes[k].LastSeen = time.Now()

			return response, nil
		}
	}
	log.Println("Given Node is not connected.")

	return &api.MultiPerson{}, nil
}

//DropOnePersonBroadcast looks through all connected Nodes and deletes specified person
func (g *GRPCBroker) DropOnePersonBroadcast(ctx context.Context, in *api.Person) (*api.Empty, error) {
	var (
		//response *api.Empty
		err     error
		success bool
	)

	for k, v := range g.Nodes {
		if !v.IsOnline {
			continue
		}
		a := api.NewServerClient(v.Connection)
		if _, err = a.DropOnePerson(context.Background(), in); err != nil {
			g.Nodes[k].IsOnline = false
			log.Println("Error while trying to run DropOnePerson: ", err)

			continue
		}
		success = true
		g.Nodes[k].LastSeen = time.Now()
	}

	if !success {
		log.Println("No Nodes connected.")
	}

	return &api.Empty{}, nil
}

//DropOnePersonNode deletes given person from specified Node
func (g *GRPCBroker) DropOnePersonNode(ctx context.Context, in *api.Person) (*api.Empty, error) {
	var (
		//response *api.Empty
		err error
	)

	for k, v := range g.Nodes {
		if k == in.Node && v.IsOnline == true {
			a := api.NewServerClient(v.Connection)
			if _, err = a.DropOnePerson(context.Background(), in); err != nil {
				g.Nodes[k].IsOnline = false
				log.Println("Error while trying to run DropOnePerson: ", err)

				return &api.Empty{}, nil
			}
			g.Nodes[k].LastSeen = time.Now()

			return &api.Empty{}, nil
		}

	}

	log.Println("Given Node is not connected")

	return &api.Empty{}, nil
}

//DropMultiPersonBroadcast deletes multiple persons going through all nodes connected.
func (g *GRPCBroker) DropMultiPersonBroadcast(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {
	var (
		response *api.Empty
		err      error
		success  bool
	)

	for k, v := range g.Nodes {
		if !v.IsOnline {
			continue
		}
		a := api.NewServerClient(v.Connection)
		response, err = a.DropMultiPerson(context.Background(), in)
		if err == nil {
			g.Nodes[k].LastSeen = time.Now()
			success = true
		}
	}
	if success {

		return &api.Empty{Response: "Successfully dropped"}, nil
	}

	return response, err

}

//DropMultiPersonNode deleted person from specified Node
func (g *GRPCBroker) DropMultiPersonNode(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {
	var (
		response *api.Empty
		err      error
	)

	for k, v := range g.Nodes {
		if k == in.Persons[0].Node && v.IsOnline == true {
			a := api.NewServerClient(v.Connection)
			if response, err = a.DropMultiPerson(context.Background(), in); err == nil {
				g.Nodes[k].LastSeen = time.Now()
				return response, err
			} else {
				return response, err
			}
		}
	}
	err = errors.New("Given Node is not connected to the server")

	return &api.Empty{}, err
}

//InsertOnePersonNode adds given Person to a specified Node
func (g *GRPCBroker) InsertOnePersonNode(ctx context.Context, in *api.Person) (*api.Empty, error) {
	var (
		response *api.Empty
		err      error
	)

	for k, v := range g.Nodes {
		if k == in.Node && v.IsOnline == true {
			a := api.NewServerClient(v.Connection)
			if response, err = a.InsertOnePerson(context.Background(), in); err == nil {
				g.Nodes[k].LastSeen = time.Now()
				return response, err
			} else {
				return response, err
			}
		}
	}
	err = errors.New("Given Node is not connected to the server")

	return response, err
}

//InsertMultiPersonNode adds given multiple persons to a specified Node
func (g *GRPCBroker) InsertMultiPersonNode(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {
	var (
		response *api.Empty
		err      error
	)

	for k, v := range g.Nodes {
		if k == in.Persons[0].Node && v.IsOnline == true {
			a := api.NewServerClient(v.Connection)
			if response, err = a.InsertMultiPerson(context.Background(), in); err == nil {
				g.Nodes[k].LastSeen = time.Now()
				return response, err
			} else {
				return response, err
			}
		}
	}
	err = errors.New("Given Node is not connected to the server")

	return response, err
}
