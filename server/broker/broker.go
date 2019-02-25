package broker

import (
	"fmt"
	"time"

	"github.com/AntanasMaziliauskas/grpc/api"
	"google.golang.org/grpc"
)

type BrokerService interface {
	Init() error
	// Startas suks cikla per visus nodes ir ieskos pas ka last_seen yra mazas
	// ir atliks DropNode jei yra timeout
	Start() error
	//Sustabdo start cikla
	Stop() error
	// api.Node | Receivinu info is Node ir uzmezgu connection.
	AddNode(*api.Node) error
	//Istrinu Node is saraso
	//DropNode(node string) error
	//Updatinu Last Seen, kuomet gaunu ping
	ReceivedPing() // paupdeitini last seen
	//Returninu Node sarasa. Kam sito reikia? HHTP, kad atvaizduoti gal
	/*	ListNodes() []broker.Node

		GetOnePersonBroadcast(ctx context.Context, person string) (api.Person, error)
		GetOnePersonNode(ctx context.Context, node, person string) (api.Person, error)

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
	Nodes map[string]Node
}

type Node struct {
	Port       string
	LastSeen   time.Time
	Connection *grpc.ClientConn
}

func (g *GRPCBroker) Init() error { return nil }

func (g *GRPCBroker) Start() error {
	fmt.Println("Start")

	return nil
}

func (g *GRPCBroker) Stop() error { return nil }

func (g *GRPCBroker) AddNode(in *api.Node) error {

	fmt.Println(in)
	//Connection padarom su Node iskarto
	//conn, err := grpc.Dial(":7778", grpc.WithInsecure()) // Portas ateina is NODE
	//	if err != nil {
	//	log.Fatalf("did not connect: %s", err)
	//	}
	//Nereikia uzdaryti connection
	//defer conn.Close()

	//a := api.NewLookForDataClient(conn)

	//TODO: Neaiski vieta. Ar tikrai galim connections deti i map? Ar dedam dar kazka?
	//Sudedam connectionus i map'a
	//g.Nodes[in.Id] = Node{
	//	Name:       in.Id,
	//		Port:       in.Port,
	//		LastSeen:   time.Now(),
	//		Connection: conn,
	//NewLookForDataClient: a (?)
	//	}
	return nil
}

//func (g *GRPCBroker) DropNode(node string) error { return nil }

func (g *GRPCBroker) ReceivedPing() {
	fmt.Println("Received Ping")
	//return nil
}

/*

func (g *GRPCBroker) ListNodes() error { return nil }

func (g *GRPCBroker) GetOnePersonBroadcast(...) {
	conn = g.nodes[node_id].Connection.

	a := api.NewLookForDataClient(conn)

	return a.GetNode() // kvieciam pas node kazka kas moka grazint person
}


func (g *GRPCBroker) GetOnePersonNode(...) {
	conn = g.nodes[node_id].Connection.

	a := api.NewLookForDataClient(conn)

	return a.GetNode() // kvieciam pas node kazka kas moka grazint person
}*/
