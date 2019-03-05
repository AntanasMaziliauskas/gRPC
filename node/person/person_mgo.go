package person

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/AntanasMaziliauskas/grpc/api"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type DataFromMgo struct {
	//	Persons map[bson.ObjectId]Person
	Data []Person
	ID   string
	Mgo  *mgo.Collection
}

//Init function does nothing
func (d *DataFromMgo) Init() error {
	//	d.Persons = make(map[bson.ObjectId]Person)
	//Connect to database and set session
	d.connectToDB() // error isideti

	return nil
}

//ListPersons function returns a list of all persons
func (d *DataFromMgo) ListPersons(ctx context.Context, in *api.Empty) (*api.MultiPerson, error) {
	var err error

	log.Println("Looking for data...")

	listOfData := &api.MultiPerson{}
	list := &[]Person{}
	if err = d.Mgo.Find(bson.M{}).All(list); err != nil {
		log.Println("Error while trying to find data: ", err)
	}

	for _, v := range *list {
		listOfData.Persons = append(listOfData.Persons, &api.Person{Id: v.ID.Hex(),
			Name: v.Name, Age: v.Age, Profession: v.Profession, Node: d.ID})
	}

	if len(listOfData.Persons) < 1 {
		log.Println("This Node has no data.")
	}

	return listOfData, nil
}

//GetOnePerson function looks for person and returns it if found
func (d *DataFromMgo) GetOnePerson(ctx context.Context, in *api.Person) (*api.Person, error) {
	result := &Person{}

	if !bson.IsObjectIdHex(in.Id) {
		log.Println("Provided ID is invalid")

		return &api.Person{}, nil
	}
	if err := d.Mgo.Find(bson.M{"_id": bson.ObjectIdHex(in.Id)}).One(&result); err != nil {
		log.Println("Unable to locate given person")

		return &api.Person{}, nil
	}

	return &api.Person{Id: result.ID.Hex(), Name: result.Name, Age: result.Age, Profession: result.Profession, Node: d.ID}, nil

}

//GetMultiPerson function looks for multiple persons and returns if found
func (d *DataFromMgo) GetMultiPerson(ctx context.Context, in *api.MultiPerson) (*api.MultiPerson, error) {
	listOfData := &api.MultiPerson{}

	for _, k := range in.Persons {
		result := &Person{}
		if !bson.IsObjectIdHex(k.Id) {
			log.Println("Provided ID is invalid")

			continue
		}
		err := d.Mgo.Find(bson.M{"_id": bson.ObjectIdHex(k.Id)}).One(&result)
		if err != nil {
			fmt.Println("Unable to locate person. Error: ", err)

			continue
		}
		listOfData.Persons = append(listOfData.Persons, &api.Person{Id: result.ID.Hex(), Name: result.Name, Age: result.Age, Profession: result.Profession, Node: d.ID})
	}

	return listOfData, nil
}

//DropOnePerson removes given person from the slice
func (d *DataFromMgo) DropOnePerson(ctx context.Context, in *api.Person) (*api.Empty, error) {
	var err error

	if !bson.IsObjectIdHex(in.Id) {
		fmt.Println("Provided ID is invalid")

		return &api.Empty{}, nil
	}
	err = d.Mgo.Remove(bson.M{"_id": bson.ObjectIdHex(in.Id)})
	if err != nil {
		fmt.Println("Unable to locate person. Error: ", err)

		return &api.Empty{}, nil
	}
	fmt.Println("Person successfully dropped")
	return &api.Empty{}, nil
}

//DropMultiPerson removes given persons from the slice
func (d *DataFromMgo) DropMultiPerson(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {
	var success bool

	fmt.Println(d.Data)
	for _, k := range in.Persons {
		if !bson.IsObjectIdHex(k.Id) {
			fmt.Println("Provided ID is invalid")

			continue
		}
		if err := d.Mgo.Remove(bson.M{"_id": bson.ObjectIdHex(k.Id)}); err != nil {
			fmt.Println("Error wile trying to remove ", k.Id, " . Error: ", err)

			continue
		}
		success = true
	}
	if success {
		fmt.Println(d.Data)
		fmt.Println("Persons successfully dropped")

		return &api.Empty{}, nil
	}

	return &api.Empty{}, nil
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
			//BUS OKAY SU KLAIDA> GRAZINAME
			if err := d.Mgo.Insert(&Person{ID: bson.NewObjectId(), Name: v.Name, Age: v.Age, Profession: v.Profession}); err != nil {
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

//connectTODB function connects to Mongo dabase.
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

//Ping function updates LastSeen for the Node that pinged the server
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
