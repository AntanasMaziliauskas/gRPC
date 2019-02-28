package control

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/AntanasMaziliauskas/grpc/api"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

type Application struct {
	client api.ControlClient
}
type Person struct {
	Name       string
	Age        int64
	Profession string
}

func (a *Application) Init() {
	//TODO: Ar paduodame serverio adresa per flag?
	source := "0.0.0.0:7778"
	conn, err := grpc.Dial(source, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	a.client = api.NewControlClient(conn)
}

func (a *Application) ListPersonsBroadcast(c *cli.Context) error {
	response, err := a.client.ListPersonsBroadcast(context.Background(), &api.Empty{})
	if err != nil {
		log.Fatalf("Error when calling ListPersonsBroadcast: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}
func (a *Application) ListPersonsNode(c *cli.Context) error {
	response, err := a.client.ListPersonsNode(context.Background(), &api.NodeInfo{Id: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling GetOnePersonBroadcast: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}
func (a *Application) GetOnePersonBroadcast(c *cli.Context) error {
	response, err := a.client.GetOnePersonBroadcast(context.Background(), &api.Person{Name: c.GlobalString("person")})
	if err != nil {
		log.Fatalf("Error when calling GetOnePersonBroadcast: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

func (a *Application) DropNode(c *cli.Context) error {
	response, err := a.client.DropNode(context.Background(), &api.NodeInfo{Id: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling DropNode: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

func (a *Application) ListNodes(c *cli.Context) error {
	response, err := a.client.ListNodes(context.Background(), &api.Empty{})
	if err != nil {
		log.Fatalf("Error when calling ListNodes: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

func (a *Application) GetOnePersonNode(c *cli.Context) error {
	response, err := a.client.GetOnePersonNode(context.Background(), &api.Person{Name: c.GlobalString("person"), Node: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling GetOnePersonNode: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

func (a *Application) GetMultiPersonBroadcast(c *cli.Context) error {
	multiPerson := &api.MultiPerson{}
	personList := a.parsePersons(c.GlobalString("person"))

	for _, v := range personList {
		multiPerson.Persons = append(multiPerson.Persons, &api.Person{Name: v})
	}
	response, err := a.client.GetMultiPersonBroadcast(context.Background(), multiPerson)
	if err != nil {
		log.Fatalf("Error when calling GetMultiPersonBroadcast: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

func (a *Application) GetMultiPersonNode(c *cli.Context) error {
	multiPerson := &api.MultiPerson{}
	personList := a.parsePersons(c.GlobalString("person"))

	for _, v := range personList {
		multiPerson.Persons = append(multiPerson.Persons, &api.Person{Name: v, Node: c.GlobalString("node")})
	}
	response, err := a.client.GetMultiPersonNode(context.Background(), multiPerson)
	if err != nil {
		log.Fatalf("Error when calling GetMultiPersonNode: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

//TODO: Sutvarkyti, kad grazintu zinute ar pavyko ar ne dropinti
func (a *Application) DropOnePersonBroadcast(c *cli.Context) error {
	response, err := a.client.DropOnePersonBroadcast(context.Background(), &api.Person{Name: c.GlobalString("person")})
	if err != nil {
		log.Fatalf("Error when calling DropOnePersonBroadcast: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

//TODO: sutvarkyti, kad grazintu zinute ar pavyko ar ne dropinti
func (a *Application) DropOnePersonNode(c *cli.Context) error {
	response, err := a.client.DropOnePersonNode(context.Background(), &api.Person{Name: c.GlobalString("person"), Node: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling DropOnePersonNode: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

//TODO: sutvarkyti, kad grazintu zinute ar pavyko ar ne dropinti
func (a *Application) DropMultiPersonBroadcast(c *cli.Context) error {
	multiPerson := &api.MultiPerson{}

	personList := a.parsePersons(c.GlobalString("person"))
	for _, v := range personList {
		multiPerson.Persons = append(multiPerson.Persons, &api.Person{Name: v})
	}
	response, err := a.client.DropMultiPersonBroadcast(context.Background(), multiPerson)
	if err != nil {
		log.Fatalf("Error when calling DropMultiPersonBroadcast: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

//TODO: sutvarkyti, kad grazintu zinute ar pavyko ar ne dropinti
func (a *Application) DropMultiPersonNode(c *cli.Context) error {
	multiPerson := &api.MultiPerson{}

	personList := a.parsePersons(c.GlobalString("person"))
	for _, v := range personList {
		multiPerson.Persons = append(multiPerson.Persons, &api.Person{Name: v, Node: c.GlobalString("node")})
	}
	response, err := a.client.DropMultiPersonNode(context.Background(), multiPerson)
	if err != nil {
		log.Fatalf("Error when calling DropMultiPersonNode: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

//TODO: sutvarkyti, kad grazintu zinute ar pavyko insert ar ne
func (a *Application) InsertOnePersonNode(c *cli.Context) error {
	personList := a.parsePerson(c.GlobalString("person"))
	response, err := a.client.InsertOnePersonNode(context.Background(), &api.Person{Name: personList[0].Name, Age: personList[0].Age, Profession: personList[0].Profession, Node: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling InsertOnePersonNode: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

//TODO: sutvarkyti, kad grazintu zinute ar pavyko insert ar ne
func (a *Application) InsertMultiPersonNode(c *cli.Context) error {
	multiPerson := &api.MultiPerson{}

	personList := a.parsePerson(c.GlobalString("person"))
	for _, v := range personList {
		multiPerson.Persons = append(multiPerson.Persons, &api.Person{Name: v.Name, Age: v.Age, Profession: v.Profession, Node: c.GlobalString("node")})
	}
	response, err := a.client.InsertMultiPersonNode(context.Background(), multiPerson)
	if err != nil {
		log.Fatalf("Error when calling InsertMultipleNode: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

func (a *Application) MoveOnePerson(c *cli.Context) error {
	response, err := a.client.GetOnePersonBroadcast(context.Background(), &api.Person{Name: c.GlobalString("person")})
	if err != nil {
		log.Fatalf("Error when calling MoveOnePerson: %s", err)
	}
	if response.Node != c.GlobalString("node") {
		//TODO: Koki patikrinima daryti Insert funkcijai?
		a.client.InsertOnePersonNode(context.Background(), &api.Person{
			Name:       response.Name,
			Age:        response.Age,
			Profession: response.Profession,
			Node:       c.GlobalString("node")})
		if err != nil {
			log.Fatalf("Error when calling MoveOnePerson: %s", err)
		}
		//TODO: Koki patikrinima daryti su Drop?
		a.client.DropOnePersonNode(context.Background(), response)
		if err != nil {
			log.Fatalf("Error when calling DropOnePersonNode: %s", err)
		}
	} else {
		fmt.Println("Person already belongs to given Node.")
	}

	return nil
}

func (a *Application) parsePersons(s string) []string {
	slice := strings.Split(s, ",")

	return slice
}

//Convert person flag into Person slice
func (a *Application) parsePerson(s string) []Person {
	var (
		list []Person
		age  int64
	)

	persons := strings.Split(s, ".")
	for _, k := range persons {
		personSlice := strings.Split(k, ",")
		if len(personSlice) > 1 {
			age, _ = strconv.ParseInt(personSlice[1], 10, 32)
		}
		list = append(list, Person{Name: personSlice[0], Age: age, Profession: personSlice[2]})
	}
	return list
}
