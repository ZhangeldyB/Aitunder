package mongodb

import (
	"Aitunder/models"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func insertOneUser(user models.User) error {
	inserted, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	fmt.Println("inserted one user in database with id: ", inserted.InsertedID)
	return nil
}

func deleteOneUser(userId string) {
	id, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"_id": id}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
}

// func getAllUsers() []primitive.D {
// 	cursor, err := collection.Find(context.Background(), bson.D{{}})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer cursor.Close(context.Background())
// 	var users []primitive.D
// 	for cursor.Next(context.Background()) {
// 		var user bson.D
// 		if err = cursor.Decode(&user); err != nil {
// 			log.Fatal(err)
// 		}
// 		users = append(users, user)
// 	}
// 	return users
// }

func getAllUsersFromDB() []models.User {
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var users []models.User
	for cursor.Next(context.Background()) {
		var user models.User
		if err = cursor.Decode(&user); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	return users
}

func getOneUserByEmail(email string) (*models.User, error) {
	var user models.User
	filter := bson.M{"email": email}
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func addProfileToUser(userId string, profile models.Portfolio) error {
	id, _ := primitive.ObjectIDFromHex(userId)

	filter := bson.M{"_id": id}

	update := bson.M{
		"$set": bson.M{
			"profile": profile,
		},
	}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	fmt.Println("Added profile to user with ID: ", userId)
	return nil
}
