package person

import (
	"context"

	"github.com/AntanasMaziliauskas/grpc/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//PersonService interface holds Init, ListPersons, GetOnePerson, GetMultiPerson
//DropOnePerson, DropMultiPerson, UpsertOnePerson, UpsertMultiperson and Ping functions
type PersonService interface {
	Init() error
	ListPersons(context.Context, *api.Empty) (*api.MultiPerson, error)
	GetOnePerson(context.Context, *api.Person) (*api.Person, error)
	GetMultiPerson(context.Context, *api.MultiPerson) (*api.MultiPerson, error)
	DropOnePerson(context.Context, *api.Person) (*api.Empty, error)
	DropMultiPerson(context.Context, *api.MultiPerson) (*api.Empty, error)
	UpsertOnePerson(context.Context, *api.Person) (*api.Empty, error)
	UpsertMultiPerson(context.Context, *api.MultiPerson) (*api.Empty, error)
	Ping(context.Context, *api.PingMessage) (*api.Empty, error)
}

//Person structure holds values of ID, Name, Age and Profession.
type Person struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	Age        int64              `bson:"age"`
	Profession string             `bson:"profession"`
}
