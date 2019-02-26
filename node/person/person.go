package person

import (
	"encoding/json"
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
