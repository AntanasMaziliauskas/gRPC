package control

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/AntanasMaziliauskas/grpc/api"
	"github.com/globalsign/mgo/bson"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

type Application struct {
	client api.ControlClient
}

type Person struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Name       string
	Age        int64
	Profession string
}

//Init function connects to a server
func (a *Application) Init() {
	//TODO: Ar paduodame serverio adresa per flag?
	source := "0.0.0.0:7778"
	conn, err := grpc.Dial(source, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	a.client = api.NewControlClient(conn)
}

//ListPersonsBroadcast function get a list of persons per every Node
func (a *Application) ListPersonsBroadcast(c *cli.Context) error {
	var (
		response *api.MultiPerson
		err      error
	)

	if response, err = a.client.ListPersonsBroadcast(context.Background(), &api.Empty{}); err != nil {
		log.Fatalf("Error when calling ListPersonsBroadcast: %s", err)
	}

	log.Printf("Response: \n")
	for _, v := range response.Persons {
		log.Println(v)
		log.Println(v.Id)
	}
	//log.Println("Response: ", response)

	return nil
}

//ListPersonsNode return list of persons from specified Node
func (a *Application) ListPersonsNode(c *cli.Context) error {
	response, err := a.client.ListPersonsNode(context.Background(), &api.NodeInfo{Id: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling GetOnePersonBroadcast: %s", err)
	}
	log.Println("Response: ", response)
	for _, v := range response.Persons {
		log.Println(v.Name)
	}

	return nil
}

//GetOnePersonBroadcast return information about the preson
//Looks through every Node that's connected to the server
func (a *Application) GetOnePersonBroadcast(c *cli.Context) error {
	response, err := a.client.GetOnePersonBroadcast(context.Background(), &api.Person{Name: c.GlobalString("person")})
	if err != nil {
		log.Fatalf("Error when calling GetOnePersonBroadcast: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

//DropNode deletes Node from the server
func (a *Application) DropNode(c *cli.Context) error {
	response, err := a.client.DropNode(context.Background(), &api.NodeInfo{Id: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling DropNode: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

//ListNodes returns a list of Nodes that are connected to the server
func (a *Application) ListNodes(c *cli.Context) error {
	response, err := a.client.ListNodes(context.Background(), &api.Empty{})
	if err != nil {
		log.Fatalf("Error when calling ListNodes: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

//GetOnePersonNode return information about the person from specified Node
func (a *Application) GetOnePersonNode(c *cli.Context) error {
	response, err := a.client.GetOnePersonNode(context.Background(), &api.Person{Name: c.GlobalString("person"), Node: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling GetOnePersonNode: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

//GetMultiPersonBroadcast returns information about multiple persons
//Looks through every Node that's connected to the server
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

//GetMultiPersonNode returns multiple persons from specified Node
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

//DropOnePersonBroadcast deletes person from any Node that is connected to the server
func (a *Application) DropOnePersonBroadcast(c *cli.Context) error {
	response, err := a.client.DropOnePersonBroadcast(context.Background(), &api.Person{Name: c.GlobalString("person")})
	if err != nil {
		log.Fatalf("Error when calling DropOnePersonBroadcast: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

//DropOnePersonNode deletes person from specified Node
func (a *Application) DropOnePersonNode(c *cli.Context) error {
	response, err := a.client.DropOnePersonNode(context.Background(), &api.Person{Name: c.GlobalString("person"), Node: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling DropOnePersonNode: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

//DropMultiPersonBroadcast drops multiple persons going through every Node that's connected to the server
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

//DropMultiPersonNode deleted multiple persons from specified Node
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

//InsertOnePersonNode adds given person to specified Node
func (a *Application) InsertOnePersonNode(c *cli.Context) error {
	personList := a.parsePerson(c.GlobalString("person"))
	response, err := a.client.InsertOnePersonNode(context.Background(), &api.Person{Name: personList[0].Name, Age: personList[0].Age, Profession: personList[0].Profession, Node: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling InsertOnePersonNode: %s", err)
	}
	log.Println("Response: ", response)

	return nil
}

//InsertMultiPersonNode adds multiple persons to specified Node
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

//MoveOnePerson function moves specified person to specific Node
func (a *Application) MoveOnePerson(c *cli.Context) error {
	response, err := a.client.GetOnePersonBroadcast(context.Background(), &api.Person{Name: c.GlobalString("person")})
	if err != nil {
		log.Fatalf("Error when calling MoveOnePerson: %s", err)
	}
	if response.Node != c.GlobalString("node") {
		a.client.InsertOnePersonNode(context.Background(), &api.Person{
			Name:       response.Name,
			Age:        response.Age,
			Profession: response.Profession,
			Node:       c.GlobalString("node")})
		if err != nil {
			log.Fatalf("Error when calling MoveOnePerson: %s", err)
		}
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

	persons := strings.Split(s, ",")
	for _, k := range persons {
		personSlice := strings.Split(k, ".")
		fmt.Println(len(personSlice))
		if len(personSlice) > 1 {
			age, _ = strconv.ParseInt(personSlice[1], 10, 32)
		}
		list = append(list, Person{Name: personSlice[0], Age: age, Profession: personSlice[2]})
	}
	return list
}
