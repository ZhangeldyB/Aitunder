package mongodb

import (
	"Aitunder/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func insertOneUser(user models.User) (string, error) {
	inserted, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return "", err
	}
	insertedID, ok := inserted.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert InsertedID to string")
	}
	log.Info("inserted one user in database with id:", inserted.InsertedID)
	return insertedID.Hex(), nil
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

func getOneUserByID(userID string) (models.UserFull, error) {
	var user models.UserFull
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return user, err
	}
	filter := bson.M{"_id": id}
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func getRandomUser(userID string) (models.UserFull, error) {
	var user models.UserFull
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return user, err
	}
	pipeline := []bson.D{
		{
			{Key: "$match", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "$ne", Value: id}}}}},
		},
		{
			{Key: "$match", Value: bson.D{{Key: "viewedUsers", Value: bson.D{{Key: "$nin", Value: bson.A{id}}}}}},
		},
		{
			{Key: "$sample", Value: bson.D{{Key: "size", Value: 1}}},
		},
	}
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Error("error in aggregation ", err)
		return user, err
	}
	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&user); err != nil {
			log.Error("error in decoding Cursor ", err)
			return user, err
		}
	}
	return user, nil
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

func updateOneUserByID(userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}

	update := bson.M{
		"$set": bson.M{
			"accountActivated": true,
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func getAllAuthorizedUsers() ([]models.User, error) {
	filter := bson.M{"accountActivated": true} // Assuming accountActivated field determines authorization
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []models.User
	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
