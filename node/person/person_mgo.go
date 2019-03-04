package person

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"

	"github.com/AntanasMaziliauskas/grpc/api"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type DataFromMgo struct {
	Data []Person
	ID   string
	Mgo  *mgo.Collection
}

//Init function does nothing
func (d *DataFromMgo) Init() error {

	//Connect to database and set session
	d.connectToDB() // error isideti

	return nil
}

//ListPersons function returns a list of all persons
func (d *DataFromMgo) ListPersons(ctx context.Context, in *api.Empty) (*api.MultiPerson, error) {
	listOfData := &api.MultiPerson{}
	list := &[]Person

	//TODO: Prideti Node ID i sarasa
	err := d.Mgo.Find(bson.M{}).All(&listOfData.Persons)

	for _, v := range listOfData.Persons {
		fmt.Println(v.Id.Hex())
	}

	/*if len(listOfData.Persons) < 1 {
		err := errors.New("There are no persons in this Node")

		return listOfData, err
	}*/

	return listOfData, err
}

//GetOnePerson function looks for person and returns it if found
func (d *DataFromMgo) GetOnePerson(ctx context.Context, in *api.Person) (*api.Person, error) {
	result := &Person{}

	err := d.Mgo.Find(bson.M{"name": in.Name}).One(&result)
	if err != nil {
		fmt.Println("Neranda", err)
	}

	return &api.Person{Name: result.Name, Age: result.Age, Profession: result.Profession, Node: d.ID}, err
}

//GetMultiPerson function looks for multiple persons and returns if found
func (d *DataFromMgo) GetMultiPerson(ctx context.Context, in *api.MultiPerson) (*api.MultiPerson, error) {
	listOfData := &api.MultiPerson{}

	for _, k := range in.Persons {
		result := &Person{}
		err := d.Mgo.Find(bson.M{"name": k.Name}).One(&result)
		if err != nil {
			fmt.Println("Neranda", err)
		}
		listOfData.Persons = append(listOfData.Persons, &api.Person{Name: result.Name, Age: result.Age, Profession: result.Profession, Node: d.ID})
	}
	if len(listOfData.Persons) < 1 {
		err := errors.New("Unable to locate given persons")

		return listOfData, err
	}
	return listOfData, nil
}

//DropOnePerson removes given person from the slice
func (d *DataFromMgo) DropOnePerson(ctx context.Context, in *api.Person) (*api.Empty, error) {
	//newData := []Person{}
	//fmt.Println(d.Data)

	err := d.Mgo.Remove(bson.M{"name": in.Name})
	if err != nil {

		return &api.Empty{}, err
	}

	return &api.Empty{Response: "Person successfully dropped"}, err
}

//DropMultiPerson removes given persons from the slice
func (d *DataFromMgo) DropMultiPerson(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {
	var success bool

	fmt.Println(d.Data)
	for _, k := range in.Persons {
		err := d.Mgo.Remove(bson.M{"name": k.Name})
		if err == nil {
			success = true
		}
	}
	if success {
		fmt.Println(d.Data)

		return &api.Empty{Response: "Persons successfully dropped"}, nil
	}
	fmt.Println(d.Data)
	err := errors.New("Unable to locate given person")

	return &api.Empty{}, err
}

//InsertOnePerson adds person to slice
func (d *DataFromMgo) InsertOnePerson(ctx context.Context, in *api.Person) (*api.Empty, error) {
	listOfData := &api.MultiPerson{}

	err := d.Mgo.Find(bson.M{}).All(&listOfData.Persons)
	fmt.Println(listOfData)

	for _, v := range listOfData.Persons {
		if v.Name == in.Name {
			err := errors.New("Given person already exist")
			fmt.Println(listOfData)

			return &api.Empty{}, err
		}
	}
	if err := d.Mgo.Insert(&Person{ID: bson.NewObjectId(), Name: in.Name, Age: in.Age, Profession: in.Profession}); err != nil {
		return &api.Empty{}, err
	}
	d.Mgo.Find(bson.M{}).All(&listOfData.Persons)
	fmt.Println(listOfData)

	return &api.Empty{Response: "Person successfully inserted"}, err
}

//InsertMultiPerson adds multiple persons to a slice
func (d *DataFromMgo) InsertMultiPerson(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {
	var duplicate string
	var inserted []string
	var dup bool

	listOfData := &api.MultiPerson{}

	err := d.Mgo.Find(bson.M{}).All(&listOfData.Persons)
	if err != nil {

		return &api.Empty{}, err
	}
	fmt.Println(listOfData)
	for _, v := range in.Persons {
		for _, k := range listOfData.Persons {
			if v.Name == k.Name {
				dup = true
				duplicate = duplicate + v.Name + " "
			}
		}
		if !dup {
			inserted = append(inserted, v.Name)
			if err := d.Mgo.Insert(&Person{Name: v.Name, Age: v.Age, Profession: v.Profession}); err != nil {
				return &api.Empty{}, err
			}
		}
	}
	if len(duplicate) < 1 {
		fmt.Println(d.Data)

		return &api.Empty{Response: "Successfully inserted"}, nil
	}
	fmt.Println(d.Data)
	message := duplicate + " are already inserted."
	return &api.Empty{Response: message}, nil
}

func (d *DataFromMgo) connectToDB() error {

	tlsConfig := &tls.Config{}
	dialInfo := &mgo.DialInfo{
		Addrs: []string{"persondb-shard-00-00-mimet.mongodb.net:27017",
			"persondb-shard-00-01-mimet.mongodb.net:27017",
			"persondb-shard-00-02-mimet.mongodb.net:27017"},
		Database: "admin",
		Username: "node",
		Password: "node01",
	}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Println(err)
	}

	d.Mgo = session.DB("Test").C("Persons")

	return nil // returininam error
}

//ReceivedPing function updates LastSeen for the Node that pinged the server
func (d *DataFromMgo) Ping(ctx context.Context, in *api.PingMessage) (*api.Empty, error) {
	fmt.Println("I Got Pinged")

	//TODO: Ka darome, jeigu Node yra dropintas ir visvien pingina?
	/*for k := range g.Nodes {
		if k == in.Id {
			fmt.Println("I Got Pinged From ", in.Id)
			g.Nodes[in.Id].LastSeen = time.Now()
		}
	}*/

	return &api.Empty{}, nil
}
