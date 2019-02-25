package person

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type PersonService interface {
	Init() error

	GetOne(name string) (Person, error)
	//	GetMulti(names []string) ([]api.Person, error)

}
type DataFromFile struct {
	Path string
	Data []Person
}

type Person struct {
	Name       string
	Age        int64
	Profession string
}

func (d *DataFromFile) Init() error {
	d.Data, _ = d.readFile()
	fmt.Println("Init veikia")
	return nil
}

//Dokumento nuskaitymas i struktura
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
