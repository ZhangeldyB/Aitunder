package mongodb

import (
	"Aitunder/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func insertOneUser(user models.User) error {
	inserted, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	log.Info("inserted one user in database with id:", inserted.InsertedID)
	return nil
}

// func deleteOneUser(userId string) {
// 	id, _ := primitive.ObjectIDFromHex(userId)
// 	filter := bson.M{"_id": id}
// 	_, err := collection.DeleteOne(context.Background(), filter)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func getFullFromDB() []models.UserFull {
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var users []models.UserFull
	for cursor.Next(context.Background()) {
		var user models.UserFull
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

func addProfileToUser(userId string, profile models.Profile) error {
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
	log.Info("Added profile to user with ID: ", userId)
	return nil
}
