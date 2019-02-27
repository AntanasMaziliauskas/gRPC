package person

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/AntanasMaziliauskas/grpc/api"
)

type PersonService interface {
	Init() error
	GetOnePerson(context.Context, *api.Person) (*api.Person, error)
	GetMultiPerson(context.Context, *api.MultiPerson) (*api.MultiPerson, error)
	DropOnePerson(context.Context, *api.Person) (*api.Empty, error)
	DropMultiPerson(context.Context, *api.MultiPerson) (*api.Empty, error)
	InsertOnePerson(context.Context, *api.Person) (*api.Empty, error)
	InsertMultiPerson(context.Context, *api.MultiPerson) (*api.Empty, error)
	GetOne(name string) (Person, error)
	//	GetMulti(names []string) ([]api.Person, error)
}

type DataFromFile struct {
	Path string
	Data []Person
	ID   string
}

type Person struct {
	Name       string
	Age        int64
	Profession string
}

//Init function reads the file
func (d *DataFromFile) Init() error {
	var err error

	d.Data, err = d.readFile()

	return err
}

//readFile function reads the file and adds the content into structure
func (d *DataFromFile) readFile() ([]Person, error) {
	var (
		data     []Person
		err      error
		jsonFile []byte
	)

	if jsonFile, err = ioutil.ReadFile(d.Path); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(jsonFile, &data); err != nil {
		return nil, err
	}

	return data, nil
}

//GetOne function looks through the structure of data and return
func (d *DataFromFile) GetOne(name string) (Person, error) {
	found := sliceContainsString(name, d.Data)
	return found, nil
}

// SliceContainsString will return true if needle has been found in haystack.
func sliceContainsString(needle string, haystack []Person) Person {
	for _, v := range haystack {
		if v.Name == needle {
			return v
		}
	}

	return Person{}
}

//GetOnePersonBroadcast function go through the list of connected Nodes
//requests to look for data with the name given, retrieved the data and return it.
func (d *DataFromFile) GetOnePerson(ctx context.Context, in *api.Person) (*api.Person, error) {

	found := sliceContainsString(in.Name, d.Data)

	return &api.Person{Name: found.Name, Age: found.Age, Profession: found.Profession}, nil
}

func (d *DataFromFile) GetMultiPerson(ctx context.Context, in *api.MultiPerson) (*api.MultiPerson, error) {
	listOfData := &api.MultiPerson{}

	for _, k := range in.Persons {
		found := sliceContainsString(k.Name, d.Data)

		listOfData.Persons = append(listOfData.Persons, &api.Person{Name: found.Name, Age: found.Age, Profession: found.Profession, Node: d.ID})
	}

	return listOfData, nil
}

func (d *DataFromFile) DropOnePerson(ctx context.Context, in *api.Person) (*api.Empty, error) {
	newData := []Person{}
	fmt.Println(d.Data)
	for _, v := range d.Data {
		if v.Name != in.Name {
			newData = append(newData, v)
		}
	}
	d.Data = newData
	fmt.Println(d.Data)
	return &api.Empty{}, nil
}

func (d *DataFromFile) DropMultiPerson(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {

	fmt.Println(d.Data)

	for _, k := range in.Persons {
		newData := []Person{}
		for _, v := range d.Data {
			if v.Name != k.Name {
				newData = append(newData, v)
			}

		}
		d.Data = newData
	}

	fmt.Println(d.Data)
	return &api.Empty{}, nil
}

func (d *DataFromFile) InsertOnePerson(ctx context.Context, in *api.Person) (*api.Empty, error) {

	fmt.Println(d.Data)
	d.Data = append(d.Data, Person{Name: in.Name, Age: in.Age, Profession: in.Profession})
	fmt.Println(d.Data)
	return &api.Empty{}, nil
}

func (d *DataFromFile) InsertMultiPerson(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {

	fmt.Println(d.Data)
	for _, v := range in.Persons {
		d.Data = append(d.Data, Person{Name: v.Name, Age: v.Age, Profession: v.Profession})
	}
	fmt.Println(d.Data)
	return &api.Empty{}, nil
}
