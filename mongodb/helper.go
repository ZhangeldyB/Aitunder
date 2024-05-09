package mongodb

import (
	"Aitunder/models"
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/gomail.v2"
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

func insertOneProject(project models.Project) (string, error) {
	inserted, err := projectCollection.InsertOne(context.Background(), project)
	if err != nil {
		return "", err
	}
	insertedID, ok := inserted.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert InsertedID to string")
	}
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
func getProjectsByID(userID string) ([]models.Project, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return make([]models.Project, 0), err
	}
	filter := bson.M{"owner": id}
	cursor, err := projectCollection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())
	var projects []models.Project
	for cursor.Next(context.Background()) {
		var project models.Project
		if err = cursor.Decode(&project); err != nil {
			log.Fatal(err)
		}
		projects = append(projects, project)
	}
	return projects, nil
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
			{Key: "$match", Value: bson.D{{Key: "viewedBy", Value: bson.D{{Key: "$nin", Value: bson.A{id}}}}}},
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
func getRandomProject(userID string) (models.Project, error) {
	var project models.Project
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return project, err
	}
	pipeline := []bson.D{
		{
			{Key: "$match", Value: bson.D{{Key: "owner", Value: bson.D{{Key: "$ne", Value: id}}}}},
		},
		{
			{Key: "$match", Value: bson.D{{Key: "viewedBy", Value: bson.D{{Key: "$nin", Value: bson.A{id}}}}}},
		},
		{
			{Key: "$sample", Value: bson.D{{Key: "size", Value: 1}}},
		},
	}
	cursor, err := projectCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Error("error in aggregation ", err)
		return project, err
	}
	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&project); err != nil {
			log.Error("error in decoding Cursor ", err)
			return project, err
		}
	}
	return project, nil
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
func setViewedUs(userID string, visitingID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	visID, err := primitive.ObjectIDFromHex(visitingID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": visID}
	update := bson.M{
		"$addToSet": bson.M{
			"viewedBy": id,
		},
	}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	log.Info("user got viewed, ", visID, " ", id)
	return nil
}

func setViewedProj(userId string, projectID string) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	projId, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": projId}
	update := bson.M{
		"$addToSet": bson.M{
			"viewedBy": id,
		},
	}
	_, err = projectCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	log.Info("project got viewed, ", projId)
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
	filter := bson.M{"accountActivated": true}
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

func addLikedUser(loggedInUserID, userID string) error {
	id, err := primitive.ObjectIDFromHex(loggedInUserID)
	if err != nil {
		return err
	}
	likeID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}

	update := bson.M{
		"$addToSet": bson.M{
			"likedUsers": likeID,
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func sendNotificationEmail(message string, users []models.User) int {
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "aitunderapp.notifications@gmail.com", "hbgr gnxq enfr zmtn")
	start := time.Now()
	count := len(users)
	var wg sync.WaitGroup
	wg.Add(len(users))
	for _, user := range users {
		go func(u models.User) {
			res := sendEmail(u.Email, message, dialer, &wg)
			if !res {
				count--
			}
		}(user)
	}
	wg.Wait()
	fmt.Printf("Time spent for sending %v emails: %v\n", count, time.Since(start))

	return count
}

func sendEmail(email, message string, dialer *gomail.Dialer, wg *sync.WaitGroup) bool {
	defer wg.Done()
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "aitunderapp.notifications@gmail.com")
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Notification")
	mailer.SetBody("text/plain", message)
	if err := dialer.DialAndSend(mailer); err != nil {
		log.Error("Failed to send an email notification for: ", email)
		return false
	}
	return true
}
