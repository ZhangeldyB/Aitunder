package mongodb

import (
	"Aitunder/models"
	"context"
	"fmt"
	"log"
)

func InsertOneUser(user models.User) {
	inserted, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted one user in database with id: ", inserted.InsertedID)
}
