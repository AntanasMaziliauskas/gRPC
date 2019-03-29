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

//BrokerService holds Init, Start, Stop, AddNode, DropNode, ListNodes, ListPersonsBroadcast
//ListPersonsNode, GetOnePersonBroadcast, GetOnePersonNode, GetMultiPersonBroadcast,
//GetMultiPersonNode, DropOnepersonBroadcast, DropOnePersonNode, DropMultiPersonBroadcast,
//DropMultiPersonNode, UpsertOnePersonNode and UpsertMultiPersonNode functions
type BrokerService interface {
	Init() error

	Start(int)

	Stop() error

	AddNode(context.Context, *api.NodeInfo) (*api.Empty, error)

	DropNode(context.Context, *api.NodeInfo) (*api.Empty, error)

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

	UpsertOnePersonNode(context.Context, *api.Person) (*api.Empty, error)
	UpsertMultiPersonNode(context.Context, *api.MultiPerson) (*api.Empty, error)
}

//GRPCBroker structure holds values of Nodes, ctx, cancel, wg and GrpcServer
type GRPCBroker struct {
	Nodes      map[string]*Node
	ctx        context.Context
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
	GrpcServer *grpc.Server
}

//Node structure holds values of Port, LastSeen, LastPing, IsOnline and Connection
type Node struct {
	Port       string
	LastSeen   time.Time
	LastPing   time.Time //TODO: Ar tikrai jos reikia?
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
					if err != nil && v.IsOnline {
						v.IsOnline = false
						fmt.Printf("Node %s went offline\n", k)
					} else if err == nil && !v.IsOnline {
						g.Nodes[k].IsOnline = true
						fmt.Printf("Node %s went online\n", k)
					}
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
		if err := v.Connection.Close(); err != nil {
			log.Println("Error while trying to close connection with Node:", err)
		}
	}
	g.cancel()
	g.wg.Wait()
	return nil
}

//AddNode function adds connected Node to a list
func (g *GRPCBroker) AddNode(ctx context.Context, in *api.NodeInfo) (*api.Empty, error) {
	log.Printf("Node %s connected.", in.Id)

	conn, err := grpc.Dial(in.Source, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Unable to connect to %s: %s", in.Id, err)
	}
	//TODO: Kas vyksta, kai iskarto nepavyksta connection uzmegsti?
	g.Nodes[in.Id] = &Node{
		Port:       in.Source,
		LastSeen:   time.Now(),
		IsOnline:   true,
		Connection: conn,
	}

	return &api.Empty{}, nil
}

//DropNode deletes Node from the list
func (g *GRPCBroker) DropNode(ctx context.Context, in *api.NodeInfo) (*api.Empty, error) {

	if _, ok := g.Nodes[in.Id]; ok {
		delete(g.Nodes, in.Id)
		fmt.Printf("Node %s was dropped.", in.Id)

		return &api.Empty{}, nil
	}
	fmt.Println("Unable to drop Node: " + in.Id + ". It is not in the list.")

	return &api.Empty{}, nil
}

//ListNodes function return a list of Nodes that are connected to the server
func (g *GRPCBroker) ListNodes(ctx context.Context, in *api.Empty) (*api.NodesList, error) {
	nodesList := &api.NodesList{}

	for k, v := range g.Nodes {
		nodesList.Nodes = append(nodesList.Nodes, &api.NodeInfo{Id: k, Source: v.Port, Isonline: v.IsOnline})
	}
	if len(nodesList.Nodes) < 1 {
		fmt.Println("There are no nodes connected to the server")
	}

	return nodesList, nil
}

//ListPersonsBroadcast gets a list of all the persons in every connected Node and return information about him
func (g *GRPCBroker) ListPersonsBroadcast(ctx context.Context, in *api.Empty) (*api.MultiPerson, error) {
	var (
		response *api.MultiPerson
		err      error
		con      bool
	)
	var wg sync.WaitGroup //
	wg.Add(len(g.Nodes))  //
	response1 := &api.MultiPerson{}
	c := make(chan api.MultiPerson) //

	for k, v := range g.Nodes {
		//Startas
		startTime := time.Now()
		if !v.IsOnline {

			continue
		}

		go func(k string, v *Node, c chan api.MultiPerson) {
			defer wg.Done()
			a := api.NewServerClient(v.Connection)
			if response, err = a.ListPersons(context.Background(), &api.Empty{}); err != nil {
				g.Nodes[k].IsOnline = false
				log.Println("Error while trying to call ListPersons: ", err)

				return
			}
			con = true
			//response1.Persons = append(response1.Persons, response.Persons...)
			g.Nodes[k].LastSeen = time.Now()
			//Finishas
			duration := time.Now().Sub(startTime)
			log.Println("Atsakas iš Node užtruko: ", duration)

			c <- *response
		}(k, v, c)

	}
	go func() {
		wg.Wait()
		close(c)
	}()
	for v := range c {
		//fmt.Println(v)
		response1.Persons = append(response1.Persons, v.Persons...)
	}

	if !con {
		log.Println("There are no Nodes Online.")
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
		//Startas
		startTime := time.Now()
		if k == in.Id && v.IsOnline {
			a := api.NewServerClient(v.Connection)
			if response, err = a.ListPersons(context.Background(), &api.Empty{}); err != nil {
				g.Nodes[k].IsOnline = false
				log.Println("Error while trying to call ListPersons for Node: ", k, ". Error:", err)

				return &api.MultiPerson{}, nil
			}
			g.Nodes[k].LastSeen = time.Now()

			duration := time.Now().Sub(startTime)
			log.Println("Atsakas iš Node užtruko: ", duration)

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
		//Startas
		startTime := time.Now()
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
		if len(response.Id) > 1 {
			duration := time.Now().Sub(startTime)
			log.Println(duration)

			return response, nil
		}
		duration := time.Now().Sub(startTime)
		log.Println("Atsakas iš Node užtruko: ", duration)
	}
	if !connect {
		log.Println("There are no Nodes connected.")
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
		//Startas
		startTime := time.Now()
		if k == in.Node && v.IsOnline {
			a := api.NewServerClient(v.Connection)
			if response, err = a.GetOnePerson(context.Background(), in); err != nil {
				g.Nodes[k].IsOnline = false
				log.Println("Error while trying to run GetOnePerson: ", err)

				continue
			}
			g.Nodes[k].LastSeen = time.Now()
			//finish
			duration := time.Now().Sub(startTime)
			log.Println("Atsakas iš Node užtruko: ", duration)
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
	}

	return response1, nil
}

//GetMultiPersonNode function get information about multiple persons from a specific Node
func (g *GRPCBroker) GetMultiPersonNode(ctx context.Context, in *api.MultiPerson) (*api.MultiPerson, error) {
	var (
		response *api.MultiPerson
		err      error
	)

	for k, v := range g.Nodes {
		if k == in.Persons[0].Node && v.IsOnline {
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
	for k, v := range g.Nodes {
		if k == in.Node && v.IsOnline {
			a := api.NewServerClient(v.Connection)
			if _, err := a.DropOnePerson(context.Background(), in); err != nil {
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
		if response, err = a.DropMultiPerson(context.Background(), in); err != nil {
			g.Nodes[k].IsOnline = false
			log.Println("Error while trying to run DropMultiPerson: ", err)
		}
		g.Nodes[k].LastSeen = time.Now()
		success = true
	}
	if !success {
		log.Println("There are no Nodes connected.")
	}

	return response, nil

}

//DropMultiPersonNode deleted person from specified Node
func (g *GRPCBroker) DropMultiPersonNode(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {
	var (
		response *api.Empty
		err      error
	)

	for k, v := range g.Nodes {
		if k == in.Persons[0].Node && v.IsOnline {
			a := api.NewServerClient(v.Connection)
			if response, err = a.DropMultiPerson(context.Background(), in); err != nil {
				g.Nodes[k].IsOnline = false
				log.Println("Error while trying to run DropMultiPerson: ", err)

				return &api.Empty{}, nil
			}
			g.Nodes[k].LastSeen = time.Now()

			return response, nil
		}
	}
	log.Println("Given Node is not connected to the server")

	return &api.Empty{}, nil
}

//UpsertOnePersonNode adds given Person to a specified Node
func (g *GRPCBroker) UpsertOnePersonNode(ctx context.Context, in *api.Person) (*api.Empty, error) {
	var err error

	for k, v := range g.Nodes {
		if k == in.Node && v.IsOnline {
			a := api.NewServerClient(v.Connection)
			if _, err = a.UpsertOnePerson(context.Background(), in); err != nil {
				g.Nodes[k].IsOnline = false
				log.Println("Error while trying to run UpsertOnePerson: ", err)

				return &api.Empty{}, nil
			}
			g.Nodes[k].LastSeen = time.Now()

			return &api.Empty{}, nil
		}
	}
	log.Println("Given Node is not connected to the server")

	return &api.Empty{}, nil
}

//UpsertMultiPersonNode adds given multiple persons to a specified Node
func (g *GRPCBroker) UpsertMultiPersonNode(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {
	for k, v := range g.Nodes {
		if k == in.Persons[0].Node && v.IsOnline {
			a := api.NewServerClient(v.Connection)
			if _, err := a.UpsertMultiPerson(context.Background(), in); err != nil {
				g.Nodes[k].IsOnline = false
				log.Println("Error while trying to run UpsertMultiPerson: ", err)

				return &api.Empty{}, nil
			}
			g.Nodes[k].LastSeen = time.Now()

			return &api.Empty{}, nil

		}
	}
	log.Println("Given Node is not connected to the server")

	return &api.Empty{}, nil
}
