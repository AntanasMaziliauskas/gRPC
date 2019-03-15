package person

import (
	"context"
	"fmt"
	"log"

	"github.com/AntanasMaziliauskas/grpc/api"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

//DataFromMgo structure holds values od Data, ID and Mgo
type DataFromMgo struct {
	Data []Person
	ID   string
	Mgo  *mgo.Collection
}

//Init function connects to database and sets session
func (d *DataFromMgo) Init() error {
	if err := d.connectToDB(); err != nil {
		log.Fatal("Error while connecting to Database: ", err)
	}

	return nil
}

//ListPersons function returns a list of all persons
func (d *DataFromMgo) ListPersons(ctx context.Context, in *api.Empty) (*api.MultiPerson, error) {
	var err error

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
	result := []Person{}
	var ids []bson.ObjectId

	for _, v := range in.Persons {
		if !bson.IsObjectIdHex(v.Id) {
			log.Println("Provided ID is invalid")

			continue
		}
		ids = append(ids, bson.ObjectIdHex(v.Id))
	}
	if err := d.Mgo.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&result); err != nil {
		fmt.Println("Error while trying to look into DB: ", err)

		return &api.MultiPerson{}, nil
	}
	if len(result) < 1 {
		fmt.Println("Unable to locate given persons")

	}
	for _, k := range result {
		listOfData.Persons = append(listOfData.Persons, &api.Person{Id: k.ID.Hex(), Name: k.Name, Age: k.Age, Profession: k.Profession, Node: d.ID})
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

//UpsertOnePerson adds person to slice
func (d *DataFromMgo) UpsertOnePerson(ctx context.Context, in *api.Person) (*api.Empty, error) {
	var p Person

	if !bson.IsObjectIdHex(in.Id) {
		fmt.Println("Provided ID is invalid")

		return &api.Empty{}, nil
	}
	p = Person{
		ID:         bson.ObjectIdHex(in.Id),
		Name:       in.Name,
		Age:        in.Age,
		Profession: in.Profession,
	}

	selector := bson.M{"_id": p.ID}
	upsertdata := bson.M{"$set": p}
	if _, err := d.Mgo.Upsert(selector, upsertdata); err != nil {
		fmt.Println("Error while trying to upsert: ", err)

		return &api.Empty{}, nil
	}
	fmt.Println("Person successfully upserted")

	return &api.Empty{}, nil
}

//UpsertMultiPerson adds multiple persons to a slice
func (d *DataFromMgo) UpsertMultiPerson(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {
	var p Person

	for _, v := range in.Persons {
		if !bson.IsObjectIdHex(v.Id) {
			fmt.Println("Provided ID is invalid")

			return &api.Empty{}, nil
		}
		p = Person{
			ID:         bson.ObjectIdHex(v.Id),
			Name:       v.Name,
			Age:        v.Age,
			Profession: v.Profession,
		}

		selector := bson.M{"_id": p.ID}
		upsertdata := bson.M{"$set": p}
		if _, err := d.Mgo.Upsert(selector, upsertdata); err != nil {
			fmt.Println("Error while trying to upsert", v.Id, ": ", err)

			continue
		}
		fmt.Println("Person ", v.Name, " successfully upserted")

	}

	return &api.Empty{}, nil
}

//connectTODB function connects to Mongo dabase.
func (d *DataFromMgo) connectToDB() error {

	session, err := mgo.Dial("192.168.200.244:27017")
	if err != nil {
		panic(err)
	}

	//defer session.Close()
	//TODO: Close the session when Node goes offline
	session.SetMode(mgo.Monotonic, true)

	// Drop Database
	/*if IsDrop {
		err = session.DB("test").DropDatabase()
		if err != nil {
			panic(err)
		}
	}*/

	// Collection People
	d.Mgo = session.DB(d.ID).C("people")
	return nil

	/*
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

		d.Mgo = session.DB(d.ID).C("Persons")

		return nil // returininam error*/
}

//Ping function updates LastSeen for the Node that pinged the server
func (d *DataFromMgo) Ping(ctx context.Context, in *api.PingMessage) (*api.Empty, error) {
	return &api.Empty{}, nil
}
