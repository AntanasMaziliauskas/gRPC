package control

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/AntanasMaziliauskas/grpc/api"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

//Application structure holds the value of client
type Application struct {
	client api.ControlClient
}

//Person structure holds the clues of ID, name, age and profession
type Person struct {
	ID         string
	Name       string
	Age        int64
	Profession string
}

//Init function connects to a server
func (a *Application) Init() {
	//TODO: Ar paduodame serverio adresa per flag?
	source := "192.168.99.1:7778"
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
	//Startas
	startTime := time.Now()

	if response, err = a.client.ListPersonsBroadcast(context.Background(), &api.Empty{}); err != nil {
		log.Fatalf("Error when calling ListPersonsBroadcast: %s", err)
	}
	//Finishas
	duration := time.Now().Sub(startTime)
	log.Println(duration)

	log.Printf("Response: \n")
	for _, v := range response.Persons {
		log.Println(v)
	}

	return nil
}

//ListPersonsNode return list of persons from specified Node
func (a *Application) ListPersonsNode(c *cli.Context) error {
	if c.GlobalString("node") == "" {
		log.Fatalf("There is NO NODE provided.")
	}
	response, err := a.client.ListPersonsNode(context.Background(), &api.NodeInfo{Id: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling GetOnePersonBroadcast: %s", err)
	}
	for _, v := range response.Persons {
		log.Println(v)
	}

	return nil
}

//GetOnePersonBroadcast return information about the preson
//Looks through every Node that's connected to the server
func (a *Application) GetOnePersonBroadcast(c *cli.Context) error {
	if c.GlobalString("person") == "" {
		log.Fatalf("There is NO PERSON provided.")
	}
	response, err := a.client.GetOnePersonBroadcast(context.Background(), &api.Person{Id: c.GlobalString("person")})
	if err != nil {
		log.Fatalf("Error when calling GetOnePersonBroadcast: %s", err)
	}
	log.Println(response)

	return nil
}

//DropNode deletes Node from the server
func (a *Application) DropNode(c *cli.Context) error {
	if c.GlobalString("node") == "" {
		log.Fatalf("There is NO NODE provided.")
	}
	_, err := a.client.DropNode(context.Background(), &api.NodeInfo{Id: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling DropNode: %s", err)
	}

	return nil
}

//ListNodes returns a list of Nodes that are connected to the server
func (a *Application) ListNodes(c *cli.Context) error {
	response, err := a.client.ListNodes(context.Background(), &api.Empty{})
	if err != nil {
		log.Fatalf("Error when calling ListNodes: %s", err)
	}
	//fmt.Println("Response:")
	for _, v := range response.Nodes {
		if v.Isonline {
			fmt.Println(v.Id, " - ONLINE")
		} else {
			fmt.Println(v.Id, " - OFFLINE")
		}

	}
	return nil
}

//GetOnePersonNode return information about the person from specified Node
func (a *Application) GetOnePersonNode(c *cli.Context) error {
	if c.GlobalString("node") == "" {
		log.Fatalf("There is NO NODE provided.")
	}
	response, err := a.client.GetOnePersonNode(context.Background(), &api.Person{Id: c.GlobalString("person"), Node: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling GetOnePersonNode: %s", err)
	}
	log.Println(response)

	return nil
}

//GetMultiPersonBroadcast returns information about multiple persons
//Looks through every Node that's connected to the server
func (a *Application) GetMultiPersonBroadcast(c *cli.Context) error {
	if c.GlobalString("person") == "" {
		log.Fatalf("There is NO PERSON provided.")
	}
	multiPerson := &api.MultiPerson{}
	personList := parsePersons(c.GlobalString("person"))

	for _, v := range personList {
		multiPerson.Persons = append(multiPerson.Persons, &api.Person{Id: v.ID})
	}
	response, err := a.client.GetMultiPersonBroadcast(context.Background(), multiPerson)
	if err != nil {
		log.Fatalf("Error when calling GetMultiPersonBroadcast: %s", err)
	}
	for _, v := range response.Persons {
		log.Println(v)
	}

	return nil
}

//GetMultiPersonNode returns multiple persons from specified Node
func (a *Application) GetMultiPersonNode(c *cli.Context) error {
	multiPerson := &api.MultiPerson{}
	personList := parsePersons(c.GlobalString("person"))
	if c.GlobalString("node") == "" || c.GlobalString("person") == "" {
		log.Fatalf("There is NO NODE provided.")
	}

	for _, v := range personList {
		multiPerson.Persons = append(multiPerson.Persons, &api.Person{Id: v.ID, Node: c.GlobalString("node")})
	}
	response, err := a.client.GetMultiPersonNode(context.Background(), multiPerson)
	if err != nil {
		log.Fatalf("Error when calling GetMultiPersonNode: %s", err)
	}
	for _, v := range response.Persons {
		log.Println(v)
	}

	return nil
}

//DropOnePersonBroadcast deletes person from any Node that is connected to the server
func (a *Application) DropOnePersonBroadcast(c *cli.Context) error {
	_, err := a.client.DropOnePersonBroadcast(context.Background(), &api.Person{Id: c.GlobalString("person")})
	if err != nil {
		log.Fatalf("Error when calling DropOnePersonBroadcast: %s", err)
	}

	return nil
}

//DropOnePersonNode deletes person from specified Node
func (a *Application) DropOnePersonNode(c *cli.Context) error {
	if c.GlobalString("node") == "" {
		log.Fatalf("There is NO NODE provided.")
	}
	_, err := a.client.DropOnePersonNode(context.Background(), &api.Person{Id: c.GlobalString("person"), Node: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling DropOnePersonNode: %s", err)
	}

	return nil
}

//DropMultiPersonBroadcast drops multiple persons going through every Node that's connected to the server
func (a *Application) DropMultiPersonBroadcast(c *cli.Context) error {
	multiPerson := &api.MultiPerson{}

	personList := parsePersons(c.GlobalString("person"))
	for _, v := range personList {
		multiPerson.Persons = append(multiPerson.Persons, &api.Person{Id: v.ID})
	}
	_, err := a.client.DropMultiPersonBroadcast(context.Background(), multiPerson)
	if err != nil {
		log.Fatalf("Error when calling DropMultiPersonBroadcast: %s", err)
	}

	return nil
}

//DropMultiPersonNode deleted multiple persons from specified Node
func (a *Application) DropMultiPersonNode(c *cli.Context) error {
	multiPerson := &api.MultiPerson{}
	if c.GlobalString("node") == "" {
		log.Fatalf("There is NO NODE provided.")
	}
	personList := parsePersons(c.GlobalString("person"))
	for _, v := range personList {
		multiPerson.Persons = append(multiPerson.Persons, &api.Person{Id: v.ID, Node: c.GlobalString("node")})
	}
	_, err := a.client.DropMultiPersonNode(context.Background(), multiPerson)
	if err != nil {
		log.Fatalf("Error when calling DropMultiPersonNode: %s", err)
	}

	return nil
}

//UpsertOnePersonNode adds given person to specified Node
func (a *Application) UpsertOnePersonNode(c *cli.Context) error {
	if c.GlobalString("node") == "" {
		log.Fatalf("There is NO NODE provided.")
	}
	personList := parsePerson(c.GlobalString("person"))

	_, err := a.client.UpsertOnePersonNode(context.Background(), &api.Person{Id: personList.ID, Name: personList.Name, Age: personList.Age, Profession: personList.Profession, Node: c.GlobalString("node")})
	if err != nil {
		log.Fatalf("Error when calling UpsertOnePersonNode: %s", err)
	}

	return nil
}

//UpsertMultiPersonNode adds multiple persons to specified Node
func (a *Application) UpsertMultiPersonNode(c *cli.Context) error {
	multiPerson := &api.MultiPerson{}
	if c.GlobalString("node") == "" {
		log.Fatalf("There is NO NODE provided.")
	}
	personList := parsePersons(c.GlobalString("person"))
	for _, v := range personList {
		multiPerson.Persons = append(multiPerson.Persons, &api.Person{Id: v.ID, Name: v.Name, Age: v.Age, Profession: v.Profession, Node: c.GlobalString("node")})
	}
	_, err := a.client.UpsertMultiPersonNode(context.Background(), multiPerson)
	if err != nil {
		log.Fatalf("Error when calling UpsertMultipleNode: %s", err)
	}

	return nil
}

//MoveOnePerson function moves specified person to specific Node
func (a *Application) MoveOnePerson(c *cli.Context) error {
	if c.GlobalString("node") == "" {
		log.Fatalf("There is NO NODE provided.")
	}
	response, err := a.client.GetOnePersonBroadcast(context.Background(), &api.Person{Id: c.GlobalString("person")})
	if err != nil {
		log.Fatalf("Error when calling MoveOnePerson: %s", err)
	}
	if response.Node != c.GlobalString("node") {
		_, err = a.client.UpsertOnePersonNode(context.Background(), &api.Person{
			Id:         response.Id,
			Name:       response.Name,
			Age:        response.Age,
			Profession: response.Profession,
			Node:       c.GlobalString("node")})
		if err != nil {
			log.Fatalf("Error when trying to insert person into Node: %s", err)
		}
		_, err = a.client.DropOnePersonNode(context.Background(), response)
		if err != nil {
			log.Fatalf("Error when calling DropOnePersonNode: %s", err)
		}
	} else {
		fmt.Println("Person already belongs to given Node.")

		return nil
	}

	fmt.Println("Person successfully moved.")
	return nil
}

//parsePerson function parses string into Person structure
func parsePerson(s string) Person {
	var (
		person Person
		age    int64
		err    error
	)

	slice := strings.Split(s, ".")
	if len(slice) > 1 {
		if age, err = strconv.ParseInt(slice[2], 10, 32); err != nil {
			fmt.Println("Error while concerting string into int: ", err)
		}
		person = Person{
			ID:         slice[0],
			Name:       slice[1],
			Age:        age,
			Profession: slice[3],
		}
		return person
	}
	person = Person{
		ID: slice[0],
	}
	return person
}

//parsePersons function parses string into []Person structure
func parsePersons(s string) []Person {
	var (
		list []Person
	)

	persons := strings.Split(s, ",")
	for _, k := range persons {
		one := parsePerson(k)
		list = append(list, one)
	}
	return list
}
