package person

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AntanasMaziliauskas/grpc/api"
	"github.com/globalsign/mgo/bson"

	//"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DataFromMgo structure holds values od Data, ID and Mgo
type DataFromMgo struct {
	Data []Person
	ID   string
	Mgo  *mongo.Collection
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
	var list []*Person

	listOfData := &api.MultiPerson{}
	//list := []Person{}
	//Start
	startTime := time.Now()

	cur, err := d.Mgo.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println(err)
	}

	//TODO: Iskarto i Person slice krauti
	for cur.Next(context.TODO()) {
		var elem Person
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
		}
		list = append(list, &elem)
	}
	cur.Close(context.TODO())
	//	fmt.Printf("Found: %+v\n", list)
	//	if err = d.Mgo.Find(bson.M{}).All(list); err != nil {
	//		log.Println("Error while trying to find data: ", err)
	//	}
	//finish
	duration := time.Now().Sub(startTime)
	log.Println("Paieška duombazėj užtruko: ", duration)
	for _, v := range list {
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
	//Start
	startTime := time.Now()
	id, _ := primitive.ObjectIDFromHex(in.Id)
	if err := d.Mgo.FindOne(context.TODO(), primitive.D{{"_id", id}}).Decode(&result); err != nil {
		log.Println("Unable to locate given person")

		return &api.Person{}, nil
	}
	//finish
	duration := time.Now().Sub(startTime)
	log.Println("Paieška duombazėj užtruko: ", duration)

	return &api.Person{Id: result.ID.Hex(), Name: result.Name, Age: result.Age, Profession: result.Profession, Node: d.ID}, nil

}

//GetMultiPerson function looks for multiple persons and returns if found
func (d *DataFromMgo) GetMultiPerson(ctx context.Context, in *api.MultiPerson) (*api.MultiPerson, error) {
	listOfData := &api.MultiPerson{}
	result := Person{}
	//var ids []bson.ObjectId

	for _, v := range in.Persons {
		if !bson.IsObjectIdHex(v.Id) {
			log.Println("Provided ID is invalid")

			continue
		}
		id, _ := primitive.ObjectIDFromHex(v.Id)
		if err := d.Mgo.FindOne(context.TODO(), primitive.D{{"_id", id}}).Decode(&result); err != nil {
			fmt.Println("Error while trying to look into DB: ", err)

			return &api.MultiPerson{}, nil
		}
		if len(result.ID) < 1 {
			fmt.Println("Unable to locate given persons")

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
	id, _ := primitive.ObjectIDFromHex(in.Id)
	_, err = d.Mgo.DeleteOne(context.TODO(), primitive.D{{"_id", id}})
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
		id, _ := primitive.ObjectIDFromHex(k.Id)
		_, err := d.Mgo.DeleteOne(context.TODO(), primitive.D{{"_id", id}})
		if err != nil {
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
	//var p Person

	if !bson.IsObjectIdHex(in.Id) {
		fmt.Println("Provided ID is invalid")

		return &api.Empty{}, nil
	}
	/*p = Person{
		ID:         bson.ObjectIdHex(in.Id),
		Name:       in.Name,
		Age:        in.Age,
		Profession: in.Profession,
	}*/
	p := primitive.D{
		{"$set", primitive.D{
			//{"_id", bson.ObjectIdHex(in.Id)},
			{"name", in.Name},
			{"age", in.Age},
			{"profession", in.Profession},
		}},
	}

	findOptions := options.Update()
	b := true
	findOptions.Upsert = &b
	id, _ := primitive.ObjectIDFromHex(in.Id)

	//selector := bson.M{"_id": p.ID}
	//upsertdata := bson.M{"$set": p}
	_, err := d.Mgo.UpdateOne(context.TODO(), primitive.D{{"_id", id}}, p, findOptions)
	if err != nil {
		fmt.Println("Error while trying to upsert: ", err)

		return &api.Empty{}, nil
	}
	fmt.Println("Person successfully upserted")

	return &api.Empty{}, nil
}

//UpsertMultiPerson adds multiple persons to a slice
func (d *DataFromMgo) UpsertMultiPerson(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {
	//var p Person

	for _, v := range in.Persons {
		if !bson.IsObjectIdHex(v.Id) {
			fmt.Println("Provided ID is invalid")

			return &api.Empty{}, nil
		}

		id, _ := primitive.ObjectIDFromHex(v.Id)

		p := primitive.D{
			{"$set", primitive.D{
				//{"_id", bson.ObjectIdHex(in.Id)},
				{"name", v.Name},
				{"age", v.Age},
				{"profession", v.Profession},
			}},
		}

		findOptions := options.Update()
		b := true
		findOptions.Upsert = &b
		//id, _ := primitive.ObjectIDFromHex(v.Id)

		//selector := bson.M{"_id": p.ID}
		//upsertdata := bson.M{"$set": p}
		_, err := d.Mgo.UpdateOne(context.TODO(), primitive.D{{"_id", id}}, p, findOptions)
		if err != nil {
			fmt.Println("Error while trying to upsert", v.Id, ": ", err)

			continue
		}
		fmt.Println("Person ", v.Name, " successfully upserted")

	}

	return &api.Empty{}, nil
}

//connectTODB function connects to Mongo dabase.
func (d *DataFromMgo) connectToDB() error {

	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://192.168.200.244:27017"))
	if err != nil {
		panic(err)
	}
	//session.SetMode(mgo.Monotonic, true)
	d.Mgo = client.Database(d.ID).Collection("people")

	return nil
}

//Ping function updates LastSeen for the Node that pinged the server
func (d *DataFromMgo) Ping(ctx context.Context, in *api.PingMessage) (*api.Empty, error) {
	return &api.Empty{}, nil
}
