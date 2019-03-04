package person

import (
	"context"
	"errors"
	"fmt"

	"github.com/AntanasMaziliauskas/grpc/api"
)

type DataFromMem struct {
	Data []Person
	ID   string
}

//Init function does nothing
func (d *DataFromMem) Init() error {

	return nil
}

//ListPersons function returns a list of all persons
func (d *DataFromMem) ListPersons(ctx context.Context, in *api.Empty) (*api.MultiPerson, error) {
	listOfData := &api.MultiPerson{}

	for _, v := range d.Data {
		listOfData.Persons = append(listOfData.Persons, &api.Person{
			Name:       v.Name,
			Age:        v.Age,
			Profession: v.Profession,
			Node:       d.ID,
		})
	}
	/*if len(listOfData.Persons) < 1 {
		err := errors.New("There are no persons in this Node")

		return listOfData, err
	}*/

	return listOfData, nil
}

//GetOnePerson function looks for person and returns it if found
func (d *DataFromMem) GetOnePerson(ctx context.Context, in *api.Person) (*api.Person, error) {

	found, err := sliceContainsString(in.Name, d.Data)

	return &api.Person{Name: found.Name, Age: found.Age, Profession: found.Profession, Node: d.ID}, err
}

//GetMultiPerson function looks for multiple persons and returns if found
func (d *DataFromMem) GetMultiPerson(ctx context.Context, in *api.MultiPerson) (*api.MultiPerson, error) {
	listOfData := &api.MultiPerson{}

	for _, k := range in.Persons {
		found, _ := sliceContainsString(k.Name, d.Data)
		listOfData.Persons = append(listOfData.Persons, &api.Person{Name: found.Name, Age: found.Age, Profession: found.Profession, Node: d.ID})
	}
	if len(listOfData.Persons) < 1 {
		err := errors.New("Unable to locate given persons")

		return listOfData, err
	}
	return listOfData, nil
}

//DropOnePerson removes given person from the slice
func (d *DataFromMem) DropOnePerson(ctx context.Context, in *api.Person) (*api.Empty, error) {
	newData := []Person{}
	fmt.Println(d.Data)
	for k, v := range d.Data {
		if v.Name == in.Name {
			newData = append(newData, d.Data[:k]...)
			newData = append(newData, d.Data[k+1:]...)
			d.Data = newData
			fmt.Println(d.Data)

			return &api.Empty{Response: "Person Successfully dropped"}, nil
		}
	}
	fmt.Println(d.Data)
	err := errors.New("Unable to locate given person")

	return &api.Empty{}, err
}

//DropMultiPerson removes given persons from the slice
func (d *DataFromMem) DropMultiPerson(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {
	var success bool

	fmt.Println(d.Data)
	for _, k := range in.Persons {
		newData := []Person{}
		for i, v := range d.Data {
			if v.Name == k.Name {
				newData = append(newData, d.Data[:i]...)
				newData = append(newData, d.Data[i+1:]...)
				success = true
				d.Data = newData
			}
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
func (d *DataFromMem) InsertOnePerson(ctx context.Context, in *api.Person) (*api.Empty, error) {

	fmt.Println(d.Data)
	for _, v := range d.Data {
		if v.Name == in.Name {
			err := errors.New("Given person already exist")
			fmt.Println(d.Data)

			return &api.Empty{}, err
		}
	}
	d.Data = append(d.Data, Person{Name: in.Name, Age: in.Age, Profession: in.Profession})
	fmt.Println(d.Data)

	return &api.Empty{Response: "Person successfully inserted"}, nil
}

//InsertMultiPerson adds multiple persons to a slice
func (d *DataFromMem) InsertMultiPerson(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {
	var duplicate string
	var inserted []string
	var dup bool

	fmt.Println(d.Data)
	for _, v := range in.Persons {
		dup = false
		for _, k := range d.Data {
			if v.Name == k.Name {
				dup = true
				duplicate = duplicate + v.Name + " "
			}
		}
		if !dup {
			inserted = append(inserted, v.Name)
			d.Data = append(d.Data, Person{Name: v.Name, Age: v.Age, Profession: v.Profession})
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

//ReceivedPing function updates LastSeen for the Node that pinged the server
func (d *DataFromMem) Ping(ctx context.Context, in *api.PingMessage) (*api.Empty, error) {
	//fmt.Println("I Got Pinged From ", in.Id)

	//TODO: Ka darome, jeigu Node yra dropintas ir visvien pingina?
	/*for k := range g.Nodes {
		if k == in.Id {
			fmt.Println("I Got Pinged From ", in.Id)
			g.Nodes[in.Id].LastSeen = time.Now()
		}
	}*/

	return &api.Empty{}, nil
}
