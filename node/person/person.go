package person

import (
	"context"
	"errors"

	"github.com/AntanasMaziliauskas/grpc/api"
	"github.com/globalsign/mgo/bson"
)

type PersonService interface {
	Init() error
	ListPersons(context.Context, *api.Empty) (*api.MultiPerson, error)
	GetOnePerson(context.Context, *api.Person) (*api.Person, error)
	GetMultiPerson(context.Context, *api.MultiPerson) (*api.MultiPerson, error)
	DropOnePerson(context.Context, *api.Person) (*api.Empty, error)
	DropMultiPerson(context.Context, *api.MultiPerson) (*api.Empty, error)
	InsertOnePerson(context.Context, *api.Person) (*api.Empty, error)
	InsertMultiPerson(context.Context, *api.MultiPerson) (*api.Empty, error)
	Ping(context.Context, *api.PingMessage) (*api.Empty, error)
}

type Person struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Name       string        `bson:"name"`
	Age        int64         `bson:"age"`
	Profession string        `bson:"profession"`
}

// SliceContainsString will return true if needle has been found in haystack.
func sliceContainsString(needle string, haystack []Person) (Person, error) {
	for _, v := range haystack {
		if v.Name == needle {
			return v, nil
		}
	}
	err := errors.New("Unable to locate given person")
	return Person{}, err
}
