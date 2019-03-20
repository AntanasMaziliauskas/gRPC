package main

import (
	"flag"
	"fmt"

	"github.com/globalsign/mgo/bson"
)

func main() {

	objectID := flag.String("id", "", "ObjectID for validation")
	flag.Parse()

	if *objectID != "" {
		if !bson.IsObjectIdHex(*objectID) {
			fmt.Println("Provided ID is invalid")
		} else {
			fmt.Println("Geras")
		}
	} else {

		id := bson.NewObjectId()

		fmt.Println(id.Hex())
	}

}
