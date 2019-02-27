package control

import (
	"context"
	"log"
	"strings"

	"github.com/AntanasMaziliauskas/grpc/api"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

type Application struct {
	conn *grpc.ClientConn
}

func (a *Application) Init() {
	var err error

	source := "0.0.0.0:7778"
	a.conn, err = grpc.Dial(source, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
}

func (a *Application) GetOnePersonBroadcast(c *cli.Context) error {
	b := api.NewControlClient(a.conn)

	response, err := b.GetOnePersonBroadcast(context.Background(), &api.Person{Name: c.GlobalString("person")})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

func (a *Application) DropNode(c *cli.Context) error {
	b := api.NewControlClient(a.conn)
	response, err := b.DropNode(context.Background(), &api.NodeInfo{Id: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

func (a *Application) ListNodes(c *cli.Context) error {
	b := api.NewControlClient(a.conn)
	response, err := b.ListNodes(context.Background(), &api.Empty{})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}
func (a *Application) GetOnePersonNode(c *cli.Context) error {
	b := api.NewControlClient(a.conn)
	response, err := b.GetOnePersonNode(context.Background(), &api.Person{Name: c.GlobalString("person"), Node: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

func (a *Application) GetMultiPersonBroadcast(c *cli.Context) error {
	b := api.NewControlClient(a.conn)
	multiPerson := &api.MultiPerson{}

	personList := a.parsePersons(c.GlobalString("person"))
	for _, v := range personList {
		multiPerson.Persons = append(multiPerson.Persons, &api.Person{Name: v})
	}
	//fmt.Println(multiPerson)
	//multiPerson.Persons = []*api.Person{Name: personList}
	response, err := b.GetMultiPersonBroadcast(context.Background(), multiPerson)
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}
func (a *Application) GetMultiPersonNode(c *cli.Context) error {
	b := api.NewControlClient(a.conn)
	multiPerson := &api.MultiPerson{}

	personList := a.parsePersons(c.GlobalString("person"))
	for _, v := range personList {
		multiPerson.Persons = append(multiPerson.Persons, &api.Person{Name: v, Node: c.GlobalString("node")})
	}
	//fmt.Println(multiPerson)
	//multiPerson.Persons = []*api.Person{Name: personList}
	response, err := b.GetMultiPersonNode(context.Background(), multiPerson)
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}
func (a *Application) DropOnePersonBroadcast(c *cli.Context) error {
	b := api.NewControlClient(a.conn)

	response, err := b.DropOnePersonBroadcast(context.Background(), &api.Person{Name: c.GlobalString("person")})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

func (a *Application) DropOnePersonNode(c *cli.Context) error {
	b := api.NewControlClient(a.conn)

	response, err := b.DropOnePersonNode(context.Background(), &api.Person{Name: c.GlobalString("person"), Node: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

func (a *Application) DropMultiPersonBroadcast(c *cli.Context) error {
	b := api.NewControlClient(a.conn)
	multiPerson := &api.MultiPerson{}

	personList := a.parsePersons(c.GlobalString("person"))
	for _, v := range personList {
		multiPerson.Persons = append(multiPerson.Persons, &api.Person{Name: v})
	}
	response, err := b.DropMultiPersonBroadcast(context.Background(), multiPerson)
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

func (a *Application) DropMultiPersonNode(c *cli.Context) error {
	b := api.NewControlClient(a.conn)
	multiPerson := &api.MultiPerson{}

	personList := a.parsePersons(c.GlobalString("person"))
	for _, v := range personList {
		multiPerson.Persons = append(multiPerson.Persons, &api.Person{Name: v, Node: c.GlobalString("node")})
	}
	response, err := b.DropMultiPersonNode(context.Background(), multiPerson)
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

func (a *Application) InsertOnePersonNode(c *cli.Context) error {
	b := api.NewControlClient(a.conn)

	response, err := b.InsertOnePersonNode(context.Background(), &api.Person{Name: c.GlobalString("person"), Node: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

func (a *Application) InsertMultiPersonNode(c *cli.Context) error {
	b := api.NewControlClient(a.conn)
	multiPerson := &api.MultiPerson{}

	personList := a.parsePersons(c.GlobalString("person"))
	for _, v := range personList {
		multiPerson.Persons = append(multiPerson.Persons, &api.Person{Name: v, Node: c.GlobalString("node")})
	}
	response, err := b.InsertMultiPersonNode(context.Background(), multiPerson)
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

func (a *Application) parsePersons(s string) []string {
	slice := strings.Split(s, ",")

	return slice
}
